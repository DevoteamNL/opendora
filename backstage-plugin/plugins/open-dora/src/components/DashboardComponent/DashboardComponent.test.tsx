import type { EntityRelation } from '@backstage/catalog-model';
import { ApiProvider } from '@backstage/core-app-api';
import { EntityProvider } from '@backstage/plugin-catalog-react';
import {
  MockConfigApi,
  renderInTestApp,
  TestApiRegistry,
} from '@backstage/test-utils';
import { act, fireEvent, screen } from '@testing-library/react';
import { rest } from 'msw';
import React from 'react';
import { baseUrl, metricUrl } from '../../../testing/mswHandlers';
import {
  DoraDataService,
  doraDataServiceApiRef,
} from '../../services/DoraDataService';
import { server } from '../../setupTests';
import {
  DashboardComponent,
  EntityDashboardComponent,
} from './DashboardComponent';

async function renderComponentWithApis(component: JSX.Element) {
  const mockConfig = new MockConfigApi({
    'open-dora': { apiBaseUrl: baseUrl },
  });

  const apiRegistry = TestApiRegistry.from([
    doraDataServiceApiRef,
    new DoraDataService({ configApi: mockConfig }),
  ]);

  return await renderInTestApp(
    <ApiProvider apis={apiRegistry}>{component}</ApiProvider>,
  );
}

describe('DashboardComponent', () => {
  function renderDashboardComponent() {
    return renderComponentWithApis(<DashboardComponent />);
  }

  it('should show a dropdown with the aggregation choices', async () => {
    const { queryByText, getByText } = await renderDashboardComponent();

    expect(queryByText('Weekly')).not.toBeNull();

    fireEvent.mouseDown(getByText('Weekly'));

    expect(queryByText('Monthly')).not.toBeNull();
    expect(queryByText('Quarterly')).not.toBeNull();
  });

  it('should show the title of the plugin', async () => {
    const { queryByText } = await renderDashboardComponent();

    expect(queryByText('OpenDORA (by Devoteam)')).not.toBeNull();
  });

  it('should show graphs for deployment frequency data', async () => {
    server.use(
      rest.get(metricUrl, (req, res, ctx) => {
        const type = req.url.searchParams.get('type');
        return res(
          ctx.json({
            aggregation: 'weekly',
            dataPoints: [
              { key: `${type}_first_key`, value: 1.0 },
              { key: `${type}_second_key`, value: 1.0 },
              { key: `${type}_third_key`, value: 1.0 },
            ],
          }),
        );
      }),
    );
    const { queryByText } = await renderDashboardComponent();

    expect(queryByText('Deployment Frequency')).not.toBeNull();
    expect(queryByText('df_count_first_key')).not.toBeNull();
    expect(queryByText('df_count_second_key')).not.toBeNull();
    expect(queryByText('df_count_third_key')).not.toBeNull();

    expect(queryByText('Deployment Frequency Average')).not.toBeNull();
    expect(queryByText('df_average_first_key')).not.toBeNull();
    expect(queryByText('df_average_second_key')).not.toBeNull();
    expect(queryByText('df_average_third_key')).not.toBeNull();
  });

  it('should retrieve new data when the aggregation is changed', async () => {
    server.use(
      rest.get(metricUrl, (req, res, ctx) => {
        const params = req.url.searchParams;
        const type = params.get('type');
        const aggregation = params.get('aggregation');

        return res(
          ctx.json({
            aggregation: aggregation,
            dataPoints: [
              { key: `${aggregation}_${type}_first_key`, value: 1.0 },
            ],
          }),
        );
      }),
    );
    const { queryByText, getByText } = await renderDashboardComponent();

    expect(queryByText('weekly_df_count_first_key')).not.toBeNull();
    expect(queryByText('weekly_df_average_first_key')).not.toBeNull();

    expect(queryByText('monthly_df_count_first_key')).toBeNull();
    expect(queryByText('monthly_df_average_first_key')).toBeNull();

    fireEvent.mouseDown(getByText('Weekly'));
    await act(async () => {
      fireEvent.click(screen.getByText('Monthly'));
    });

    expect(queryByText('weekly_df_count_first_key')).toBeNull();
    expect(queryByText('weekly_df_average_first_key')).toBeNull();

    expect(queryByText('monthly_df_count_first_key')).not.toBeNull();
    expect(queryByText('monthly_df_average_first_key')).not.toBeNull();
  });

  it('should show loading indicator when waiting on data to return', async () => {
    jest.useFakeTimers({
      legacyFakeTimers: true,
    });

    server.use(
      rest.get(metricUrl, async (_, res, ctx) => {
        await new Promise(resolve => setTimeout(resolve, 10000));
        return res(
          ctx.json({
            aggregation: 'weekly',
            dataPoints: [{ key: `first_key`, value: 1.0 }],
          }),
        );
      }),
    );

    const { queryByText, queryByRole, findAllByRole, queryAllByText } =
      await renderDashboardComponent();

    expect(await findAllByRole('progressbar')).toHaveLength(2);
    expect(queryByText('first_key')).toBeNull();

    await act(async () => {
      jest.runAllTimers();
    });

    expect(queryByRole('progressbar')).toBeNull();
    expect(queryAllByText('first_key')).toHaveLength(2);
  });

  it('should show the error returned from the service', async () => {
    server.use(
      rest.get(metricUrl, (_, res, ctx) => {
        return res(ctx.status(401));
      }),
    );
    const { queryAllByText, getByText } = await renderDashboardComponent();
    expect(queryAllByText('Error: Unauthorized')).toHaveLength(2);

    server.use(
      rest.get(metricUrl, (_, res) => {
        return res.networkError('Host unreachable');
      }),
    );

    // Trigger another request
    fireEvent.mouseDown(getByText('Weekly'));
    await act(async () => {
      fireEvent.click(screen.getByText('Monthly'));
    });

    expect(queryAllByText('Error: Failed to fetch')).toHaveLength(2);
  });

  it('should show error if there are no datapoints', async () => {
    server.use(
      rest.get(metricUrl, async (_, res, ctx) => {
        return res(
          ctx.json({
            aggregation: 'weekly',
            dataPoints: [],
          }),
        );
      }),
    );
    const { queryAllByText } = await renderDashboardComponent();
    expect(queryAllByText('No data found')).not.toBeNull();
    expect(queryAllByText('No data found')).toHaveLength(2);
  });
});

describe('EntityDashboardComponent', () => {
  function renderEntityDashboardComponent(relations?: EntityRelation[]) {
    server.use(
      rest.get(metricUrl, (req, res, ctx) => {
        const params = req.url.searchParams;
        const type = params.get('type');
        const aggregation = params.get('aggregation');
        const project = params.get('project');
        const team = params.get('team');

        return res(
          ctx.json({
            aggregation: aggregation,
            dataPoints: [
              {
                key: `${project}_${team}_${aggregation}_${type}_first_key`,
                value: 1.0,
              },
            ],
          }),
        );
      }),
    );

    return renderComponentWithApis(
      <EntityProvider
        entity={{
          apiVersion: 'version',
          kind: 'entity-kind',
          metadata: {
            name: 'entity-name',
          },
          relations: relations,
        }}
      >
        <EntityDashboardComponent />
      </EntityProvider>,
    );
  }

  it('should send component info to the service from the context', async () => {
    const { queryByText } = await renderEntityDashboardComponent([
      { targetRef: 'kind:namespace/owner-name', type: 'ownedBy' },
    ]);

    expect(
      queryByText('entity-name_owner-name_weekly_df_average_first_key'),
    ).not.toBeNull();
  });

  it('should send component info without owner info', async () => {
    const { queryByText } = await renderEntityDashboardComponent();

    expect(
      queryByText('entity-name_null_weekly_df_average_first_key'),
    ).not.toBeNull();
  });
});

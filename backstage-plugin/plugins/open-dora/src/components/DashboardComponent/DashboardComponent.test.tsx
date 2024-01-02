import type { EntityRelation } from '@backstage/catalog-model';
import { EntityProvider } from '@backstage/plugin-catalog-react';
import { act, fireEvent, screen } from '@testing-library/react';
import { rest } from 'msw';
import React from 'react';
import { renderComponentWithApis } from '../../../testing/component-test-utils';
import { metricUrl } from '../../../testing/mswHandlers';
import { server } from '../../setupTests';
import {
  DashboardComponent,
  EntityDashboardComponent,
} from './DashboardComponent';

const NUMBER_OF_METRICS = 3;

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

  it('should show graphs for metric data', async () => {
    server.use(
      rest.get(metricUrl, (req, res, ctx) => {
        const type = req.url.searchParams.get('type');
        return res(
          ctx.json({
            aggregation: 'weekly',
            dataPoints: [{ key: `${type}_first_key`, value: 1.0 }],
          }),
        );
      }),
    );
    const { queryByText } = await renderDashboardComponent();

    expect(queryByText('Deployment Frequency')).not.toBeNull();
    expect(queryByText('df_count_first_key')).not.toBeNull();

    expect(queryByText('Deployment Frequency Average')).not.toBeNull();
    expect(queryByText('df_average_first_key')).not.toBeNull();

    expect(queryByText('Median Lead Time for Changes')).not.toBeNull();
    expect(queryByText('mltc_first_key')).not.toBeNull();
  });

  it('should retrieve new data when the aggregation is changed', async () => {
    server.use(
      rest.get(metricUrl, (req, res, ctx) => {
        const params = req.url.searchParams;
        const aggregation = params.get('aggregation');

        return res(
          ctx.json({
            aggregation: aggregation,
            dataPoints: [{ key: `${aggregation}_first_key`, value: 1.0 }],
          }),
        );
      }),
    );
    const { queryAllByText, getByText } = await renderDashboardComponent();

    expect(queryAllByText('weekly_first_key')).toHaveLength(NUMBER_OF_METRICS);
    expect(queryAllByText('monthly_first_key')).toHaveLength(0);

    fireEvent.mouseDown(getByText('Weekly'));
    await act(async () => {
      fireEvent.click(screen.getByText('Monthly'));
    });

    expect(queryAllByText('weekly_first_key')).toHaveLength(0);
    expect(queryAllByText('monthly_first_key')).toHaveLength(NUMBER_OF_METRICS);
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

import { act } from '@testing-library/react';
import { rest } from 'msw';
import React from 'react';
import { renderComponentWithApis } from '../../../testing/component-render-utils';
import { metricUrl } from '../../../testing/mswHandlers';
import '../../i18n';
import { MetricContext } from '../../services/MetricContext';
import { server } from '../../setupTests';
import { ChartGridItem } from './ChartGridItem';

describe('ChartGridItem', () => {
  function renderChartGridItem(
    { type, label }: { type: string; label: string } = {
      type: 'df_count',
      label: 'Deployment Frequency',
    },
  ) {
    return renderComponentWithApis(<ChartGridItem type={type} label={label} />);
  }

  it('should show the a bar chart and title when data is returned from the server', async () => {
    const { queryByText } = await renderChartGridItem();

    expect(queryByText('Deployment Frequency')).not.toBeNull();
    expect(queryByText('10/23')).not.toBeNull();
  });

  it('should show loading indicator when waiting on the data to return', async () => {
    jest.useFakeTimers({
      legacyFakeTimers: true,
    });

    server.use(
      rest.get(metricUrl, async (_, res, ctx) => {
        await new Promise(resolve => setTimeout(resolve, 10000));
        return res(
          ctx.json({
            aggregation: 'weekly',
            dataPoints: [{ key: '10/23', value: 1.0 }],
          }),
        );
      }),
    );

    const { queryByText, queryByRole, findByRole } =
      await renderChartGridItem();

    expect(queryByText('Deployment Frequency')).not.toBeNull();
    expect(await findByRole('progressbar')).not.toBeNull();
    expect(queryByText('10/23')).toBeNull();

    await act(async () => {
      jest.runAllTimers();
    });

    expect(queryByText('Deployment Frequency')).not.toBeNull();
    expect(queryByRole('progressbar')).toBeNull();
    expect(queryByText('10/23')).not.toBeNull();
  });

  it('should show the error returned from the service', async () => {
    server.use(
      rest.get(metricUrl, (_, res, ctx) => {
        return res(ctx.status(401));
      }),
    );
    const { queryByText } = await renderChartGridItem();

    expect(queryByText('Deployment Frequency')).not.toBeNull();
    expect(queryByText('Error: Unauthorized')).not.toBeNull();
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
    const { queryByText } = await renderChartGridItem();
    expect(queryByText('No data found')).not.toBeNull();
  });

  it('should use details from the metric context in data requests', async () => {
    server.use(
      rest.get(metricUrl, async (req, res, ctx) => {
        const params = req.url.searchParams;
        const aggregation = params.get('aggregation');
        return res(
          ctx.json({
            aggregation: aggregation,
            dataPoints: [{ key: `${aggregation}_10/23`, value: 1.0 }],
          }),
        );
      }),
    );

    const { queryByText } = await renderComponentWithApis(
      <MetricContext.Provider value={{ aggregation: 'monthly' }}>
        <ChartGridItem type="df_count" label="Deployment frequency" />
      </MetricContext.Provider>,
    );

    expect(queryByText('weekly_10/23')).toBeNull();
    expect(queryByText('monthly_10/23')).not.toBeNull();
  });
});

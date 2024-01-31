import { act } from '@testing-library/react';
import { rest } from 'msw';
import React from 'react';
import {
  delayRequest,
  renderComponentWithApis,
} from '../../../testing/component-test-utils';
import { metricUrl } from '../../../testing/mswHandlers';
import '../../i18n';
import { MetricContext } from '../../services/MetricContext';
import { server } from '../../setupTests';
import { ChartGridItem } from './ChartGridItem';

describe('ChartGridItem', () => {
  function renderChartGridItem() {
    return renderComponentWithApis(
      <ChartGridItem type="df_count" label="Deployment Frequency" />,
    );
  }

  it('should show the a bar chart and title when data is returned from the server', async () => {
    const { queryByText, queryAllByText } = await renderChartGridItem();

    expect(queryByText('Deployment Frequency')).toBeInTheDocument();
    expect(queryAllByText('10/23')[0]).toBeInTheDocument();
  });

  it('should show loading indicator when waiting on the data to return', async () => {
    delayRequest(
      { aggregation: 'weekly', dataPoints: [{ key: '10/23', value: 1.0 }] },
      metricUrl,
    );

    const { queryByText, queryAllByText, queryByRole, findByRole } =
      await renderChartGridItem();

    expect(queryByText('Deployment Frequency')).toBeInTheDocument();
    expect(await findByRole('progressbar')).toBeInTheDocument();

    await act(async () => {
      jest.runAllTimers();
    });

    expect(queryByText('Deployment Frequency')).toBeInTheDocument();
    expect(queryByRole('progressbar')).not.toBeInTheDocument();
    expect(queryAllByText('10/23')[0]).toBeInTheDocument();

    jest.runOnlyPendingTimers();
    jest.useRealTimers();
  });

  it('should show the error returned from the service', async () => {
    server.use(
      rest.get(metricUrl, (_, res, ctx) => {
        return res(ctx.status(401));
      }),
    );
    const { queryByText } = await renderChartGridItem();

    expect(queryByText('Deployment Frequency')).toBeInTheDocument();
    expect(queryByText('Error: Unauthorized')).toBeInTheDocument();
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
    expect(queryByText('No data found')).toBeInTheDocument();
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

    const { queryByText, queryAllByText } = await renderComponentWithApis(
      <MetricContext.Provider value={{ aggregation: 'monthly' }}>
        <ChartGridItem type="df_count" label="Deployment frequency" />
      </MetricContext.Provider>,
    );

    expect(queryByText('weekly_10/23')).not.toBeInTheDocument();
    expect(queryAllByText('monthly_10/23')[0]).toBeInTheDocument();
  });
});

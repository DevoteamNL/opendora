import { render } from '@testing-library/react';
import React from 'react';
import { BarChartComponent } from './BarChartComponent';

describe('BarChartComponent', () => {
  it('should create a bar chart with data from input', async () => {
    const { queryAllByText } = render(
      <BarChartComponent
        metricData={{
          aggregation: 'weekly',
          dataPoints: [
            { key: 'data_key', value: 1.0 },
            { key: 'data_key_2', value: 1.0 },
          ],
        }}
      />,
    );

    expect(queryAllByText('data_key')[0]).toBeInTheDocument();
    expect(queryAllByText('data_key_2')[0]).toBeInTheDocument();
  });
});

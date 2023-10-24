import React from 'react';
import { render } from '@testing-library/react';
import { BarChartComponent } from './BarChartComponent';

describe('BarChartComponent', () => {
  it('should create a bar chart with data from input', async () => {
    const { queryByText } = render(
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

    expect(queryByText('data_key')).not.toBeNull();
    expect(queryByText('data_key_2')).not.toBeNull();
  });
});

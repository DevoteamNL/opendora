import React from 'react';
import { DashboardComponent } from './DashboardComponent';
import { renderInTestApp, TestApiRegistry } from '@backstage/test-utils';
import { groupDataServiceApiRef } from '../../services/GroupDataService';
import { ApiProvider } from '@backstage/core-app-api';
import { fireEvent, screen, act } from '@testing-library/react';
import { MetricData } from '../../models/MetricData';

describe('DashboardComponent', () => {
  async function renderDashboardComponent(mockData?: jest.Mock<MetricData>) {
    const groupDataApiMock = {
      retrieveMetricDataPoints:
        mockData ??
        jest.fn().mockResolvedValue({
          aggregation: 'weekly',
          dataPoints: [{ key: '10/23', value: 1.0 }],
        }),
    };
    const apiRegistry = TestApiRegistry.from([
      groupDataServiceApiRef,
      groupDataApiMock,
    ]);

    return await renderInTestApp(
      <ApiProvider apis={apiRegistry}>
        <DashboardComponent />
      </ApiProvider>,
    );
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

  it('should show a graph for deployment frequency data', async () => {
    const { queryByText } = await renderDashboardComponent(
      jest.fn().mockResolvedValue({
        aggregation: 'weekly',
        dataPoints: [
          { key: 'first_key', value: 1.0 },
          { key: 'second_key', value: 1.0 },
          { key: 'third_key', value: 1.0 },
        ],
      }),
    );

    expect(queryByText('Deployment Frequency')).not.toBeNull();
    expect(queryByText('first_key')).not.toBeNull();
    expect(queryByText('second_key')).not.toBeNull();
    expect(queryByText('third_key')).not.toBeNull();
  });

  it('should retrieve new data when the aggregation is changed', async () => {
    const apiMock = jest
      .fn()
      .mockResolvedValueOnce({
        aggregation: 'weekly',
        dataPoints: [
          { key: 'first_key', value: 1.0 },
          { key: 'second_key', value: 1.0 },
          { key: 'third_key', value: 1.0 },
        ],
      })
      .mockResolvedValueOnce({
        aggregation: 'monthly',
        dataPoints: [
          { key: 'new_first_key', value: 1.0 },
          { key: 'new_second_key', value: 1.0 },
          { key: 'new_third_key', value: 1.0 },
        ],
      });
    const { queryByText, getByText } = await renderDashboardComponent(apiMock);

    expect(apiMock).toHaveBeenLastCalledWith({
      type: 'df_count',
      aggregation: 'weekly',
    });

    expect(queryByText('first_key')).not.toBeNull();
    expect(queryByText('second_key')).not.toBeNull();
    expect(queryByText('third_key')).not.toBeNull();

    fireEvent.mouseDown(getByText('Weekly'));
    await act(async () => {
      fireEvent.click(screen.getByText('Monthly'));
    });

    expect(apiMock).toHaveBeenLastCalledWith({
      type: 'df_count',
      aggregation: 'monthly',
    });

    expect(queryByText('first_key')).toBeNull();
    expect(queryByText('second_key')).toBeNull();
    expect(queryByText('third_key')).toBeNull();

    expect(queryByText('new_first_key')).not.toBeNull();
    expect(queryByText('new_second_key')).not.toBeNull();
    expect(queryByText('new_third_key')).not.toBeNull();
  });

  it('should show loading indicator when waiting on data to return', async () => {
    jest.useFakeTimers();

    const apiMock = jest.fn().mockImplementation(() => {
      return new Promise(resolve => {
        setTimeout(() => {
          resolve({
            aggregation: 'monthly',
            dataPoints: [{ key: 'data_key', value: 1.0 }],
          });
        }, 1000);
      });
    });
    const { queryByText, queryByRole, findByRole } =
      await renderDashboardComponent(apiMock);

    expect(await findByRole('progressbar')).not.toBeNull();
    expect(queryByText('data_key')).toBeNull();

    await act(async () => {
      jest.runAllTimers();
    });

    expect(queryByRole('progressbar')).toBeNull();
    expect(queryByText('data_key')).not.toBeNull();
  });
});
import { MockConfigApi } from '@backstage/test-utils';
import { http, HttpResponse } from 'msw';
import { server } from '../setupTests';
import { GroupDataService } from './GroupDataService';

function createService() {
  server.use(
    http.get('http://localhost:10666/dora/api/metric', ({ request }) => {
      const params = new URL(request.url).searchParams;
      const type = params.get('type');
      const aggregation = params.get('aggregation');
      const project = params.get('project');
      const team = params.get('team');

      return HttpResponse.json({
        aggregation: aggregation || 'weekly',
        dataPoints: [
          {
            key: `${project}_${team}_${aggregation}_${type}_first_key`,
            value: 1.0,
          },
        ],
      });
    }),
  );
  const mockConfig = new MockConfigApi({
    'open-dora': { apiBaseUrl: 'http://localhost:10666' },
  });

  return new GroupDataService({ configApi: mockConfig });
}

describe('GroupDataService', () => {
  it('should return data from the server', async () => {
    const service = createService();

    expect(
      await service.retrieveMetricDataPoints({ type: 'df_count' }),
    ).toEqual({
      aggregation: 'weekly',
      dataPoints: [{ key: 'null_null_null_df_count_first_key', value: 1 }],
    });
  });

  it('should use provided details in the query parameters', async () => {
    const service = createService();

    expect(
      await service.retrieveMetricDataPoints({
        type: 'df_count',
        aggregation: 'monthly',
        project: 'project1',
        team: 'team1',
      }),
    ).toEqual({
      aggregation: 'monthly',
      dataPoints: [
        { key: 'project1_team1_monthly_df_count_first_key', value: 1 },
      ],
    });
  });

  it('should throw an error if the response does not contain metric data', async () => {
    const service = createService();

    server.use(
      http.get('http://localhost:10666/dora/api/metric', () => {
        return HttpResponse.json({ other: 'data' });
      }),
    );
    await expect(
      service.retrieveMetricDataPoints({
        type: 'df_count',
      }),
    ).rejects.toEqual(new Error('Unexpected response'));
  });

  it('should throw an error when the server is unreachable', async () => {
    const service = createService();

    server.use(
      http.get('http://localhost:10666/dora/api/metric', () => {
        return HttpResponse.error();
      }),
    );
    await expect(
      service.retrieveMetricDataPoints({
        type: 'df_count',
      }),
    ).rejects.toEqual(new Error('Failed to fetch'));
  });

  it('should throw an error when the server returns a non-ok status', async () => {
    const service = createService();

    server.use(
      http.get('http://localhost:10666/dora/api/metric', () => {
        return new HttpResponse(null, { status: 401 });
      }),
    );
    await expect(
      service.retrieveMetricDataPoints({
        type: 'df_count',
      }),
    ).rejects.toEqual(new Error('Unauthorized'));

    server.use(
      http.get('http://localhost:10666/dora/api/metric', () => {
        return new HttpResponse(null, { status: 500 });
      }),
    );
    await expect(
      service.retrieveMetricDataPoints({
        type: 'df_count',
      }),
    ).rejects.toEqual(new Error('Internal Server Error'));
  });
});

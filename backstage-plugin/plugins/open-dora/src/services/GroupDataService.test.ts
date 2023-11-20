import { MockConfigApi } from '@backstage/test-utils';
import { rest } from 'msw';
import { baseUrl, metricUrl } from '../../testing/mswHandlers';
import { server } from '../setupTests';
import { GroupDataService } from './GroupDataService';

function createService() {
  server.use(
    rest.get(metricUrl, (req, res, ctx) => {
      const params = req.url.searchParams;
      const type = params.get('type');
      const aggregation = params.get('aggregation');
      const project = params.get('project');
      const team = params.get('team');

      return res(
        ctx.json({
          aggregation: aggregation || 'weekly',
          dataPoints: [
            {
              key: `${project}_${team}_${aggregation}_${type}_first_key`,
              value: 2.3,
            },
          ],
        }),
      );
    }),
  );
  const mockConfig = new MockConfigApi({
    'open-dora': { apiBaseUrl: baseUrl },
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
      dataPoints: [{ key: 'null_null_null_df_count_first_key', value: 2.3 }],
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
        { key: 'project1_team1_monthly_df_count_first_key', value: 2.3 },
      ],
    });
  });

  it('should throw an error if the response does not contain metric data', async () => {
    const service = createService();

    server.use(
      rest.get(metricUrl, (_, res, ctx) => {
        return res(ctx.json({ other: 'data' }));
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
      rest.get(metricUrl, (_, res) => {
        return res.networkError('Host unreachable');
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
      rest.get(metricUrl, (_, res, ctx) => {
        return res(ctx.status(401));
      }),
    );
    await expect(
      service.retrieveMetricDataPoints({
        type: 'df_count',
      }),
    ).rejects.toEqual(new Error('Unauthorized'));

    server.use(
      rest.get(metricUrl, (_, res, ctx) => {
        return res(ctx.status(500));
      }),
    );
    await expect(
      service.retrieveMetricDataPoints({
        type: 'df_count',
      }),
    ).rejects.toEqual(new Error('Internal Server Error'));
  });
});

import { rest } from 'msw';

export const baseUrl = 'http://localhost:10666';
export const metricUrl = `${baseUrl}/dora/api/metric`;

export const handlers = [
  rest.get(metricUrl, (_, res, ctx) => {
    return res(
      ctx.json({
        aggregation: 'weekly',
        dataPoints: [{ key: '10/23', value: 1.0 }],
      }),
    );
  }),
];

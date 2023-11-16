import { rest } from 'msw';

export const handlers = [
  rest.get('http://localhost:10666/dora/api/metric', (_, res, ctx) => {
    return res(
      ctx.json({
        aggregation: 'weekly',
        dataPoints: [{ key: '10/23', value: 1.0 }],
      }),
    );
  }),
];

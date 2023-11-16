import { http, HttpResponse } from 'msw';

export const handlers = [
  http.get('http://localhost:10666/dora/api/metric', () => {
    return HttpResponse.json({
      aggregation: 'weekly',
      dataPoints: [{ key: '10/23', value: 1.0 }],
    });
  }),
];

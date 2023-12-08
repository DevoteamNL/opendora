import { ApiProvider } from '@backstage/core-app-api';
import {
  MockConfigApi,
  renderInTestApp,
  TestApiRegistry,
} from '@backstage/test-utils';
import { rest } from 'msw';
import React from 'react';
import {
  DoraDataService,
  doraDataServiceApiRef,
} from '../src/services/DoraDataService';
import { server } from '../src/setupTests';
import { baseUrl } from './mswHandlers';

export async function renderComponentWithApis(component: JSX.Element) {
  const mockConfig = new MockConfigApi({
    'open-dora': { apiBaseUrl: baseUrl },
  });

  const apiRegistry = TestApiRegistry.from([
    doraDataServiceApiRef,
    new DoraDataService({ configApi: mockConfig }),
  ]);

  return await renderInTestApp(
    <ApiProvider apis={apiRegistry}>{component}</ApiProvider>,
  );
}

export function delayRequest(request: Object, url: string, delay = 10000) {
  jest.useFakeTimers({
    legacyFakeTimers: true,
  });

  server.use(
    rest.get(url, async (_, res, ctx) => {
      await new Promise(resolve => setTimeout(resolve, delay));
      return res(ctx.json(request));
    }),
  );
}

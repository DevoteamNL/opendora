import { ApiProvider } from '@backstage/core-app-api';
import {
  MockConfigApi,
  renderInTestApp,
  TestApiRegistry,
} from '@backstage/test-utils';
import React from 'react';
import {
  DoraDataService,
  doraDataServiceApiRef,
} from '../src/services/DoraDataService';
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

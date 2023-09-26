import {
  configApiRef,
  createApiFactory,
  createPlugin,
  createRoutableExtension,
} from '@backstage/core-plugin-api';
import { rootRouteRef } from './routes';
import {
  GroupDataService,
  groupDataServiceApiRef,
} from './services/GroupDataService';

export const doraPluginPlugin = createPlugin({
  id: 'dora-plugin',
  routes: {
    root: rootRouteRef,
  },
  apis: [
    createApiFactory({
      api: groupDataServiceApiRef,
      deps: { configApi: configApiRef },
      factory: ({ configApi }) => new GroupDataService({ configApi }),
    }),
  ],
});

export const DoraPluginPage = doraPluginPlugin.provide(
  createRoutableExtension({
    name: 'DoraPluginPage',
    component: () =>
      import('./components/DashboardComponent').then(m => m.DashboardComponent),
    mountPoint: rootRouteRef,
  }),
);

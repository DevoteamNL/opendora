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

export const openDoraPlugin = createPlugin({
  id: 'opendora',
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

export const OpenDoraPluginPage = openDoraPlugin.provide(
  createRoutableExtension({
    name: 'OpenDoraPluginPage',
    component: () =>
      import('./components/DashboardComponent').then(m => m.DashboardComponent),
    mountPoint: rootRouteRef,
  }),
);

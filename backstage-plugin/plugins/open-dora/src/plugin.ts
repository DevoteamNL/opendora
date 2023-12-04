import {
  configApiRef,
  createApiFactory,
  createPlugin,
  createRoutableExtension,
} from '@backstage/core-plugin-api';
import { rootRouteRef } from './routes';
import {
  DoraDataService,
  doraDataServiceApiRef,
} from './services/DoraDataService';

export const openDoraPlugin = createPlugin({
  id: 'opendora',
  routes: {
    root: rootRouteRef,
  },
  apis: [
    createApiFactory({
      api: doraDataServiceApiRef,
      deps: { configApi: configApiRef },
      factory: ({ configApi }) => new DoraDataService({ configApi }),
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

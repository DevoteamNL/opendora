import { createPlugin, createRoutableExtension } from '@backstage/core-plugin-api';

import { rootRouteRef } from './routes';

export const doraPluginPlugin = createPlugin({
  id: 'dora-plugin',
  routes: {
    root: rootRouteRef,
  },
});

export const DoraPluginPage = doraPluginPlugin.provide(
  createRoutableExtension({
    name: 'DoraPluginPage',
    component: () =>
      import('./components/ExampleComponent').then(m => m.ExampleComponent),
    mountPoint: rootRouteRef,
  }),
);

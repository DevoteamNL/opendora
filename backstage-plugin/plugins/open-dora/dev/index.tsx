import { createDevApp } from '@backstage/dev-utils';
import { MockConfigApi } from '@backstage/test-utils';
import React from 'react';
import { openDoraPlugin, OpenDoraPluginPage } from '../src';
import {
  DoraDataService,
  doraDataServiceApiRef,
} from '../src/services/DoraDataService';

const mockConfig = new MockConfigApi({
  'open-dora': {
    apiBaseUrl: 'http://localhost:10666',
  },
});

createDevApp()
  .registerPlugin(openDoraPlugin)
  .registerApi({
    api: doraDataServiceApiRef,
    deps: {},
    factory: () => new DoraDataService({ configApi: mockConfig }),
  })
  .addPage({
    element: <OpenDoraPluginPage />,
    title: 'OpenDORA Dev Page',
    path: '/open-dora-plugin',
  })
  .render();

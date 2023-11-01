import { createDevApp } from '@backstage/dev-utils';
import React from 'react';
import { openDoraPlugin, OpenDoraPluginPage } from '../src';
import { MockConfigApi } from '@backstage/test-utils';
import {
  GroupDataService,
  groupDataServiceApiRef,
} from '../src/services/GroupDataService';

const mockConfig = new MockConfigApi({
  'open-dora': {
    apiBaseUrl: 'http://localhost:10666',
  },
});

createDevApp()
  .registerPlugin(openDoraPlugin)
  .registerApi({
    api: groupDataServiceApiRef,
    deps: {},
    factory: () => new GroupDataService({ configApi: mockConfig }),
  })
  .addPage({
    element: <OpenDoraPluginPage />,
    title: 'OpenDORA Dev Page',
    path: '/open-dora-plugin',
  })
  .render();

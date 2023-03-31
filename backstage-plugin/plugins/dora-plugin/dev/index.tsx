import React from 'react';
import { createDevApp } from '@backstage/dev-utils';
import { doraPluginPlugin, DoraPluginPage } from '../src/plugin';

createDevApp()
  .registerPlugin(doraPluginPlugin)
  .addPage({
    element: <DoraPluginPage />,
    title: 'Root Page',
    path: '/dora-plugin'
  })
  .render();

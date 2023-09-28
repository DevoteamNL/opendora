import { createDevApp } from '@backstage/dev-utils';
import React from 'react';
import { DoraPluginPage, doraPluginPlugin } from '../src/plugin';

createDevApp()
  .registerPlugin(doraPluginPlugin)
  .addPage({
    element: <DoraPluginPage />,
    title: 'Root Page',
    path: '/dora-plugin',
  })
  .render();

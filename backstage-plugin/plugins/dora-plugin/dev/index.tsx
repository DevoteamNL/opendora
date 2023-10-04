import { createDevApp } from '@backstage/dev-utils';
import React from 'react';
import { OpenDoraPluginPage, openDoraPlugin } from '../src/plugin';

createDevApp()
  .registerPlugin(openDoraPlugin)
  .addPage({
    element: <OpenDoraPluginPage />,
    title: 'Root Page',
    path: '/open-dora-plugin',
  })
  .render();

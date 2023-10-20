import { createDevApp } from '@backstage/dev-utils';
import React from 'react';
import { openDoraPlugin, OpenDoraPluginPage } from '../src';

createDevApp()
  .registerPlugin(openDoraPlugin)
  .addPage({
    element: <OpenDoraPluginPage />,
    title: 'OpenDORA Dev Page',
    path: '/open-dora-plugin',
  })
  .render();

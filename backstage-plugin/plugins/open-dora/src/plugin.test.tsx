import { renderInTestApp } from '@backstage/test-utils';
import { OpenDoraPluginPage, openDoraPlugin } from './plugin';
import React from 'react';
import { Route, Routes } from 'react-router';

describe('dora-plugin', () => {
  it('should export plugin', () => {
    expect(openDoraPlugin).toBeDefined();
  });

  it('should mount the plugin page in the test app', async () => {
    const { getByText } = await renderInTestApp(
      <Routes>
        <Route path="/open-dora" element={<OpenDoraPluginPage />} />
      </Routes>,
      {
        routeEntries: ['/open-dora'],
      },
    );

    expect(getByText('OpenDORA (by Devoteam)')).toBeInTheDocument();
  });
});

import { OpenDoraPluginPage, openDoraPlugin } from './plugin';

describe('dora-plugin', () => {
  it('should export plugin', () => {
    expect(openDoraPlugin).toBeDefined();
  });

  it('should export plugin page', () => {
    expect(OpenDoraPluginPage).toBeDefined();
  });
});

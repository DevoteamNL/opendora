# OpenDORA Plugin for Backstage

Welcome to the OpenDORA plugin!

This plugin allows you to see DORA metrics for the teams within Backstage.

## Setup

1. Install this plugin:

```bash
# From your Backstage root directory
yarn --cwd packages/app add open-dora-backstage-plugin
```

<!-- TODO Replace with a link to installation instructions for the helm chart -->

2. Make sure the [OpenDORA backend](https://github.com/DevoteamNL/opendora) is deployed.

3. Configure the url from which the OpenDORA API is accessible.

```yaml
# app-config.yaml
open-dora:
  apiBaseUrl: http://localhost:10666
```

### Entity Pages

1. Add a route to the plugin page:

```jsx
// In packages/app/src/App.tsx
import { OpenDoraPluginPage } from 'open-dora-backstage-plugin';

...
const routes = (
  <FlatRoutes>
    {/* other routes... */}
    <Route path="/open-dora" element={<OpenDoraPluginPage />} />
  </FlatRoutes>
);
```

2. Add the plugin as a tab to your side-navigation:

```jsx
// In packages/app/src/components/Root/Root.tsx
export const Root = ({ children }: PropsWithChildren<{}>) => (
  <SidebarPage>
    <Sidebar>
      {/* other sidebar groups... */}
      <SidebarGroup label="Menu" icon={<MenuIcon />}>
        {/* other sidebar items... */}
        <SidebarItem icon={ExtensionIcon} to="open-dora" text="OpenDORA" />
      </SidebarGroup>

      {/* other sidebar groups... */}
    </Sidebar>
    {children}
  </SidebarPage>
);
```

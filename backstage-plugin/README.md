# OpenDORA Plugin

## Setup

Follow [installation instructions](./plugins/open-dora/README.md).

## Local development

First install dependencies:

```
yarn install
```

Then you can start the local dev page at `plugins/open-dora/dev` with:

```
yarn dev
```

For testing this package in a local Backstage installation you can use [local package linking](https://backstage.io/docs/local-dev/linking-local-packages/)

Local Backstage `package.json`:

```json
"packages": [
  "packages/*",
  "plugins/*",
  "../open-dora/backstage-plugin/plugins/open-dora", // New path added to work on OpenDORA
],
```

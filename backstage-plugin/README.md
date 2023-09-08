# Backstage Dora Plugin

This is a example Backstage app used as an environment to develop the plugin. It is made up of 2 packages and 1 plugin.

### Prerequisites
Backstage Docs: https://backstage.io/docs/getting-started/#prerequisites

**Note: If you are using WSL and the project is on a Windows NTFS drive make sure to use WSL1. Yarn will be extremely slow installing dependencies otherwise.**

## Installation

```
yarn install
```

## [`app`](packages/app)
This package contains all of the frontend components to make up the main UI of the Backstage dashboard. 

Start up the webpack server by running:
```
yarn start 
```
This should be available at http://localhost:3000/

## [`backend`](packages/backend)
This package contains the Node backend to serve data and run actions for the Backstage dashboard.

Start up the server by running:
```
yarn start-backend
```
This should be available at http://localhost:7007/

[More details including configuration of tokens for deploying components](packages\backend\README.md)

## [`dora-plugin`](plugins/dora-plugin/)

This is the frontend plugin that includes each page and component to visualize the Dora metrics. No installation or starting is needed for this as it will be started with the `app` package. 

This should be available at http://localhost:3000/dora-plugin

[More details](plugins\dora-plugin\README.md)
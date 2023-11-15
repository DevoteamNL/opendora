# OpenDORA (by Devoteam)

This repo contains all the projects needed to start development on the OpenDORA backstage plugin for DORA metrics.

![Backstage Plugin](.github/workflows/pr-backstage-plugin-workflow.yaml/badge.svg)
![API](.github/workflows/pr-go-workflow.yaml/badge.svg)
[![Plugin version](https://img.shields.io/github/package-json/v/devoteamnl/opendora?label=plugin&filename=backstage-plugin%2Fplugins%2Fopen-dora%2Fpackage.json)](https://www.npmjs.com/package/@devoteam-nl/open-dora-backstage-plugin)

![Screenshot of the main OpenDORA dashboard](screenshot.png)

## Contents

### dora-api-mock

This is a basic Spring application used to provide a mock for the DORA API metric datapoints. This can be used to develop the frontend plugin without needing to setup the local backend API.

[More details](dora-api-mock/README.md)

### backstage-plugin

This is a [Backstage](https://backstage.io) plugin setup used to develop the frontend plugin.

[More details](backstage-plugin/README.md)

### devlake-go

This contains the docker image (todo) and configuration scripts to setup DevLake to properly ingest the DORA metrics from a GitLab repo and group them according to Backstage groups.

[More details](devlake-go/README.md)

### dev-environment

Contains an initial docker compose with the services needed to test the metrics in the dev's local environment.

[More details](dev-environment/README.md)

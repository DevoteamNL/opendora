![Backstage Plugin](https://github.com/DevoteamNL/opendora/actions/workflows/pr-backstage-plugin-workflow.yaml/badge.svg?branch=main)
![API](https://github.com/DevoteamNL/opendora/actions/workflows/pr-go-workflow.yaml/badge.svg?branch=main)
[![Plugin version](https://img.shields.io/github/package-json/v/devoteamnl/opendora?label=plugin&filename=backstage-plugin%2Fplugins%2Fopen-dora%2Fpackage.json)](https://www.npmjs.com/package/@devoteam-nl/open-dora-backstage-plugin)

# OpenDORA

_Team performance observability for your organization._

OpenDORA includes an open-source plugin for [Backstage](https://backstage.io), a popular developer portal platform. It integrates with [Apache DevLake](https://devlake.apache.org) to organize and aggregate data from deployment and project management tooling like Gitlab, GitHub, Jira, and Jenkins. OpenDORA extracts meaningful insights from this data through its API, and renders dashboards within Backstage that provide insights on the teams' performance.

## What are DORA metrics?

DORA, short for [DevOps Research and Assessment](https://dora.dev), was created 6 years ago by researchers based on data from thousands of teams, to find a reliable way to measure software team performance.

The key 4 DORA metrics are:
- Deployment frequency
- Lead time to changes
- Change failure rate
- Mean time to recovery

Today, DORA metrics are widely accepted as a framework to determine stability and velocity of software teams. They provide a benchmark for determining the maturity of software teams, helping set a path towards high performance.

## How does it work?

OpenDORA has largely a pluggable architecture, with some opinionated tooling choices that we can talk about at length. Ingestor scripts running on schedule fetch data from external tools (Gitlab, Jira, etc) and push them to DevLake. The `devlake-go` API exposes this data and provides endpoints to retrieve it from DevLake. The `backstage-plugin` utilizes this API and renders the results on the Backstage frontend. This is a React plugin based on Material UI.

![Screenshot of the main OpenDORA dashboard](architecture-diagram.png)

The resulting dashboard as rendered within Backstage:

![Screenshot of the main OpenDORA dashboard](screenshot-plugin.png)

Below is a description of the main components of this repo.

### backstage-plugin

This is a [Backstage](https://backstage.io) plugin setup used to develop the frontend plugin.

[More details](backstage-plugin/README.md)

### devlake-go

This contains the docker image and configuration scripts to setup DevLake to properly ingest the DORA metrics from a GitLab repo and group them according to Backstage groups.

[More details](devlake-go/README.md)

### dora-api-mock

This is a basic Spring application used to provide a mock for the DORA API metric datapoints. This can be used to develop the frontend plugin without needing to setup the local backend API.

[More details](dora-api-mock/README.md)

### dev-environment

Contains an initial docker compose with the services needed to test the metrics in the dev's local environment.

[More details](dev-environment/README.md)

## Setup

Goto [Plugin setup documentation](backstage-plugin/plugins/open-dora/README.md) and follow the steps to install and setup the plugin in your Backstage environment.

## Contributor Guide

The OpenDORA team are looking for feedback and contributions. We encourage you to try out the Backstage plugin or other parts of the solution. Please create an issue on GitHub or reach out to us at opendora@devoteam.com.

We also accept pull requests for the code and documentation. Kindly use the discussions in the issue, and please read the [contribution guidelines](CONTRIBUTING.md) page.

![Backstage Plugin](https://github.com/DevoteamNL/opendora/actions/workflows/pr-backstage-plugin-workflow.yaml/badge.svg?branch=main)
![API](https://github.com/DevoteamNL/opendora/actions/workflows/pr-go-workflow.yaml/badge.svg?branch=main)
[![Plugin version](https://img.shields.io/github/package-json/v/devoteamnl/opendora?label=plugin&filename=backstage-plugin%2Fplugins%2Fopen-dora%2Fpackage.json)](https://www.npmjs.com/package/@devoteam-nl/open-dora-backstage-plugin)

# OpenDORA

_Team performance observability for your organization._

OpenDORA includes an open-source plugin for [Backstage](https://backstage.io), a popular developer portal platform. It integrates with [Apache DevLake](https://devlake.apache.org) to organize and aggregate data from deployment and project management tooling like Gitlab, GitHub, Jira, and Jenkins. OpenDORA extracts meaningful insights from this data through its API, and renders dashboards within Backstage that provide insights on the teams' performance.

## What are DORA metrics?

DORA, short for [DevOps Research and Assessment](https://dora.dev), is a research program that began in 2013 investigating ways to measure and improve software delivery, operations, team, and organizational performance. The program has collected data from thousands of professionals and [publishes findings in annual reports](https://dora.dev/publications/).

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

# DORA Plugin Development Environment

## Prerequisites

- Docker Runtime
- Docker Compose

## Configuring the Environment

### Apache DevLake

Before running the Dev Environment, some minor configurations needs to be performed. This development environment
already provides default values for the majority of them but at least an encryption key (`DEVLAKE_ENCRYPTION_SECRET`) needs to be provided.

To generate one, run the below:

```shell
openssl rand -base64 2000 | tr -dc 'A-Z' | fold -w 128 | head -n 1
```

It can be configured globally, and it will be picked up by the docker compose, or it can be added to the `.env` file.

Some environment variables can be customised. The list below has all of them as well as their default values:

| Variable Name                  | Default Value                                                          |
|--------------------------------|------------------------------------------------------------------------|
| `DEVLAKE_ENCRYPTION_SECRET`    |                                                                        |
| `DEVLAKE_MYSQL_USER`           | `merico`                                                               |
| `DEVLAKE_MYSQL_PASSWORD`       | `merico`                                                               |
| `DEVLAKE_MYSQL_CONNECTION_URL` | `mysql://merico:merico@mysql:3306/lake?charset=utf8mb4&parseTime=True` |
| `DEVLAKE_ADMIN_USER`           | `devlake`                                                              |
| `DEVLAKE_ADMIN_PASSWORD`       | `merico`                                                               |

For more information on how to configure Apache DevLake with docker compose, check their official reference: [Launch DevLake with Docker Compose](https://devlake.apache.org/docs/v0.18/GettingStarted/DockerComposeSetup).

Also, make sure to check [DevLake's Configuration documentation](https://devlake.apache.org/docs/v0.18/Configuration) to correctly configure
the connection, data scope, etc. For the DORA metrics, the most important part is the [transformations](https://devlake.apache.org/docs/v0.18/Configuration/Tutorial#step-3---add-transformations-optional), 
so make sure to configure it accordingly (job names) for the metrics to be properly calculated.

## Running the Dev Environment

- Docker Compose v1.*

```shell
docker-compose up -d
```

- Docker Compose v2.*

```shell
docker compose up -d
```

A Grafana container is also available in case you want to access the default dashboards delivered with Apache DevLake.

To start it, run the below:

- Docker Compose v1.*

```shell
docker-compose --profile dashboard up -d
```

- Docker Compose v2.*

```shell
docker compose --profile dashboard up -d
```

After that, the services will be accessible on the below:

- DevLake Configuration UI: http://localhost:4000
- Grafana Dashboards (if started with profile `dashboard`): http://localhost:3000 (initial credential: `admin`/`admin`)

For more information about Docker Compose, check their official documentation: [Docker Compose](https://docs.docker.com/compose/)

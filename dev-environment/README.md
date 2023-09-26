# DORA Plugin Development Environment

## Prerequisites

- Docker Runtime
- Docker Compose

## Configuring the Environment

### Apache DevLake

Before running the Dev Environment, some minor configurations needs to be performed. This development environment
already provide default values for the majority of them but at least an encryption key (`DEVLAKE_ENCRYPTION_SECRET`) needs to be provided.

To generate one, run the below:

```shell
openssl rand -base64 2000 | tr -dc 'A-Z' | fold -w 128 | head -n 1
```

It can be configured globally, and it will be pick up by the docker compose, or it can be added to the `.env` file.

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
docker-compose -p dashboard up -d
```

- Docker Compose v2.*

```shell
docker compose -p dashboard up -d
```

For more information about Docker Compose, check their official documentation: [Docker Compose](https://docs.docker.com/compose/)

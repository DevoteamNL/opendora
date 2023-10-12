# OpenDORA

This repo contains all the projects needed to start development on the OpenDORA backstage plugin for DORA metrics.

## Contents

### dora-api-mock

This is a basic Spring application used to provide a mock for the DORA API metric datapoints. This can be used to develop the frontend plugin without needing to setup the local backend API.

[More details](dora-api-mock/README.md)

### backstage-plugin

This is a scaffold [Backstage](https://backstage.io) app used to configure an environment to develop the plugin.

[More details](backstage-plugin/README.md)

### devlake-go

This contains the docker image (todo) and configuration scripts to setup DevLake to properly ingest the DORA metrics from a GitLab repo and group them according to Backstage groups.

[More details](devlake-go/README.md)

### dev-environment

Contains an initial docker compose with the services needed to test the metrics in the dev's local environment.

[More details](dev-environment/README.md)

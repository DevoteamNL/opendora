# Backstage Dora Plugin

This repo contains all the projects needed to start development on the backstage plugin for Dora metrics.

## Contents

### backstage-mock-main

This is a basic Spring application used to provide a mock API with data representing metrics stored in a DevLake. This can be used to develop the frontend plugin without needing to setup the local backend API or DevLake.

[More details](backstage-mock-main/README.md)

### backstage-plugin

This is a scaffold [Backstage](https://backstage.io) app used to configure an environment to develop the plugin.

[More details](backstage-plugin/README.md)

### devlake-go

This contains the docker image (todo) and configuration scripts to setup DevLake to properly ingest the Dora metrics from a GitLab repo and group them according to Backstage groups.

[More details](devlake-go/README.md)

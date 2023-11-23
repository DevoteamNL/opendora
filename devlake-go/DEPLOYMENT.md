# OpenDORA API Deployment

This REST API provides an endpoint to return DORA metrics from DevLake. The [OpenAPI spec](https://github.com/DevoteamNL/opendora/blob/main/dora-api-mock/src/main/resources/openapi.yaml) outlines the path and query parameters to use the endpoint as well as the expected response.

# Docker Deployment

## Prerequisites

- GitHub account
- Docker
- DevLake: https://devlake.apache.org/docs/GettingStarted/DockerComposeSetup
- DevLake project with connections setup for DORA Metrics: https://devlake.apache.org/docs/DORA#how-to-implement-dora-metrics-with-apache-devlake

## Deployment

- First login to GitHub container registry using a personal access token: https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry#authenticating-with-a-personal-access-token-classic
- Pull the latest container:
```sh
docker pull ghcr.io/devoteamnl/opendora/opendora-api:latest
```
- Run the container with env arguments to connect to DevLake:
    - Replace env details with ones used to setup your DevLake instance
    - If running DevLake on Docker you may need to create a bridge network to connect the two
    - The server exposes port 10666 for the REST endpoints, make sure to bind/expose this: `-p 127.0.0.1:10666:10666/tcp`
```sh
docker run --name go-api -p 127.0.0.1:10666:10666/tcp -e DEVLAKE_DBUSER=merico -e DEVLAKE_DBPASS=merico -e DEVLAKE_DBADDRESS=localhost:3306 -e DEVLAKE_DBNAME=lake -d ghcr.io/devoteamnl/opendora/opendora-api:latest
```

# Kubernetes/Helm Deployment

## Prerequisites

- GitHub account
- Kubernetes
- Helm
- DevLake: https://devlake.apache.org/docs/GettingStarted/HelmSetup
- DevLake project with connections setup for DORA Metrics: https://devlake.apache.org/docs/DORA#how-to-implement-dora-metrics-with-apache-devlake

## Deployment

- First login to GitHub container registry using a personal access token (also make sure to save this token for a later step): https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry#authenticating-with-a-personal-access-token-classic
- Then, either using your docker config or manually through providing the PAT again, add the GHCR credentials as a secret to Kubernetes: https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/#registry-secret-existing-credentials
- Pull the chart:
```sh
helm pull oci://ghcr.io/devoteamnl/opendora/open-dora-helm-chart 
```
- Install the chart:

- If your DevLake or GHCR secret is setup with different details than the default, make sure to pass the appropriate arguments:
    - The MySQL database address and port used for DevLake: `--set openDora.devlakeDbAddress=devlake-mysql:3306`
    - The database username: `--set openDora.devlakeDbUser=merico`
    - The database password: `--set openDora.devlakeDbPass=merico`
    - The database name: `--set openDora.devlakeDbName=lake`
    - The secret name: `--set imagePullSecrets[0].name=regcred`
```sh
helm install opendora oci://ghcr.io/devoteamnl/opendora/open-dora-helm-chart --set openDora.devlakeDbAddress=devlake-mysql:3306 --set openDora.devlakeDbUser=merico --set openDora.devlakeDbPass=merico --set openDora.devlakeDbName=lake --set imagePullSecrets[0].name=regcred
```
# devlake-go

## group-sync

This will retrieve all Backstage entities of type `Group` and insert/remove them in the `teams` table of DevLake. This also updates the `parentId` for teams where the `childOf` and `parentOf` relationship points to an existing group/team in Backstage or DevLake.

**Note:** group names are will be saved with the correct case in DevLake but group relationships are case insensitive. If you have groups: `groupa` and `GroupA`, one of them parenting `groupb`, `groupb` will be parented to the first group found in DevLake (order not guaranteed). I recommend to just make group names unique ignoring case.

By default, the script will look for DevLake at http://localhost:4000/ and Backstage at http://localhost:7007/ but this can be changed using environment variables, respectively: `BACKSTAGE_URL` and `DEVLAKE_URL`. **Make sure to include the trailing forward slash `/`**

By default, the script will look for use DevLake admin credentials user:`devlake` and password:`merico`. This can be changed using environment variables, respectively: `DEVLAKE_ADMIN_USER` and `DEVLAKE_ADMIN_PASS`.

### To run the script

```shell
make run-sync
```

## api

This HTTP server provides an endpoint to return metrics from DevLake. The [OpenAPI spec](../dora-api-mock/src/main/resources/openapi.yaml) outlines the path and query parameters to use the endpoint as well as the expected response.

### To start the server

- Run the server, informing the variables for connecting to the database (your DevLake setup may be different)

```shell
DEVLAKE_DBUSER=merico \
    DEVLAKE_DBPASS=merico \
    DEVLAKE_DBADDRESS=localhost:3306 \
    DEVLAKE_DBNAME=lake \
    make run-api
```

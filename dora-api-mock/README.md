# Dora Mock API

This is a mock set up to represent the API used to retrieve Dora metrics from DevLake. It is a simple Spring Rest application.

### Prerequisites

- JDK 17
- Maven

### Running the application

You can run the application directly with the Spring Boot plugin. For that, just run the below:

```
./mvnw spring-boot:run
```

The API is now available at http://localhost:10666/ (or at the port as configured in `wiremock.server.port`), exposing the following endpoints:

- GET `/dora/api/openapi.yaml`: retrieves the DORA API OpenAPI 3.0 specification
- GET `/dora/api/metric`: retrieves the DORA metrics (only the `type` and `aggregation` parameters are supported as of now)

For instance, you can retrieve the Deployment Frequency Average datapoints, weekly aggregated, by running the following command:

```shell
curl -X GET "http://localhost:10666/dora/api/metric?type=df_average&aggregation=weekly" -H "accept: application/json"
```

### Changing or adding more Mock Endpoints/Data

The mocked data served by the API endpoints are stored as JSON files at [src/main/resources/stubs](src/main/resources/stubs).

The mock API is using [Spring Cloud Contract Wiremock Module](https://cloud.spring.io/spring-cloud-contract/reference/html/project-features.html#features-wiremock) to handle the HTTP request/response.

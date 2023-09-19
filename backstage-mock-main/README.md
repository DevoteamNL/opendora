# Backstage Dora Plugin API Mock

This is a mock set up to represent the API used to retrieve Dora metrics from DevLake. It is a simple Spring Rest application.

### Prerequisites:

- JDK 17
- Maven

### Running the application

You can run the application directly with the Spring Boot plugin. For that, just run the below:

```
./mvnw spring-boot:run
```

The API should now be available at http://localhost:8080/

### Changing or adding more Mock Endpoints/Data

Some of the mocked data served by the API endpoints are stored as JSON files at [src/main/resources/mock-data.json](src/main/resources/stubs).

The mock API is using [Spring Cloud Contract](https://spring.io/projects/spring-cloud-contract) to create the response stubs.
The contracts are defined at [src/main/resources/contracts](src/main/resources/contracts).

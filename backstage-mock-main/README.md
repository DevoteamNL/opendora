# Backstage Dora Plugin Backend Mock

This is a mock set up to represent the API used to retrieve Dora metrics from DevLake. It is a simple Spring Rest application.

### Prerequisites:
- JDK 11
- Maven

### Running the application

First build the project:
```
mvn clean
```
Then to run the application :
```
mvn spring-boot:run
```
The API should now be available at http://localhost:8080/

### Changing mock data

The mock data served by the api is stored in a JSON file at [src/main/resources/mock-data.json](src/main/resources/mock-data.json).

It contains a list of deployments with details about which projects they are from, the pass/fail status, dates, users involved with the deployment etc.
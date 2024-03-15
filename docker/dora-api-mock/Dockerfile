FROM eclipse-temurin:17-jdk-alpine as dev

WORKDIR /app

COPY ./dora-api-mock/.mvn/ .mvn/

COPY ./dora-api-mock/mvnw .

COPY ./dora-api-mock/pom.xml .

COPY ./dora-api-mock/src ./src

RUN ./mvnw install

EXPOSE 10666

ENTRYPOINT ["./mvnw", "spring-boot:run"]
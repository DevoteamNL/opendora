FROM eclipse-temurin:17-jdk-alpine as dev

WORKDIR /app

COPY .mvn/ .mvn/

COPY mvnw .

COPY pom.xml .

COPY src ./src

RUN ./mvnw install

EXPOSE 10666

ENTRYPOINT ["./mvnw", "spring-boot:run"]
FROM ubuntu:latest

LABEL Author="containers@devoteam.com"

COPY dist/open-dora-api /app/open-dora-api

WORKDIR /app

EXPOSE 10666

ENTRYPOINT ["/app/open-dora-api"]
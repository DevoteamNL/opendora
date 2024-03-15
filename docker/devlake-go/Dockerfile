ARG HOMEDIR=/opt/opendora

FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.21 AS builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

ARG HOMEDIR

ADD ./devlake-go /app/
WORKDIR /app/

RUN make build-static

WORKDIR ${HOMEDIR}

RUN cp /app/bin/opendora-api api

RUN echo "opendora:x:1000:opendora" >> /etc/group && \
    echo "opendora:x:1000:1000:opendora user:${HOMEDIR}:/sbin/nologin" >> /etc/passwd && \
    chown -R opendora:opendora ${HOMEDIR} && \
    chmod -R g+rw ${HOMEDIR} && \
    chmod +x api

FROM --platform=${TARGETPLATFORM:-linux/amd64} scratch
ARG HOMEDIR

LABEL Name="OpenDORA API" \
      Release=https://github.com/devoteamnl/opendora/releases \
      Url=https://github.com/devoteamnl/opendora \
      Help=https://github.com/devoteamnl/opendora/issues \
      org.opencontainers.image.source=https://github.com/devoteamnl/opendora

COPY --chown=1000:1000 --from=builder ${HOMEDIR} ${HOMEDIR}
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /usr/share/ca-certificates /usr/share/ca-certificates
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

WORKDIR ${HOMEDIR}

USER 1000

EXPOSE 10666

ENTRYPOINT [ "/opt/opendora/api" ]

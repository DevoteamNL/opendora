ARG HOMEDIR=/opt/opendora

FROM golang:1.21 AS builder
ARG HOMEDIR

ADD . /app/
WORKDIR /app/

RUN make static

WORKDIR ${HOMEDIR}

RUN cp /app/bin/opendora-api .

RUN echo "opendora:x:1000:opendora" >> /etc/group && \
    echo "opendora:x:1000:1000:opendora user:${HOMEDIR}:/sbin/nologin" >> /etc/passwd && \
    chown -R opendora:opendora ${HOMEDIR} && \
    chmod -R g+rw ${HOMEDIR} && \
    chmod +x opendora-api

FROM scratch
ARG HOMEDIR

LABEL Name="OpenDORA API" \
      Release=https://github.com/devoteamnl/opendora \
      Url=https://github.com/devoteamnl/opendora \
      Help=https://github.com/devoteamnl/opendora/issues

COPY --chown=1000:1000 --from=builder ${HOMEDIR} ${HOMEDIR}
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /usr/share/ca-certificates /usr/share/ca-certificates
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

WORKDIR ${HOMEDIR}
USER 1000
ENTRYPOINT [ "/opt/opendora/api" ]

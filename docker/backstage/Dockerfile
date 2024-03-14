FROM node:18-bookworm-slim AS node

FROM node as dev

ENV GIT_DISCOVERY_ACROSS_FILESYSTEM=1

RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt,sharing=locked \
    apt-get update && \
    apt-get install -y --no-install-recommends git python3 g++ build-essential && \
    yarn config set python /usr/bin/python3

WORKDIR /app

COPY .git .git

COPY --chown=node:node ./backstage-plugin ./backstage-plugin

COPY --chown=node:node --chmod=755 docker/backstage/entrypoint.sh /bin

RUN chown -R node:node /app

WORKDIR /app/backstage-plugin

# From here on we use the least-privileged `node` user to run the backend.
USER node

# This switches many Node.js dependencies to production mode.
# ENV NODE_ENV production

EXPOSE 3000

ENTRYPOINT [ "entrypoint.sh" ]

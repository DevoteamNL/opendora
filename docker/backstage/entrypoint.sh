#!/bin/bash

[ -z "$(ls -A /app)" ] && yes app | npx @backstage/create-app@latest --path /app

cp app-config.yaml /app
cp App.tsx /app/packages/app/src/App.tsx
cp Root.tsx /app/packages/app/src/components/Root/Root.tsx

cd /app && yarn --cwd packages/app add @devoteam-nl/open-dora-backstage-plugin && yarn install && yarn tsc && yarn --cwd packages/backend build && yarn dev
#!/bin/bash

[ -z "$(ls -A /app/node_modules)" ] && yarn install

yarn dev
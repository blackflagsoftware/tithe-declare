#!/bin/bash

npx nuxi generate --dotenv .env.prod
mkdir -p production
cp -r .output/public/ production
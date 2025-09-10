#!/bin/bash

npx nuxi generate
mkdir -p production
cp -r .output/public/ production
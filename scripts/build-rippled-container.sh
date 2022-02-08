#!/bin/bash

mkdir tmp && cd tmp/
git clone https://github.com/WietseWind/docker-rippled.git

cd docker-rippled
docker build --tag local-rippled:latest .

cd ../../
rm -rf tmp
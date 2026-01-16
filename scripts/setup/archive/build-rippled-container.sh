#!/bin/bash

mkdir tmp && cd tmp/ || exit
git clone https://github.com/WietseWind/docker-rippled.git

cd docker-rippled || exit
docker build --tag local-rippled:latest .

cd ../../
rm -rf tmp

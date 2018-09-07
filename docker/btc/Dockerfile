FROM ubuntu:18.04

RUN apt-get update && \
    apt-get install -y software-properties-common

RUN add-apt-repository ppa:bitcoin/bitcoin && \
    add-apt-repository -y ppa:ubuntu-toolchain-r/test

RUN apt-get update && \
    apt-get install -y libstdc++-7-dev bitcoind

EXPOSE 8332 8333 18332 18333

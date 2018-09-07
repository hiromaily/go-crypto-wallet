FROM ubuntu:18.04

RUN apt-get update && \
    apt-get install -y software-properties-common

RUN add-apt-repository ppa:bitcoin-abc/ppa && \
    add-apt-repository -y ppa:ubuntu-toolchain-r/test

RUN apt-get update && \
    apt-get install -y libstdc++-7-dev bitcoind

#RUN mkdir /root/.bitcoin
#COPY bitcoin.conf /root/.bitcoin/bitcoin.conf

EXPOSE 8332 8333 18332 18333

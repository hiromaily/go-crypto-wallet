FROM wallet-ubuntu:18.04
#FROM ubuntu:18.04

#RUN apt-get update && \
#    apt-get install -y software-properties-common

RUN add-apt-repository ppa:bitcoin/bitcoin && \
    add-apt-repository -y ppa:ubuntu-toolchain-r/test

RUN apt-get update && \
    apt-get install -y libstdc++-7-dev bitcoind

#RUN mkdir /root/.bitcoin
#COPY bitcoin.conf /root/.bitcoin/bitcoin.conf

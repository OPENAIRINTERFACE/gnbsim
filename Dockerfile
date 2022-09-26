# Copyright 2021-present Open Networking Foundation
#
# SPDX-License-Identifier: Apache-2.0
#

FROM golang:1.19.1-bullseye AS sim

LABEL maintainer="ONF <omec-dev@opennetworking.org>"

WORKDIR $GOPATH/src/gnbsim

COPY . $GOPATH/src/gnbsim 

RUN cd $GOPATH/src/gnbsim && \
    go build -mod=vendor && \
    ldd gnbsim

FROM debian:bullseye-slim AS gnbsim
ENV GOPATH=/go

RUN apt-get update && \
    apt-get -y install ethtool

WORKDIR /gnbsim/bin

COPY --from=sim $GOPATH/src/gnbsim/gnbsim /gnbsim/bin/
RUN ldconfig && \
    ldd /gnbsim/bin/

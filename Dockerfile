# Copyright 2021-present Open Networking Foundation
#
# SPDX-License-Identifier: Apache-2.0
#

FROM golang:1.19.1-bullseye AS sim

LABEL maintainer="ONF <omec-dev@opennetworking.org>"

RUN apt-get update && apt-get -y install vim ethtool
RUN cd $GOPATH/src && mkdir -p gnbsim
COPY . $GOPATH/src/gnbsim 

RUN cd $GOPATH/src/gnbsim && \
    go build -buildvcs=false -mod=vendor && \
    ldd gnbsim

FROM alpine:3.16 AS gnbsim
ENV GOPATH=/go
RUN apk update && apk add -U gcompat strace net-tools curl netcat-openbsd bind-tools bash tcpdump

RUN mkdir -p /gnbsim/bin

# Copy executable
WORKDIR /gnbsim/bin
COPY --from=sim $GOPATH/src/gnbsim/gnbsim /gnbsim/bin/
RUN ldconfig && \
    ldd /gnbsim/bin/gnbsim

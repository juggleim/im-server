FROM registry.cn-hangzhou.aliyuncs.com/juggleim/alpine:3.23

WORKDIR /opt

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add bash

ADD . /opt
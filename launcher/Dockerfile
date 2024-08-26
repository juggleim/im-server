FROM registry.cn-hangzhou.aliyuncs.com/juggleim/alpine:3.19

WORKDIR /opt

RUN mkdir /opt/conf

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add bash

ADD ./build/* /opt

RUN chmod +x /opt/run.sh

RUN chmod +x /opt/imserver

ENTRYPOINT /opt/run.sh

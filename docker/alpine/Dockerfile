FROM postgres:9.6.1-alpine
MAINTAINER <miaolizhao@126.com>

ENV SAMARITAN_VERSION 0.1.2
RUN apk add --no-cache --virtual .build-deps wget gnupg tar ca-certificates && \
    update-ca-certificates && \
    wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://raw.githubusercontent.com/sgerrand/alpine-pkg-glibc/master/sgerrand.rsa.pub && \
    wget -q https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.23-r3/glibc-2.23-r3.apk && \
    apk add glibc-2.23-r3.apk && \
    wget -q https://github.com/miaolz123/samaritan/releases/download/v${SAMARITAN_VERSION}/samaritan_linux_amd64.tar.gz && \
    tar -C /usr/src -xzf samaritan_linux_amd64.tar.gz && \
    echo $'#!/bin/sh\n\
nohup su-exec 999 postgres >/dev/null 2>&1 &\n\
sleep 10s\n\
cd /usr/src/samaritan_linux_amd64/\n\
./samaritan\n' >> /usr/src/samaritan_linux_amd64/cmd.sh && \
    chmod +x /usr/src/samaritan_linux_amd64/samaritan && \
    chmod +x /usr/src/samaritan_linux_amd64/cmd.sh && \
    rm -f samaritan_linux_amd64.tar.gz glibc-2.23-r3.apk /root/.gnupg && \
    apk del .build-deps && \
    apk add --no-cache --virtual .fetch-deps su-exec

EXPOSE 9876

VOLUME /var/lib/postgresql/data

COPY docker/config.ini /usr/src/samaritan_linux_amd64/custom/config.ini

CMD ["/usr/src/samaritan_linux_amd64/cmd.sh"]

FROM postgres:9.6.1
MAINTAINER <miaolizhao@126.com>

ENV SAMARITAN_VERSION 0.1.1
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates wget && rm -rf /var/lib/apt/lists/* && \
    wget -q https://github.com/miaolz123/samaritan/releases/download/v${SAMARITAN_VERSION}/samaritan_linux_amd64.tar.gz && \
    tar -C /usr/src -xzf samaritan_linux_amd64.tar.gz && \
    echo $'#!/bin/sh\n\
nohup gosu postgres postgres >/dev/null 2>&1 &\n\
sleep 10s\n\
cd /usr/src/samaritan_linux_amd64/\n\
./samaritan\n' >> /usr/src/samaritan_linux_amd64/cmd.sh && \
    chmod +x /usr/src/samaritan_linux_amd64/samaritan && \
    chmod +x /usr/src/samaritan_linux_amd64/cmd.sh && \
    rm -f samaritan_linux_amd64.tar.gz && \
    apt-get purge -y --auto-remove ca-certificates wget

EXPOSE 9876

VOLUME /var/lib/postgresql/data

COPY docker/config.ini /usr/src/samaritan_linux_amd64/custom/config.ini

CMD ["/usr/src/samaritan_linux_amd64/cmd.sh"]

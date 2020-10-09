FROM alpine

ADD redis /usr/local/soft/redis

WORKDIR /usr/local/soft/redis

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    mv sgerrand.rsa.pub /etc/apk/keys/sgerrand.rsa.pub && \
    apk --no-cache add tzdata expect glibc-2.29-r0.apk && \
    cp -rf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' >/etc/timezone && \
    echo "export LC_ALL=zh_CN.UTF-8" >> /etc/profile && source /etc/profile && \
    ln -s /usr/local/soft/redis/bin/redis-cli /usr/bin/redis-cli && \
    ln -s /usr/local/soft/redis/bin/redis-server /usr/bin/redis-server && \
    rm -f glibc-2.29-r0.apk && \
    apk del tzdata
EXPOSE 6379
EXPOSE 16379
ENTRYPOINT ["/bin/sh","/usr/local/soft/redis/start.sh"]


FROM alpine

COPY ./trx /trx

RUN apk add upx tzdata && \
    cp /usr/share/zoneinfo/Europe/Berlin /etc/localtime && \
    echo "Europe/Berlin" > /etc/timezone && \
    apk del tzdata

ENTRYPOINT ["/trx"]
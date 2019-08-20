FROM golang:1.12-alpine

ENV TIMEZONE="Europe/Berlin"

RUN apk update \
    && apk add \
        tzdata \
        git \
        build-base \
    && cp /usr/share/zoneinfo/${TIMEZONE} /etc/localtime \
    && echo "${TIMEZONE}" > /etc/timezone \
    && rm -rf /var/cache/apk/* \
    && go get github.com/cespare/reflex \
    && echo "-r '(\.go$|go\.mod)' -- go build" > /reflex.conf \
    && wget -O - -q https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | BINDIR=/usr/local/bin sh -s v1.17.1

WORKDIR /build
VOLUME ["/build"]
CMD ["reflex", "-v", "-c", "/reflex.conf"]

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
    && wget -O - -q https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | BINDIR=/usr/local/bin sh -s v1.17.1

RUN echo "-r '(\.go$|go\.mod)' -- golangci-lint run" > /reflex.conf \
    && echo "-r '(\.go$|go\.mod)' -- go build" >> /reflex.conf

WORKDIR /build
VOLUME ["/build"]
CMD ["reflex", "-v", "-c", "/reflex.conf"]

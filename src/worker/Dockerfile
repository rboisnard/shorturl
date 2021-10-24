ARG base_image=alpine:3.14

FROM ${base_image}

RUN  apk update \
  && apk add    \
    bash        \
    less        \
    go

COPY app/ /build

RUN  cd /build                  \
  && go build shorturl.go       \
  && mkdir -p /app              \
  && cp /build/shorturl /app/   \
  && cp /build/home.html /app/  \
  && rm -rf /build

COPY entrypoint.sh /

WORKDIR /app
ENTRYPOINT ["/entrypoint.sh"]

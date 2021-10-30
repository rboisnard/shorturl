ARG base_image=alpine:3.14

FROM ${base_image} as dev

RUN  apk update \
  && apk add    \
    bash        \
    less        \
    go

COPY app/ /build

RUN  cd /build  \
  && go get github.com/go-redis/redis/v8

FROM dev as builder

RUN  cd /build  \
  && go build shorturl.go

FROM dev

COPY --from=builder ["/build/shorturl", "/build/home.html", "/app/"]
COPY ["entrypoint.sh", "/"]

WORKDIR /app
ENTRYPOINT ["/entrypoint.sh"]

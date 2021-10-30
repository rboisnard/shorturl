ARG base_image=alpine:3.14

FROM ${base_image}

RUN  apk update \
  && apk add    \
    bash        \
    less        \
    redis

COPY entrypoint.sh /

ENTRYPOINT ["/entrypoint.sh"]

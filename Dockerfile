FROM alpine:latest

RUN apk add --update --no-cache \
    ca-certificates \
    fuse \
    openssh-client \
    bash

COPY ./dynamicflare /usr/bin/dynamicflare

CMD [ "/usr/bin/dynamicflare" ]
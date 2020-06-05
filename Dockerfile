FROM golang AS BUILD

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build

FROM alpine:latest

RUN apk add --update --no-cache \
    ca-certificates \
    fuse \
    openssh-client \
    bash

COPY --from=BUILD /app/dynamicflare /usr/bin/dynamicflare

CMD [ "/usr/bin/dynamicflare" ]
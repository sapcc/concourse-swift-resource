FROM alpine:3.6

RUN apk add --no-cache ca-certificates

ENV PATH /opt/resource:$PATH
COPY bin/ /opt/resource/

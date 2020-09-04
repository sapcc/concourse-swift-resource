FROM alpine:3.6
LABEL source_repository="https://github.com/sapcc/concourse-swift-resource"

RUN apk add --no-cache ca-certificates

ENV PATH /opt/resource:$PATH
COPY bin/ /opt/resource/

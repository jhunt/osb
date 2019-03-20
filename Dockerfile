FROM golang:1.11 AS build
WORKDIR /go/src/github.com/jhunt/osb
COPY . .
ARG VERSION
RUN make build CGO_ENABLED=0



FROM alpine:3.5
MAINTAINER James Hunt <james@niftylogic.com>

COPY --from=build /go/src/github.com/jhunt/osb/osb /usr/bin

ENTRYPOINT ["/usr/bin/osb"]
CMD        []

FROM golang:1.23-bookworm as builder
LABEL stage=builder

ENV REPO=thoas/picfit

ADD . /go/src/github.com/${REPO}

WORKDIR /go/src/github.com/${REPO}

RUN make docker-build-static && mv bin/picfit /picfit

###

FROM debian:buster-slim

RUN apt-get update && apt-get install -y ca-certificates

COPY --from=builder /picfit /picfit

CMD ["/picfit"]

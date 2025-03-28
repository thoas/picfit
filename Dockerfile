FROM golang:1.24-bookworm AS builder
LABEL stage=builder

ENV REPO=thoas/picfit

ADD . /go/src/github.com/${REPO}

WORKDIR /go/src/github.com/${REPO}

RUN make docker-build-static && mv bin/picfit /picfit

###

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates

COPY --from=builder /picfit /picfit

CMD ["/picfit"]

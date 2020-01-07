FROM ubuntu:18.04

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

ADD bin/picfit /picfit

CMD ["/picfit"]

FROM golang:1.25-trixie AS builder
LABEL stage=builder

ENV REPO=thoas/picfit

# libvips >= 8.16 is required by github.com/cshum/vipsgen/vips816
RUN apt-get update \
    && apt-get install -y --no-install-recommends libvips-dev \
    && rm -rf /var/lib/apt/lists/*

ADD . /go/src/github.com/${REPO}

WORKDIR /go/src/github.com/${REPO}

RUN make docker-build-static && mv bin/picfit /picfit

###

FROM debian:trixie-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates libvips42t64 gifsicle \
    && rm -rf /var/lib/apt/lists/*

# Limits glibc malloc arenas, libvips threads otherwise fragment the memory
ENV MALLOC_ARENA_MAX=2
# Disables the libvips loaders not designed for untrusted input (magick, pdf, openslide...)
ENV VIPS_BLOCK_UNTRUSTED=1

COPY --from=builder /picfit /picfit

CMD ["/picfit"]

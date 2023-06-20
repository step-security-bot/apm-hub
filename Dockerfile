FROM golang:1.20 as builder
WORKDIR /app
ARG VERSION
COPY go.mod /app/go.mod
COPY go.sum /app/go.sum
RUN go mod download
COPY ./ ./
WORKDIR /app
RUN go version
RUN make build

FROM ubuntu:jammy
WORKDIR /app

# install CA certificates
RUN apt-get update && \
  apt-get install -y ca-certificates && \
  rm -Rf /var/lib/apt/lists/*  && \
  rm -Rf /usr/share/doc && rm -Rf /usr/share/man  && \
  apt-get clean

COPY --from=builder /app/.bin/apm-hub /app
ENV ASSUME_NO_MOVING_GC_UNSAFE_RISK_IT_WITH=go1.20
ENTRYPOINT ["/app/apm-hub"]

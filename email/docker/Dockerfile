FROM golang:1.18-buster as builder

WORKDIR /app

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY cmd cmd
COPY internal internal

COPY Makefile .

RUN make build

FROM debian:buster

WORKDIR /app
USER nobody

COPY --from=builder /app/worker /app/worker

ENTRYPOINT ["/app/worker"]


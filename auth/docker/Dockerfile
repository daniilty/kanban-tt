FROM golang:1.17.2-buster as builder

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

COPY --from=builder /app/server /app/server

EXPOSE 8000

ENTRYPOINT ["/app/server"]


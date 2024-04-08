FROM golang:1.22-alpine3.19 AS builder

WORKDIR /usr/local/src
RUN apk --no-cache add bash git make gcc gettext musl-dev

# dependencies
COPY ["app/go.mod", "app/go.sum", "./"]
COPY app ./
RUN go mod download

# build
RUN go build -o ./bin/app main.go


FROM alpine:3.19 AS runner

COPY --from=builder /usr/local/src/bin/app /
COPY rconfig/config.yml /config/config.yml
CMD ["/app"]
FROM golang:1.23.3-alpine3.20 AS builder

WORKDIR /usr/local/src
RUN apk --no-cache add bash git make gcc gettext musl-dev

# dependencies
COPY ["app/go.mod", "app/go.sum", "./"]
RUN go mod download

# build
COPY app ./
RUN go build -o ./bin/balance main.go


FROM alpine:3.20 AS runner

RUN apk add --no-cache tzdata
ENV TZ=America/Winnipeg
COPY --from=builder /usr/local/src/bin/balance /
CMD ["/balance"]
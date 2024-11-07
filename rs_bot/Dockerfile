FROM golang:1.23.1-alpine3.20 AS builder

WORKDIR /usr/local/src
RUN apk --no-cache add bash git make gcc gettext musl-dev

# dependencies
COPY ["app/go.mod", "app/go.sum", "./"]
COPY app ./
RUN go mod download

# build
RUN go build -o ./bin/rs_bot main.go


FROM alpine:3.20 AS runner

RUN apk add --no-cache tzdata
ENV TZ=America/Winnipeg
COPY --from=builder /usr/local/src/bin/rs_bot /
CMD ["/rs_bot"]
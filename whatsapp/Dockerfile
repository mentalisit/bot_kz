FROM golang:1.25.3-alpine3.22 AS builder

WORKDIR /usr/local/src
RUN apk --no-cache add bash git make gcc gettext musl-dev

# dependencies
COPY ["app/go.mod", "app/go.sum", "./"]
COPY app ./
RUN go mod download

# build
RUN go build -o ./bin/whatsapp main.go


FROM alpine:3.22 AS runner

RUN apk add --no-cache tzdata
ENV TZ=America/Winnipeg
COPY --from=builder /usr/local/src/bin/whatsapp /
CMD ["/whatsapp"]
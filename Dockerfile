FROM golang:1.19-alpine AS build

RUN apk update  \
    && apk add --update upx gcc git musl-dev

RUN adduser -DHg scratchuser scratchuser

WORKDIR /app

ENV CGO_ENABLED=1
ENV GOOS=linux

# Downloading dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

ADD . ./

RUN go build -a -ldflags='-linkmode external -extldflags "-static" -w -s' -o traefik-u2f .

# Compress executable
RUN upx traefik-u2f
RUN chmod 555 traefik-u2f

FROM scratch

WORKDIR /
COPY --from=build /app/traefik-u2f /
COPY --from=build /app/config.yaml /
COPY --from=build /etc/passwd /etc/passwd
USER scratchuser
CMD ["/traefik-u2f"]

FROM golang:1.19-alpine AS build

RUN apk update && apk add --update upx

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux

# Downloading dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

ADD . ./

#RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Building the app
RUN go build -o traefik-u2f .

# Compress executable
RUN upx traefik-u2f

FROM scratch
WORKDIR /
COPY --from=build /app/traefik-u2f .
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
USER 405
CMD ["./traefik-u2f"]

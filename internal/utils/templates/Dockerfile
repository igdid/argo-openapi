FROM golang:1.26-alpine AS builder

ARG UID=10005
ARG GID=10005

RUN addgroup -g ${GID} sender && \
    adduser -D -H -u ${UID} -G sender sender && \
    apk add --no-cache ca-certificates

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
      -trimpath \
      -ldflags="-s -w" \
      -o /out/{{ .AppName }} ./cmd

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
USER sender:sender
EXPOSE 8080
EXPOSE 8081

COPY --from=builder /out/{{ .AppName }} /usr/bin/{{ .AppName }}
ENTRYPOINT ["/usr/bin/{{ .AppName }}"]

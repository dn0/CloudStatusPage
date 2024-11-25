# Usage: TARGET=mon-app-name envsubst '$TARGET' < Dockerfile
FROM golang:alpine AS builder
ARG VERSION
RUN apk update && apk add --no-cache make ca-certificates
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . ./
RUN GOARCH=amd64 VERSION=$VERSION make bin/${TARGET}

FROM scratch AS final
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder --chmod=755 /app/bin/${TARGET}.amd64 /${TARGET}
COPY srv/${TARGET} /srv/${TARGET}
CMD ["/${TARGET}"]

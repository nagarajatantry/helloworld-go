FROM alpine:latest as certs
RUN apk --update add ca-certificates go

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o main .

FROM busybox
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=certs /app/main /main
COPY cert.pem /cert.pem
COPY key.pem /key.pem
ENTRYPOINT ["/main"]
EXPOSE 9098
ARG VERSION=latest

FROM golang:1.16 as builder
ARG VERSION

WORKDIR /go/src

COPY . /go/src

RUN go get
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -o emr-idle-check -ldflags "-s -w -X main.version=${VERSION}"
RUN md5sum emr-idle-check

# Create image from scratch
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/emr-idle-check /emr-idle-check
COPY --from=builder /tmp /tmp

ENTRYPOINT [ "/emr-idle-check" ]

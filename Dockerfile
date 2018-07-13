FROM golang:alpine AS builder
MAINTAINER "Cuong Manh Le <cuong.manhle.vn@gmail.com>"

RUN apk update && \
    apk add git build-base && \
    rm -rf /var/cache/apk/* && \
    mkdir -p "$GOPATH/src/github.com/bitnami-labs/kubewatch"

ADD . "$GOPATH/src/github.com/bitnami-labs/kubewatch"

RUN cd "$GOPATH/src/github.com/bitnami-labs/kubewatch" && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a --installsuffix cgo --ldflags="-s" -o /kubewatch

FROM alpine:3.4
RUN apk add --update ca-certificates

COPY --from=builder /kubewatch /bin/kubewatch

# ENV KW_GOOGLECHAT_URL=https://chat.googleapis.com/v1/spaces/AAAAosXkqiA/messages?key=AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI&token=soK8YP_o6jn_ShJi_QErkn95XVIWOc7wletgqtHUNnY%3Dhttps://chat.googleapis.com/v1/spaces/AAAAosXkqiA/messages?key=AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI&token=soK8YP_o6jn_ShJi_QErkn95XVIWOc7wletgqtHUNnY%3D

ENTRYPOINT ["/bin/kubewatch"]

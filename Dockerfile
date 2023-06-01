FROM golang:1.17.2-buster AS build

WORKDIR /app

COPY go-plugins/ ./

RUN cd cmd/go-runner && CGO_ENABLED=0 GOOS=linux go build -o /go-plugins .

FROM apache/apisix:2.15.0-alpine

USER root

RUN mkdir /runner && chmod +x /runner

COPY --from=build /go-plugins /runner/

WORKDIR /usr/local/apisix

EXPOSE 9080/tcp

EXPOSE 9443/tcp

CMD ["/bin/sh", "-c", "/usr/bin/apisix init && /usr/bin/apisix init_etcd && /usr/local/openresty/bin/openresty -p /usr/local/apisix -g 'daemon off;'"]

STOPSIGNAL SIGQUIT
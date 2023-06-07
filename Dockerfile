FROM golang:1.17.2-buster AS build

WORKDIR /app

COPY go-plugins/ ./

RUN cd cmd/go-runner && CGO_ENABLED=0 GOOS=linux go build -o /go-plugins .

FROM apache/apisix:3.2.0-debian

RUN mkdir /runner && chmod +x /runner

COPY --from=build /go-plugins /runner/

WORKDIR /usr/local/apisix

USER apisix

EXPOSE 9080/tcp

EXPOSE 9443/tcp

COPY ./docker-entrypoint.sh /docker-entrypoint.sh # buildkit

ENTRYPOINT ["/docker-entrypoint.sh"]

CMD ["docker-start"]

STOPSIGNAL SIGQUIT

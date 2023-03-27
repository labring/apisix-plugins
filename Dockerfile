FROM golang:1.17.2-buster AS build
WORKDIR /app
COPY go-plugins/ ./
RUN cd cmd/go-runner && CGO_ENABLED=0 GOOS=linux go build -o /go-plugins .

FROM apache/apisix:3.2.0-debian
USER root
RUN mkdir /runner && chown apisix:apisix /runner
COPY --from=build /go-plugins /runner/
USER apisix
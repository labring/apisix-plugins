FROM golang:1.17.2-buster AS build
WORKDIR /app
COPY go-runner/ ./
RUN cd cmd/go-runner && CGO_ENABLED=0 GOOS=linux go build -o /go-runner .

FROM apache/apisix:3.2.0-debian
USER root
RUN mkdir /runner && chown apisix:apisix /runner
COPY --from=build /go-runner /runner/
USER apisix
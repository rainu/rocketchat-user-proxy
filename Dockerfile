FROM golang:1.10 as buildContainer

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOPATH=/

COPY . /src/rocketchat-user-proxy
WORKDIR /src/rocketchat-user-proxy

RUN go get ./... &&\
    go build -ldflags -s -a -installsuffix cgo


FROM alpine

COPY --from=buildContainer /src/rocketchat-user-proxy/rocketchat-user-proxy /rocketchat-user-proxy

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

USER 10000:10000
EXPOSE 8080

ENTRYPOINT ["/rocketchat-user-proxy"]
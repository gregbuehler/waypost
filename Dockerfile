FROM golang:latest as builder

WORKDIR /go/src/github.com/gregbuehler/waypost
COPY . .
RUN cd /go/src/github.com/gregbuehler/waypost && \
    CGO_ENABLED=0 GOOS=linux \
    make compile

##
FROM alpine:latest

RUN apk update && apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder go/src/github.com/gregbuehler/waypost/bin/waypost .
COPY etc /etc/

EXPOSE 53
EXPOSE 5554

ENTRYPOINT [ "./waypost" ]
CMD [ "serve"]
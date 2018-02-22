FROM alpine

RUN apk add --update ca-certificates

COPY go-desync /

EXPOSE 8080

CMD ["/go-desync"]

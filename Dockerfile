FROM docker.io/golang:1.17-alpine

EXPOSE 8888

WORKDIR /go/src/go-torrent
COPY . .

RUN go install go-torrent

CMD /go/bin/go-torrent
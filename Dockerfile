FROM golang:1.17.5 AS builder

#ENV GOPROXY https://goproxy.cn,direct

WORKDIR /www

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o filebeat

FROM alpine

LABEL author=github.com/wolanx
ENV TZ utc-8

WORKDIR /usr/share/filebeat

COPY --from=builder /www/filebeat .

ENTRYPOINT ["./filebeat"]

# docker build -f Dockerfile -t zx5435/wolan:logging .
# docker run --restart=unless-stopped --name wlog -d -p 20100:20100 zx5435/wolan:logging

FROM golang AS builder

ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on

RUN go version
RUN go install github.com/cortesi/modd/cmd/modd@latest

CMD [ "modd" ]
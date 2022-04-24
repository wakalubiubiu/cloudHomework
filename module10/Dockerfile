FROM golang:1.17.8-alpine3.15 as build
RUN mkdir /go/src/homework/
COPY main.go metrics.go /go/src/homework/
RUN export GOPATH=/go 
WORKDIR /go/src/homework
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go mod init &&go env -w GOPROXY=https://goproxy.cn,direct && go mod tidy && go build -o /go/bin/homework/httpserver
EXPOSE 8080

FROM busybox
RUN mkdir /apps/
COPY --from=build /go/bin/homework/httpserver /apps/
WORKDIR /apps
RUN pwd && ls -l
ENTRYPOINT /apps/httpserver -log_dir="./"

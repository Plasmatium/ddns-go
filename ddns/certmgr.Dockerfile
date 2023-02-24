FROM golang:alpine AS BUILD
WORKDIR /work_dir
ADD . /work_dir/
RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    CGO_ENABLED=0 go build -o dist/certmgr bins/certbot-dns-aliyun/main.go

FROM alpine
RUN sed -i 's/dl-cdn.alpinelinux.org/mirror.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && \
    apk add --no-cache certbot
WORKDIR /work_dir
COPY --from=BUILD /work_dir/dist/certmgr .
COPY ./do_cert.sh /work_dir/
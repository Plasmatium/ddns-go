FROM golang:alpine AS BUILD
WORKDIR /work_dir
ADD . /work_dir/
RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    CGO_ENABLED=0 go build -o dist/ddns

FROM alpine
WORKDIR /work_dir
COPY --from=BUILD /work_dir/dist/ddns .
ENTRYPOINT [ "./ddns" ]

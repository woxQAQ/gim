FROM golang:alpine AS builder
LABEL authors="shirakami_yuki"
ENV GOPROXY=https://goproxy.cn,direct \
GO111MODULE=on \
CGO_ENABLE=0 \
GOOS=linux \
GOARCH=arm64
WORKDIR /gIM
COPY . /gIM
RUN go build -o ./cmd/app ./build
FROM scratch
COPY --from=builder /gIM/cmd/app /
EXPOSE 8000
ENTRYPOINT ["/app"]
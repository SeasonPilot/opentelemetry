# Build the app binary
FROM golang:1.18 as builder

WORKDIR /workspace
# Copy the go source
COPY . .

ENV GOPROXY="https://goproxy.cn"
# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -mod=vendor -o app


FROM alpine
WORKDIR /
COPY --from=builder /workspace/app .
ENTRYPOINT ["/app"]
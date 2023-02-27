# Build the app binary
FROM golang as builder

WORKDIR /workspace
# Copy the go source
COPY . .
# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build  -o app


FROM alpine
WORKDIR /
COPY --from=builder /workspace/app .
ENTRYPOINT ["/app"]
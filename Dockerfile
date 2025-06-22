FROM golang:1.23 AS builder
WORKDIR /go/src/github.com/robrohan/laffaire
COPY . .
RUN make build

FROM builder
WORKDIR /root/
COPY --from=builder /go/src/github.com/robrohan/laffaire/build/ ./
RUN ls -alFh
CMD ["./server"]

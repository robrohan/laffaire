# Note: alpine linux does not have the libraris to CGO_ENABLED to build 
# sqlite3. If you want to build that you need to change to a linux image
FROM golang:alpine as builder
RUN apk --no-cache add gcc g++ make ca-certificates git
WORKDIR /go/src/github.com/robrohan/go-web-template
COPY . .
RUN make build

FROM golang:alpine
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/robrohan/go-web-template/build/ ./
RUN ls -alFh
CMD ["./server"]

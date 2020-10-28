FROM golang:alpine as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main github.com/codemicro/contactFormProcessor/cmd/contactFormProcessor
FROM alpine
COPY --from=builder /build/main /
WORKDIR /run
CMD ["../main"]
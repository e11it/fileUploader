FROM golang:latest AS builder

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 && chmod +x /usr/local/bin/dep

RUN mkdir -p /go/src/github.com/e11it/fileUploader
WORKDIR /go/src/github.com/e11it/fileUploader

ADD . ./

RUN dep ensure -vendor-only
# install the dependencies without checking for go code
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /fileUploader .

FROM scratch
COPY --from=builder fileUploader /app/
EXPOSE 8080
VOLUME /upload
WORKDIR /upload
CMD ["/app/fileUploader"]
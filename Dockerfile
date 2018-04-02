FROM golang:latest
ADD . /go/src/url-shortener
WORKDIR /go/src/url-shortener
RUN go get -v
RUN go install url-shortener

FROM golang:latest
COPY --from=0 /go/bin/url-shortener .
ENV PORT 8080
CMD ["./url-shortener"]

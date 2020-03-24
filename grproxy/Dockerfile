FROM golang
COPY ./src/grproxy .
COPY github.com .
RUN go get github.com/talbright/go-zookeeper/zk
RUN CGO_ENABLED=0 go build -o /grproxy
RUN echo "nobody:x:65534:65534:Nobody:/:" > /etc_passwd



FROM scratch
COPY --from=0 /grproxy /grproxy
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /etc_passwd /etc/passwd
USER nobody
CMD ["./grproxy"]

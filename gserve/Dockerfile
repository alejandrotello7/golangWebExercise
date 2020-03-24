FROM golang
COPY ./src/gserve .
COPY github.com .
RUN go get github.com/talbright/go-zookeeper/zk
RUN CGO_ENABLED=0 go build -o /gserve
COPY page.html ./gserve
RUN echo "nobody:x:65534:65534:Nobody:/:" > /etc_passwd

FROM scratch
COPY page.html .
COPY --from=0 /gserve /gserve
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /etc_passwd /etc/passwd
USER nobody
CMD ["./gserve"]

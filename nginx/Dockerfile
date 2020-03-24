FROM alpine:latest
RUN apk upgrade
RUN apk update
RUN apk add nginx
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 80
CMD ["nginx"]
RUN mkdir -p /run/nginx



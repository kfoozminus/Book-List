FROM busybox:glibc
COPY ./web /web
EXPOSE 8080
CMD ["./web"]

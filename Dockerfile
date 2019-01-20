FROM golang:1.11-alpine
WORKDIR /go/src/github.com/fbcbarbosa/drone-ignore-config/
RUN apk add -U --no-cache ca-certificates
ADD . .
RUN GOOS=linux CGO_ENABLED=0 go build -o /bin/changeset \
    github.com/fbcbarbosa/drone-ignore-config/cmd/drone-ignore-config

FROM scratch
COPY --from=0 /bin/changeset /bin/changeset
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 3000
ENTRYPOINT ["/bin/changeset"]
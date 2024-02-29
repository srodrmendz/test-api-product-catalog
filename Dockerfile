# BUILD IMAGE
FROM golang:1.18 as builder

WORKDIR /go/src/app/

COPY . .

ARG version
ARG githubToken

RUN git config --global url."${githubToken}".insteadOf "https://github.com/"
RUN pkgPath="$(go list -m)/conf.Version=${version}" && GOPRIVATE="github.com/basset-la" CGO_ENABLED=0 GOOS=linux go build -ldflags="-X $pkgPath" -a -installsuffix cgo -o app

# PROD/DEV IMAGE
FROM alpine:latest


WORKDIR /go/src/app/

COPY --from=builder /go/src/app/app .
COPY --from=builder /go/pkg/mod /go/pkg/mod
COPY --from=builder /go/src/app/config/ config/

CMD ./app -e $env

EXPOSE 8080

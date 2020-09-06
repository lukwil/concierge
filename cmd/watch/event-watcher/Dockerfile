FROM golang:1.14-alpine as build

RUN mkdir -p /go/src/handler
WORKDIR /go/src/handler
COPY . .

ENV CGO_ENABLED=0


RUN go test ./... -cover

RUN go build --ldflags "-s -w" -a -installsuffix cgo -o handler .

FROM alpine:3.11
# Add non root user and certs
RUN apk --no-cache add ca-certificates \
    && addgroup -S app && adduser -S -g app app \
    && mkdir -p /home/app \
    && chown app /home/app

WORKDIR /home/app

COPY --from=build /go/src/handler/handler .

RUN chown -R app /home/app

RUN touch /tmp/.lock #  Write a health check for OpenFaaS here or in the HTTP server start-up

USER app

EXPOSE 8080

HEALTHCHECK --interval=3s CMD [ -e /tmp/.lock ] || exit 1

CMD ["/home/app/handler"]

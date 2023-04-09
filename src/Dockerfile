FROM golang:alpine as build

RUN apk --no-cache add make
RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/app

# Now copy it into our base image.
FROM gcr.io/distroless/static-debian11 as cloudrun
COPY --from=build /go/bin/app /
CMD ["/app", "-cloudrun"]

FROM gcr.io/distroless/static-debian11 as default
COPY --from=build /go/bin/app /
CMD ["/app"]
FROM golang:1.12.9 AS builder
WORKDIR /go/src/app
COPY . .
RUN go get github.com/gobuffalo/packr/v2
RUN go get -u github.com/gobuffalo/packr/v2/packr2
RUN go get github.com/victorspringer/http-cache
RUN packr2
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o conform

FROM scratch
WORKDIR /go/src/app
COPY --from=builder /go/src/app/conform .
COPY --from=builder /go/src/app/assets ./assets
EXPOSE 8812
CMD ["./conform"]

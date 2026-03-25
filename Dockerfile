FROM golang:1.21-alpine

WORKDIR /app
COPY go.mod ./
RUN go mod download && go mod verify

COPY . .
RUN go build

FROM scratch
COPY --from=0 /app/bannermd /app/bannermd
COPY ./banner.tmpl ./banner.tmpl

CMD ["/app/bannermd"]

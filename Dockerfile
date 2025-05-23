FROM golang:1.24-alpine AS build

WORKDIR /work

COPY cmd cmd
COPY internal internal
COPY pkg pkg
COPY go.mod go.mod
COPY go.sum go.sum

ENV CGO_ENABLED=0
RUN go build -o /tmp ./cmd/server

FROM scratch AS deploy
COPY --from=build /tmp/server  /server  
EXPOSE 8080

CMD ["/server"]
FROM golang:latest AS build

WORKDIR /src/
COPY go.mod go.sum main.go /src/
COPY service/* /src/service/
RUN CGO_ENABLED=0 go build -o /bin/service

FROM scratch
COPY --from=build /bin/service /bin/service
ENTRYPOINT ["/bin/service"]
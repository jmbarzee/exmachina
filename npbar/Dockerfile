FROM golang:latest AS build

WORKDIR /src/
COPY go.mod go.sum main.go /src/
COPY grpc/* /src/grpc
COPY service/* /src/service
ENV CGO_ENABLED=0
RUN go build -o /bin/service

FROM scratch
COPY --from=build /bin/service /bin/service
ENTRYPOINT ["/bin/service"]
FROM golang:latest AS build

WORKDIR /src/
COPY go.mod go.sum main.go /src/
COPY service/ /src/service/

# Install dependecies for github.com/faiface/pixel
RUN apt-get -y update
RUN apt-get -y install xorg-dev libgl1-mesa-dev

ENV CGO_ENABLED=1
RUN go build -o /bin/service



FROM scratch

COPY --from=build /bin/service /bin/service
ENTRYPOINT ["/bin/service"]
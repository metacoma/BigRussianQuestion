FROM golang:1.9.4-alpine
RUN apk add --update git bind-tools curl jq imagemagick sqlite gcc sqlite-dev musl-dev coreutils
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build
ADD ./images /images
ADD ./image_generator /image_generator
CMD ["app"]

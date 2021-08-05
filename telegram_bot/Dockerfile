FROM golang:1.9.4-alpine
RUN apk add --update git bind-tools curl jq imagemagick sqlite gcc sqlite-dev musl-dev coreutils python3
WORKDIR /go/src/app
COPY . .
RUN pip3 install -r requirements.txt
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build
ADD ./images /images
ADD ./image_generator /image_generator
ADD ./brq_upload_vk.py /usr/local/bin/brq_upload_vk.py
CMD ["app"]

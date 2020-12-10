FROM golang

RUN mkdir -p /go/videosapp

RUN go get github.com/gorilla/mux

RUN go get github.com/lib/pq

RUN go get github.com/lib/pq

RUN go get github.com/subosito/gotenv

ADD .env /go/videosapp/

ADD *.go /go/videosapp/

WORKDIR /go/videosapp

RUN go build .

EXPOSE 80

CMD [ "./videosapp" ]
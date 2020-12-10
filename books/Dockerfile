FROM golang

RUN mkdir -p /go/booksapp

RUN go get github.com/gorilla/mux

RUN go get github.com/lib/pq

RUN go get github.com/lib/pq

RUN go get github.com/subosito/gotenv

ADD .env /go/booksapp/

ADD *.go /go/booksapp/

WORKDIR /go/booksapp

RUN go build .

EXPOSE 80

CMD [ "./booksapp" ]
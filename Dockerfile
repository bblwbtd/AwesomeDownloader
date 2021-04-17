FROM golang

WORKDIR /go/src/app
COPY . .

RUN go mod vendor

EXPOSE 1234

CMD /bin/bash -c 'go run .'

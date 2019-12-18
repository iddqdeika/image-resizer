#версия образа
FROM golang:latest

#создаем папку
RUN mkdir -p /go/src/image-resizer

#идем в папку
WORKDIR /go/src/image-resizer

#копируем файлы
COPY . /go/src/image-resizer

#зависимости
RUN go get -d -v ./...

#билдим
RUN go build -o main

#запускаем
CMD ["/go/src/image-resizer/main"]
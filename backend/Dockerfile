FROM golang:latest

LABEL base.name="backend"


WORKDIR /app

COPY . .


RUN go build -o server .


EXPOSE 8000


ENTRYPOINT ["./server"]
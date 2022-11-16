FROM golang:1.18-bullseye

RUN mkdir /App
WORKDIR /App
RUN cd /App
COPY . .
RUN go build .

CMD ["./Go-Filter-Bot"]
FROM golang:1.18-bullseye

RUN mkdir /App
WORKDIR /App
RUN cd /App
COPY . .

CMD ["go","run","."]

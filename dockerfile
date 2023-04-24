FROM golang:1.19

WORKDIR /app

COPY . .

ENV MONGO_DB_URI=mongodb+srv://airscrum_admin:kSF34_LdAD6f-kB@cluster0.vya2xvh.mongodb.net/?retryWrites=true&w=majority \
    RABBITMQ_URI=amqp://admin:admin@localhost:5672/ \
    PORT=8002

RUN go get -d -v ./...

RUN go install -v ./...

CMD ["go", "run", "main.go"]

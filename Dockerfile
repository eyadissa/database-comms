FROM golang:1.17-alpine

#FOR Cloud run
EXPOSE 8080

WORKDIR /app
COPY . ./
COPY go.mod go.sum ./
RUN go mod tidy
RUN go build -o /main
RUN chmod +x /start.sh
CMD [ "./start.sh" ]

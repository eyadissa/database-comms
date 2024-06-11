FROM golang:1.17-alpine

WORKDIR /app
COPY . ./
COPY go.mod go.sum ./
RUN go mod tidy
RUN go build -o /main
CMD [ "./start.sh" ]

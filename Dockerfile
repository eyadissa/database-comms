FROM golang:1.17-alpine
ENV PORT 2345
ENV HOSTDIR 0.0.0.0

EXPOSE 2345
WORKDIR ./
#COPY go.mod go.sum ./
RUN go mod tidy
RUN go build -o /main
CMD [ "/main" ]

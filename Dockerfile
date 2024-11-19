FROM golang:1.22.2-alpine

#set workdir
WORKDIR /app

#copy go.mod and go.sum to workdir
COPY go.mod go.sum ./

#install dependencies
RUN go mod download

# copy soure code to workdir
COPY . .

#build go application
RUN go build -o main .

#expose port
EXPOSE 5000

#set entry point to container
CMD ["./main"]



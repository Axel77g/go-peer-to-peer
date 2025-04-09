FROM golang:latest

# Set the working directory
WORKDIR /app
# Copy the Go modules manifests
# Copy the source code into the container
COPY . .

EXPOSE 9999/udp

CMD [ "sh", "-c", "go run main.go" ]

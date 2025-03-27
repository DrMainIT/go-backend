# Start with the official Go image
FROM golang:1.23

# # Set the working directory inside the container
WORKDIR /app

# # Copy go.mod and go.sum files
COPY go.mod go.sum ./

# # Download dependencies
RUN go mod download

# # Copy the rest of the application code
COPY . .

ENV POSTGRES_USER=postgresql://postgres:postgres@localhost:5432/postgres
ENV FRONTEND_URL=http://127.0.0.1:3000
# # Build the Go application
RUN go build -o main .

# # Expose the port the app runs on
EXPOSE 8080

# # Command to run the application
CMD ["./main"]
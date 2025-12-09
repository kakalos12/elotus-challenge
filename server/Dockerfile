FROM golang:latest

WORKDIR /app

# Install air for hot reloading
RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# We don't need to build the binary here for development with Air,
# but keeping the build step can be useful for checking compilation errors.
# RUN go build -o server .

EXPOSE 8080

# Run Air
CMD ["air"]
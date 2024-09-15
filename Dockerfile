# Stage 1: Build the Go application
FROM golang:alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY server .
COPY sites ./sites
RUN go build -o main .

# Stage 2: Set up the Nginx server and Go application
FROM nginx:alpine

# Copy the Nginx configuration file
COPY nginx.conf /etc/nginx/nginx.conf

# Copy the built Go application and sites from the builder stage
COPY --from=builder /app/main /app/
COPY --from=builder /app/sites /sites

# Expose the port Nginx will be running on
EXPOSE 80

# Add a script to run both Nginx and the Go application
COPY start.sh /start.sh
RUN chmod +x /start.sh

CMD ["/start.sh"]

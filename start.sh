#!/bin/sh

# Start the Go application in the background
/app/main &

# Start Nginx in the foreground
nginx -g 'daemon off;'

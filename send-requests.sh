#!/bin/bash

# Ensure the script has exactly three arguments
if [ "$#" -ne 3 ]; then
    echo "Usage: $0 <initial_requests> <file_path> <domain>"
    exit 1
fi

# Read command-line arguments
initial_requests=$1
file_path=$2
domain=$3

# Validate the number of requests
if ! [[ "$initial_requests" =~ ^[0-9]+$ ]] || [ "$initial_requests" -le 0 ]; then
    echo "Error: Initial number of requests must be a positive integer."
    exit 1
fi

# Validate the file path
if [ ! -f "$file_path" ]; then
    echo "Error: File not found at path $file_path."
    exit 1
fi

# Time to wait between each step (in seconds)
wait_time=5
division_factor=2

# Record the start time
start_time=$(date +%s%3N)

# Function to send requests
send_requests() {
    local request_count=$1
    for ((i=1; i<=request_count; i++)); do
        curl -s -o /dev/null -F "file=@$file_path" http://$domain/upload &
    done
    wait
}

# Loop to gradually decrease the number of requests by division
current_requests=$initial_requests
while [ $current_requests -ge 1 ]; do
    echo "Sending $current_requests requests..."
    send_requests $current_requests

    # Divide the number of requests
    current_requests=$((current_requests / division_factor))
    if [ $current_requests -lt 1 ]; then
        current_requests=0
    fi

    # Wait before sending the next batch
    sleep $wait_time
done

# Record the end time
end_time=$(date +%s%3N)

# Calculate elapsed time in milliseconds
elapsed_time=$((end_time - start_time))

echo "All requests have been sent."
echo "Elapsed time: ${elapsed_time} ms"

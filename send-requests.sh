#!/bin/bash

# Ensure the script has at least two arguments
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <number_of_requests> <file_path>"
    exit 1
fi

# Read command-line arguments
num_requests=$1
file_path=$2

# Validate the number of requests
if ! [[ "$num_requests" =~ ^[0-9]+$ ]]; then
    echo "Error: Number of requests must be a positive integer."
    exit 1
fi

# Validate the file path
if [ ! -f "$file_path" ]; then
    echo "Error: File not found at path $file_path."
    exit 1
fi

# Record the start time
start_time=$(date +%s%3N)

# Loop to send the requests
for ((i=1; i<=num_requests; i++)); do
    curl -s -o ./output.csv -F "file=@$file_path" http://localhost:3333/upload &
done

# Wait for all background processes to complete
wait

# Record the end time
end_time=$(date +%s%3N)

# Calculate elapsed time in milliseconds
elapsed_time=$((end_time - start_time))

echo "All requests have been sent."
echo "Elapsed time: ${elapsed_time} ms"

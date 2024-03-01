#!/bin/bash

# Start time of the script
start_time=$(date +%s)

# Function to get elapsed time since start of script
get_elapsed_time() {
    current_time=$(date +%s)
    elapsed=$((current_time - start_time))
    echo "$elapsed seconds"
}

# Loop to execute commands every 15 seconds
while true; do
    # Get elapsed time
    elapsed_time=$(get_elapsed_time)

    # Output file names
    autoscaling="./output/autoscaling.log"
    resources="./output/resources.log"

    # Execute commands and append results to respective output files
    {
        echo "Elapsed Time: $elapsed_time"
        autoscaling_output=$(kubectl get horizontalpodautoscaler.autoscaling)
        echo "$autoscaling_output"
        echo "----------------------------------------"
    } >> "$autoscaling"

    {
        echo "Elapsed Time: $elapsed_time"
        resources_output=$(kubectl top pods | grep -E '^(NAME|subtask|task)')
        echo "$resources_output"
        echo "----------------------------------------"
    } >> "$resources"

    # Wait for 10 seconds
    sleep 10
done
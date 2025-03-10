#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Check if a process is running on port 8080
echo "Checking for a process on port 8080..."
PID=$(lsof -ti:8080 || true)

# If a process is found, terminate it
if [[ -n "$PID" ]]; then
  echo "A process was found on port 8080 with PID: $PID. Terminating it..."
  kill -9 "$PID"
else
  echo "No processes found on port 8080"
fi

# Build the application
echo "Building the application..."
go build -o server ./cmd

# Run the application in the background, redirecting output to a log file
echo "Starting the application..."
nohup ./server > output.log 2>&1 &
echo "Application started with PID: $!"

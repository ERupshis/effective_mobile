#!/bin/bash
# wait-for-it.sh

# Usage: ./wait-for-it.sh host:port [host:port] ... -- command args
# Example: ./wait-for-it.sh db:5432 -- my_command --arg1 --arg2

set -e

# Extract the last argument, which should be the command to run
cmd="${!#}"

# Remove the last argument from the argument list
args=("${@:1:$#-1}")

# Function to check if a host:port is available
wait_for() {
  local hostport="$1"
  local timeout=${2:-15} # Default timeout: 15 seconds
  local start_time=$(date +%s)
  local end_time=$((start_time + timeout))

  while true; do
    nc -z ${hostport} >/dev/null 2>&1
    result=$?

    if [ $result -eq 0 ]; then
      return 0
    fi

    current_time=$(date +%s)

    if [ $current_time -ge $end_time ]; then
      echo "Timeout! ${hostport} is not available after ${timeout} seconds."
      return 1
    fi

    sleep 1
  done
}

for arg in "${args[@]}"; do
  if [ "$arg" == "--" ]; then
    # When we encounter "--", it means the end of host:port arguments
    break
  fi

  wait_for "$arg"
done

# Remove the host:port arguments from the command line
args=("${args[@]#*--}")

# Execute the specified command
exec "${args[@]}"

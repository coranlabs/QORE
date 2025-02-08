#!/bin/bash

# Check which service to run based on the SERVICE_NAME environment variable
case "$SERVICE_NAME" in
  sbi)
    echo "Running SBI service..."
    exec /nwdaf/coranlabs_sbi
    ;;
  events)
    echo "Running Events service..."
    exec /nwdaf/coranlabs_events
    ;;
  analytics)
    echo "Running Analytics service..."
    exec /nwdaf/coranlabs_analytics
    ;;
  engine)
    echo "Running Engine service..."
    exec /nwdaf/coranlabs_engine
    ;;
  nbiml)
    echo "Running NBIML service..."
    exec /nwdaf/coranlabs_nbiml
    ;;
  *)
    echo "Error: Unknown or unset SERVICE_NAME environment variable."
    echo "Please set SERVICE_NAME to one of: sbi, events, analytics, engine, nbiml."
    exit 1
    ;;
esac

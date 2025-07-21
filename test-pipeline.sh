#!/bin/bash

echo "🚀 Starting pipeline test..."

# Start ingestor in background
echo "📊 Starting log ingestor..."
go run ./cmd/log-ingestor/main.go &
INGESTOR_PID=$!

# Wait for ingestor to start
sleep 2

# Run generator for 5 seconds
echo "📝 Starting log generator..."
timeout 5 go run ./cmd/log-generator/main.go

# Wait a bit for processing
sleep 2

# Check database
echo "🔍 Checking database..."
docker exec postgres psql -U loguser -d logdb -c 'SELECT COUNT(*) as total_logs FROM log_raw;'
docker exec postgres psql -U loguser -d logdb -c 'SELECT * FROM log_raw LIMIT 5;'

# Kill ingestor
kill $INGESTOR_PID

echo "✅ Test complete!"
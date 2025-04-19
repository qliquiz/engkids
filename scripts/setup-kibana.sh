#!/bin/bash

echo "Waiting for Kibana to start..."
while ! curl -s http://kibana:5601/api/status | grep -q '"status":{"overall":{"level":"available"'; do
  sleep 5
done

echo "Kibana is up. Setting up index patterns..."

curl -X POST "http://kibana:5601/api/saved_objects/index-pattern/engkids-logs" \
  -H "kbn-xsrf: true" \
  -H "Content-Type: application/json" \
  -d'{"attributes":{"title":"engkids-logs-*","timeFieldName":"@timestamp"}}'

echo "Index pattern created!"

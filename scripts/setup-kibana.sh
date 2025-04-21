#!/bin/bash

EXISTS=$(curl -s -o /dev/null -w "%{http_code}" http://kibana:5601/api/saved_objects/index-pattern/engkids-logs)

if [ "$EXISTS" -ne 200 ]; then
  echo "Creating index pattern..."
  curl -X POST "http://kibana:5601/api/saved_objects/index-pattern/engkids-logs" \
    -H 'kbn-xsrf: true' \
    -H 'Content-Type: application/json' \
    -d '{
      "attributes": {
        "title": "engkids-logs",
        "timeFieldName": "@timestamp"
      }
    }'
else
  echo "Index pattern already exists, skipping"
fi

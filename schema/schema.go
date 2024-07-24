package schema

import (
	"github.com/ONSdigital/dp-kafka/v3/avro"
)

var contentPublishedEvent = `{
  "type": "record",
  "name": "content-published",
  "fields": [
    {"name": "uri", "type": "string", "default": ""},
    {"name": "data_type", "type": "string", "default": ""},
    {"name": "collection_id", "type": "string", "default": ""},
    {"name": "job_id", "type": "string", "default": ""},
    {"name": "search_index", "type": "string", "default": ""},
    {"name": "trace_id", "type": "string", "default": ""}
  ]
}`

// ContentPublishedEvent is the Avro schema for contentPublishedEvent messages.
var ContentPublishedEvent = &avro.Schema{
	Definition: contentPublishedEvent,
}

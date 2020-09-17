package producer

import (
	"encoding/binary"
	"json-kafka/parser"
	"log"
	"time"

	"github.com/Shopify/sarama"
	"github.com/hamba/avro"
)

var (
	SASLUsername  = ""
	SASLPassword  = ""
	Topic         = ""
	MessageKey    = ""
	SchemaVersion = 1
	Brokers       []string
	Producer      sarama.AsyncProducer = nil
)

func NewProducer() sarama.AsyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond
	if SASLUsername != "" || SASLPassword != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = SASLUsername
		config.Net.SASL.Password = SASLPassword
	}
	if Topic == "" {
		log.Fatalln("Failed to start producer: topic must be non-null")
	}
	producer, err := sarama.NewAsyncProducer(Brokers, config)
	if err != nil {
		log.Fatalln("Failed to start producer:", err)
	}

	go func() {
		for err := range producer.Errors() {
			log.Println("Failed to add to Kafka:", err)
		}
	}()

	parser.KafkaChannel = make(chan []byte, 10)
	var schemaIdentifier uint32 = uint32(SchemaVersion)
	var schemaIdentifierBuffer []byte = make([]byte, 4)
	var confluentAvroHeader []byte
	binary.BigEndian.PutUint32(schemaIdentifierBuffer, schemaIdentifier)
	confluentAvroHeader = append([]byte{0}, schemaIdentifierBuffer...)
	go func() {
		for object := range parser.KafkaChannel {
			confluentMessage := append(confluentAvroHeader, object...)
			producer.Input() <- &sarama.ProducerMessage{
				Topic: Topic,
				Key:   sarama.StringEncoder(MessageKey),
				Value: sarama.ByteEncoder(confluentMessage),
			}
		}
	}()

	return producer
}

func NewAvroCodec() (avro.Schema, error) {
	schema := `{
		"type": "record",
		"name": "dns_record",
		"fields" : [ {
			"name" : "timestamp",
			"type" : "long"
		  }, {
			"name" : "ip_version",
			"type" : "int"
		  }, {
			"name" : "ip_src",
			"type" : "string"
		  }, {
			"name" : "ip_dst",
			"type" : "string"
		  }, {
			"name" : "dst_port",
			"type" : "int"
		  }, {
			"name" : "txid",
			"type" : "int"
		  }, {
			"name" : "rcode",
			"type" : "int"
		  }, {
			"name" : "qtype",
			"type" : "int"
		  }, {
			"name" : "qname",
			"type" : "string"
		  }, {
			"name" : "recursion_desired",
			"type" : "boolean"
		  }, {
			"name" : "response",
			"type" : "boolean"
		  }, {
			"name" : "answer",
			"type" : [ "null", "boolean" ],
			"default": null
		  }, {
			"name" : "authority",
			"type" : [ "null", "boolean" ],
			"default": null
		  }, {
			"name" : "additional",
			"type" : [ "null", "boolean" ],
			"default": null
		  }, {
			"name" : "rname",
			"type" : [ "null", "string" ],
			"default": null
		  }, {
			"name" : "rtype",
			"type" : [ "null", "int" ],
			"default": null
		  }, {
			"name" : "rdata",
			"type" : [ "null", "string" ],
			"default": null
		  }, {
			"name" : "ttl",
			"type" : [ "null", "long" ],
			"default": null
		  }, {
			"name" : "ecs_client",
			"type" : [ "null", "string" ],
			"default": null
		  }, {
			"name" : "ecs_source",
			"type" : [ "null", "string" ],
			"default": null
		  }, {
			"name" : "ecs_scope",
			"type" : [ "null", "string" ],
			"default": null
		  }, {
			"name" : "source",
			"type" : "string"
		  }, {
			"name" : "sensor",
			"type" : [ "null", "string" ],
			"default": null
		  } ]
		}`
	return avro.Parse(schema)
}

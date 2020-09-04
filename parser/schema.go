package parser

import (
	"github.com/hamba/avro"
	"github.com/hamba/avro/ocf"
)

var (
	AvroFileEncoder *ocf.Encoder
	KafkaChannel    chan []byte
	AvroSchema      avro.Schema
)

type DnsSchema struct {
	Timestamp          int64   `json:"timestamp" avro:"timestamp"`
	Sha256             string  `json:"sha256" avro:"-"`
	Udp                bool    `json:"udp" avro:"udp"`
	Ipv4               bool    `json:"ipv4" avro:"-"`
	SourceAddress      string  `json:"src_address" avro:"ip_src"`
	SourcePort         int     `json:"src_port" avro:"-"`
	DestinationAddress string  `json:"dst_address" avro:"ip_dst"`
	DestinationPort    int     `json:"dst_port" avro:"dst_port"`
	Id                 int     `json:"id" avro:"txid"`
	Rcode              int     `json:"rcode" avro:"rcode"`
	Truncated          bool    `json:"truncated" avro:"truncated"`
	Response           bool    `json:"response" avro:"response"`
	RecursionDesired   bool    `json:"recursion_desired" avro:"recursion_desired"`
	Answer             bool    `json:"answer" avro:"-"`
	Authority          bool    `json:"authority" avro:"-"`
	Additional         bool    `json:"additional" avro:"-"`
	Qname              string  `json:"qname" avro:"qname"`
	Qtype              int     `json:"qtype" avro:"qtype"`
	Ttl                *uint32 `json:"ttl" avro:"-"`
	Rname              *string `json:"rname" avro:"-"`
	Rtype              *uint16 `json:"rtype" avro:"-"`
	Rdata              *string `json:"rdata" avro:"-"`
	EcsClient          *string `json:"ecs_client" avro:"-"`
	EcsSource          *uint8  `json:"ecs_source" avro:"-"`
	EcsScope           *uint8  `json:"ecs_scope" avro:"-"`
	Source             string  `json:"source,omitempty" avro:"source"`
	Sensor             string  `json:"sensor,omitempty" avro:"-"`

	IPVersionAvro  int                    `avro:"ip_version" json:"-"`
	AnswerAvro     map[string]interface{} `avro:"answer" json:"-"`
	AuthorityAvro  map[string]interface{} `avro:"authority" json:"-"`
	AdditionalAvro map[string]interface{} `avro:"additional" json:"-"`
	TTLAvro        longUnion              `avro:"ttl,omitempty" json:"-"`
	RnameAvro      stringUnion            `avro:"rname,omitempty" json:"-"`
	RtypeAvro      intUnion               `avro:"rtype,omitempty" json:"-"`
	RdataAvro      stringUnion            `avro:"rdata,omitempty" json:"-"`
	EcsClientAvro  stringUnion            `avro:"ecs_client,omitempty" json:"-"`
	EcsSourceAvro  stringUnion            `avro:"ecs_source,omitempty" json:"-"`
	EcsScopeAvro   stringUnion            `avro:"ecs_scope,omitempty" json:"-"`
	SensorAvro     stringUnion            `avro:"sensor,omitempty" json:"-"`
}

type stringUnion struct {
	String string `avro:"string"`
}
type intUnion struct {
	Int int `avro:"int"`
}
type longUnion struct {
	Long int `avro:"long"`
}

func GetAvroSchema() string {
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
	return schema
}

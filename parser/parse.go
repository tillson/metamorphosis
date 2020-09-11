package parser

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"os"

	"github.com/hamba/avro"
	"gopkg.in/urfave/cli.v2"
)

func ParseFile(file string, c *cli.Context) {
	var s *bufio.Scanner
	if file == "-" {
		s = bufio.NewScanner(os.Stdin)
	} else {
		f, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		s = bufio.NewScanner(f)
	}

	if c.String("input") == "json" {
		for s.Scan() {
			var schema DnsSchema
			if err := json.Unmarshal(s.Bytes(), &schema); err != nil {
				log.Fatal(err)
			}
			schema.ToAvro(c)
		}
	} else if c.String("input") == "avro" {
		s.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			if atEOF && len(data) == 0 {
				return 0, nil, nil
			}
			if i := bytes.Index(data, []byte{0xC3, 0x01}); i >= 0 {
				return i + 1, data[0:i], nil
			}
			if atEOF {
				return len(data), data, nil
			}
			return
		})
		for s.Scan() {
			KafkaChannel <- s.Bytes()
		}

	}
	if s.Err() != nil {
		log.Fatal(s.Err())
	}

}

func (d DnsSchema) ToAvro(c *cli.Context) {
	d.AnswerAvro = map[string]interface{}{"boolean": d.Answer}
	d.AuthorityAvro = map[string]interface{}{"boolean": d.Authority}
	d.AdditionalAvro = map[string]interface{}{"boolean": d.Additional}

	if d.Rname != nil {
		d.RnameAvro = stringUnion{String: *d.Rname}
	}
	if d.Rtype != nil {
		d.RtypeAvro = intUnion{Int: int(*d.Rtype)}
	}
	if d.Rdata != nil {
		d.RdataAvro = stringUnion{String: *d.Rdata}
	}
	if d.Ttl != nil {
		d.TTLAvro = longUnion{Long: int(*d.Ttl)}
	}
	if d.EcsClient != nil {
		d.EcsClientAvro = stringUnion{String: *d.EcsClient}
	}
	if d.EcsSource != nil {
		d.EcsSourceAvro = stringUnion{String: string(*d.EcsSource)}
	}
	if d.EcsScope != nil {
		d.EcsScopeAvro = stringUnion{String: string(*d.EcsScope)}
	}
	if d.Source != "" {
		d.SensorAvro = stringUnion{String: d.Sensor}
	}

	if d.Ipv4 {
		d.IPVersionAvro = 4
	} else {
		d.IPVersionAvro = 6
	}
	if c.String("output") == "kafka" {
		data, err := avro.Marshal(AvroSchema, d)
		if err != nil {
			log.Fatal("Error while encoding avro %d", err)
		}
		KafkaChannel <- data
	} else {
		err := AvroFileEncoder.Encode(d)
		if err != nil {
			log.Fatal("Error while encoding avro %d", err)
		}
	}
}

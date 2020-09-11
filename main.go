package main

import (
	"json-kafka/parser"
	"json-kafka/producer"
	"log"
	"os"
	"time"

	"github.com/hamba/avro"
	"github.com/hamba/avro/ocf"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "metamorphosis"
	app.Description = "Convert JSON to avro and send it to Kafka"
	app.Usage = "cat dns.json | metamorphosis"
	app.Version = "1.0.0"
	app.Compiled = time.Now()

	app.Authors = []*cli.Author{
		{
			Name:  "Tillson Galloway",
			Email: "tillson@gatech.edu",
		},
	}

	app.Action = func(c *cli.Context) error {
		setupAvro(c)
		for _, f := range c.Args().Slice() {
			parser.ParseFile(f, c)
		}
		return nil
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "input",
			Usage: "Input type (json, avro)",
			Value: "json",
		},
		&cli.StringFlag{
			Name:  "output",
			Usage: "Output type (stdout (default), file, kafka)",
			Value: "stdout",
		},
		&cli.StringFlag{
			Name:  "config-file",
			Usage: "specify location of config file containing Kafka options (default config.yml)",
			Value: "config.yml",
		},
		&cli.StringFlag{
			Name:  "output-file",
			Usage: "specify location of output file",
		},
		// Kafka
		&cli.StringSliceFlag{
			Name:  "brokers",
			Usage: "list of kafka brokers",
		},
		&cli.StringFlag{
			Name:  "username",
			Usage: "kafka sasl username",
		},
		&cli.StringFlag{
			Name:  "password",
			Usage: "kafka sasl password",
		},
		&cli.StringFlag{
			Name:  "topic",
			Usage: "kafka topic",
		},
		&cli.StringFlag{
			Name:  "message-key",
			Usage: "kafka message key",
		},
		&cli.IntFlag{
			Name:  "schema-version",
			Usage: "kafka schema version",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func setupAvro(c *cli.Context) {
	codec := parser.GetAvroSchema()
	if c.String("output") == "kafka" {
		producer.Brokers = c.StringSlice("brokers")
		producer.SASLUsername = c.String("username")
		producer.SASLPassword = c.String("password")
		producer.Topic = c.String("topic")
		producer.MessageKey = c.String("message-key")
		producer.SchemaVersion = c.Int("schema-version")

		producer.Producer = producer.NewProducer()
		parser.AvroSchema = avro.MustParse(parser.GetAvroSchema())

	} else {
		var (
			encoder *ocf.Encoder
			err     error
		)
		if c.String("output") == "stdout" {
			encoder, err = ocf.NewEncoder(codec, os.Stdout)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			f, err := os.Create(c.String("output-file"))
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			encoder, err = ocf.NewEncoder(codec, f)
			if err != nil {
				log.Fatal("Failed to initialize avro encoder.")
			}
		}
		parser.AvroFileEncoder = encoder
	}
}

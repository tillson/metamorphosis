# Metamorphosis 
Transforms JSON to avro format and feeds to a kafka consumer. Options are specified by command-line arguments. 

## Usage
```
NAME:
   metamorphosis - json->avro->kafka transformer

USAGE:
   metamorphosis [global options] command [command options] [arguments...]

VERSION:
   1.0.0

DESCRIPTION:
   Convert JSON to avro and send it to Kafka

AUTHOR:
   Tillson Galloway <tillson@gatech.edu>

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --output value          Output type (stdout (default), file, kafka) (default: "stdout")
   --config-file value     specify location of config file containing Kafka options (default config.yml) (default: "config.yml")
   --output-file value     specify location of output file
   --brokers value         list of kafka brokers
   --username value        kafka sasl username
   --password value        kafka sasl password
   --topic value           kafka topic
   --message-key value     kafka message key
   --schema-version value  kafka schema version (default: 0)
   --help, -h              show help (default: false)
   --version, -v           print the version (default: false)
```

### Examples
`cat data.json | metamorphosis - `  
`metamorphosis --output file --output-file out.avro data.json`  
`metamorphosis --output kafka --brokers localhost:9092 --topic active-data data.avro data.json`
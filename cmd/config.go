package cmd

type Configs struct {
	GeoClientURL                  string
	KafkaHost                     string
	ConsumerGroup                 string
	KafkaOrdersCreateTopic        string
	KafkaOrdersStatusChangedTopic string
}

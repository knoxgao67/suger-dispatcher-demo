package kafka

// Init call after config.Init.
func Init() {
	initConsumer()
	initProducer()
}

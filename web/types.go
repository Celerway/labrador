package web

type mqttMessage struct {
	Topic   string
	Payload []byte
}

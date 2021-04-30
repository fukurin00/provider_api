package main

import (
	"encoding/json"
	"log"
	"sync"

	synerex "github.com/fukurin00/provider_api"
	sxmqtt "github.com/synerex/proto_mqtt"
	"google.golang.org/protobuf/proto"
)

func main() {
	wg := sync.WaitGroup{}
	channels := []uint32{synerex.MQTT_GATEWAY_SVC, synerex.JSON_DATA_SVC}
	names := []string{"mqtt_sample", "json_sample"}
	wg.Add(1)
	s, err := synerex.NewSynerexConfig("sample", channels, names)
	if err != nil {
		log.Print("failure on Starting Synerex Provider ..", err)
	} else {
		topic := "test/sample"
		msg := `{"sample": "test"}`
		jmsg, err := json.Marshal(msg)
		if err != nil {
			log.Print(err)
		}
		mqttRec := sxmqtt.MQTTRecord{
			Topic:  topic,
			Record: jmsg,
		}

		out, err := proto.Marshal(&mqttRec)
		if err != nil {
			log.Print(err)
		}
		id, err := s.NotifySupply(out, synerex.MQTT_GATEWAY_SVC, "testMQTTMessage")
		if err != nil {
			log.Print(err)
		}
		log.Printf("send message to id %d", id)
	}
	wg.Wait()
}

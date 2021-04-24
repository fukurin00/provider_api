package main

import (
	"fmt"
	"log"
	"sync"

	synerex "github.com/fukurin00/provider_api/src/api"
	sxmqtt "github.com/synerex/proto_mqtt"
	api "github.com/synerex/synerex_api"
	sxutil "github.com/synerex/synerex_sxutil"
	"google.golang.org/protobuf/proto"
)

func mqttCallback(clt *sxutil.SXServiceClient, sp *api.Supply) {
	record := sxmqtt.MQTTRecord{}
	err := proto.Unmarshal(sp.Cdata.Entity, &record)
	if err != nil {
		log.Print(err)
	}
	log.Printf("Receive MQTT Topic:%s", record.Topic)
	fmt.Println(record.Record)
}

func main() {
	wg := sync.WaitGroup{}
	channels := []uint32{synerex.MQTT_GATEWAY_SVC}
	names := []string{"mqtt_sample"}
	wg.Add(1)
	s, err := synerex.NewSynerexConfig("sample", channels, names)
	if err != nil {
		log.Print("Failure on Starting Synerex Provider ..", err)
	} else {
		s.SubscribeSupply(channels[0], mqttCallback)
	}
	wg.Wait()
}

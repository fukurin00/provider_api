package main

import (
	"fmt"
	"log"

	synerex "bitbucket.org/uclabnu/synerex_provider_api/api"
	sxmqtt "github.com/synerex/proto_mqtt"
	api "github.com/synerex/synerex_api"
	sxutil "github.com/synerex/synerex_sxutil"
	"google.golang.org/protobuf/proto"
)

func mqttCallback(clt *sxutil.SXServiceClient, sp *api.Supply) {
	record := sxmqtt.MQTTRecord{}
	err := proto.Unmarshal(sp.Cdata.Entity, &record)
	if err != nil {
		log.Printf("Receive MQTT Topic:%s")
		fmt.Print(record.Record)
	}
}

func main() {
	channels := [1]uint32{synerex.MQTT_GATEWAY_SVC}
	names := [1]string{"mqtt_sample"}
	s, err := synerex.NewSynerexConfig("sample", channels, names)
	if err != nil {
		log.Print("Failure on Starting Synerex Provider ..", err)
	} else {
		s.SubscribeSupply(channels[0], mqttCallback)
	}
}

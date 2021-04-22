package api

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	proto "github.com/synerex/synerex_proto"
	sxutil "github.com/synerex/synerex_sxutil"
)

const (
	RIDE_SHARE         = proto.RIDE_SHARE
	AD_SERVICE         = proto.AD_SERVICE
	LIB_SERVICE        = proto.LIB_SERVICE
	PT_SERVICE         = proto.PT_SERVICE
	ROUTING_SERVICE    = proto.ROUTING_SERVICE
	MARKETING_SERVICE  = proto.MARKETING_SERVICE
	FLUENTD_SERVICE    = proto.FLUENTD_SERVICE
	MEETING_SERVICE    = proto.MEETING_SERVICE
	STORAGE_SERVICE    = proto.STORAGE_SERVICE
	RETRIEVAL_SERVICE  = proto.RETRIEVAL_SERVICE
	PEOPLE_COUNTER_SVC = proto.PEOPLE_COUNTER_SVC
	AREA_COUNTER_SVC   = proto.AREA_COUNTER_SVC
	PEOPLE_AGENT_SVC   = proto.PEOPLE_AGENT_SVC
	GEOGRAPHIC_SVC     = proto.GEOGRAPHIC_SVC
	JSON_DATA_SVC      = proto.JSON_DATA_SVC
	MQTT_GATEWAY_SVC   = proto.MQTT_GATEWAY_SVC
	WAREHOUSE_SVC      = proto.WAREHOUSE_SVC
	PEOPLE_FLOW_SVC    = proto.PEOPLE_FLOW_SVC
	GRIDEYE_SVC        = proto.GRIDEYE_SVC
	LATENT_DMD_SVC     = proto.LATENT_DMD_SVC
	LATENT_DMD_DSP_SVC = proto.LATENT_DMD_DSP_SVC
	ALT_PT_SVC         = proto.ALT_PT_SVC
)

var (
	Nodesrv = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	Mu      sync.Mutex
)

type SynerexConfig struct {
	Nodesrv         string
	SxServerAddress string
	ChannelList     map[uint32]*sxutil.SXServiceClient
}

// Constructor for synerex config
func NewSynerexConfig(nodeName string, channelTypes []uint32, channelNames []string) (s *SynerexConfig, oerr error) {
	s = new(SynerexConfig)
	s.ChannelList = make(map[uint32]*sxutil.SXServiceClient)
	flag.Parse() // load command line arguments

	//kill all process by Crtl+C
	go sxutil.HandleSigInt()
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)

	srv, err := sxutil.RegisterNode(*Nodesrv, nodeName, channelTypes, nil)
	s.SxServerAddress = srv
	if err != nil {
		oerr = errors.New("Cannot register node...")
		return s, oerr
	}

	// connect each channel
	for i := 0; i < len(channelTypes); i++ {
		s.startSingleChannel(channelTypes[i], channelNames[i])
	}
	return s, nil
}

func (s *SynerexConfig) startSingleChannel(channelType uint32, clientName string) {
	client := sxutil.GrpcConnectServer(s.SxServerAddress)
	argJSON := fmt.Sprintf("{Client:%s}", clientName)
	sxClient := sxutil.NewSXServiceClient(client, channelType, argJSON)

	s.ChannelList[channelType] = sxClient
}

func (s SynerexConfig) ReconnectClient(client *sxutil.SXServiceClient) {
	Mu.Lock()
	if client.SXClient != nil {
		client.SXClient = nil
		log.Printf("Client reset \n")
	}
	Mu.Unlock()
	time.Sleep(5 * time.Second) // wait 5 seconds to reconnect
	Mu.Lock()
	if client.SXClient == nil {
		newClt := sxutil.GrpcConnectServer(s.SxServerAddress)
		if newClt != nil {
			// log.Printf("Reconnect server [%s]\n", s.SxServerAddress)
			client.SXClient = newClt
		}
	}
	Mu.Unlock()
}

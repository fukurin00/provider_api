package api

import (
	"context"

	api "github.com/synerex/synerex_api"
	sxutil "github.com/synerex/synerex_sxutil"
)

// Subscribe all supply in specified channel
func (s SynerexConfig) SubscribeSupply(channelType uint32, callback func(clt *sxutil.SXServiceClient, sp *api.Supply)) {
	go s.callSubscribeSupply(channelType, callback)
}

func (s SynerexConfig) callSubscribeSupply(channelType uint32, callback func(clt *sxutil.SXServiceClient, sp *api.Supply)) {
	ctx := context.Background()
	for { // make it continuously working..
		client := s.ChannelList[channelType]
		err := client.SubscribeSupply(ctx, callback)
		if err != nil {
			s.ReconnectClient(client)
		}
	}
}


// Send supply to all providers
func (s SynerexConfig) NotifySupply(protocolBuffer []byte, channelType uint32, supplyName string) (uint64, error) {
	client := s.ChannelList[channelType]
	cData := api.Content{Entity: protocolBuffer}
	supplyOpt := sxutil.SupplyOpts{
		Name:  supplyName,
		JSON:  client.ArgJson,
		Cdata: &cData,
	}
	id, err := client.NotifySupply(&supplyOpt)
	return id, err
}

// send suuply to target provider
func (s SynerexConfig) ProposeSupply(protocolBuffer []byte, channelType uint32, target uint64, supplyName string) uint64 {
	client := s.ChannelList[channelType]
	cData := api.Content{Entity: protocolBuffer}
	supplyOpt := sxutil.SupplyOpts{
		Name:   supplyName,
		Target: target,
		JSON:   client.ArgJson,
		Cdata:  &cData,
	}
	id := client.ProposeSupply(&supplyOpt)
	return id
}

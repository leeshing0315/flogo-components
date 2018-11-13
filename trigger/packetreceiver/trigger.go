package packetreceiver

import (
	"fmt"
	"time"

	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// MyTriggerFactory My Trigger factory
type MyTriggerFactory struct {
	metadata *trigger.Metadata
}

// NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &MyTriggerFactory{metadata: md}
}

// New Creates a new trigger instance for a given id
func (t *MyTriggerFactory) New(config *trigger.Config) trigger.Trigger {
	return &MyTrigger{metadata: t.metadata, config: config}
}

// MyTrigger is a stub for your Trigger implementation
type MyTrigger struct {
	metadata *trigger.Metadata
	config   *trigger.Config
}

// Initialize implements trigger.Init.Initialize
func (t *MyTrigger) Initialize(ctx trigger.InitContext) error {
	return nil
}

// Metadata implements trigger.Trigger.Metadata
func (t *MyTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Start implements trigger.Trigger.Start
func (t *MyTrigger) Start() error {
	version := pcap.Version()
	fmt.Println(version)
	devices, _ := pcap.FindAllDevs()
	fmt.Println(devices)
	handle, _ := pcap.OpenLive("\\Device\\NPF_{CFD7F7F5-C1F9-428B-8CAE-CB86996ED02E}", int32(65535), false, -1*time.Second)
	// handle.SetBPFFilter("tcp and port 8033")
	defer handle.Close()
	packetSource := gopacket.NewPacketSource(
		handle,
		handle.LinkType(),
	)
	for packet := range packetSource.Packets() {
		// Let's see if the packet is IP (even though the ether type told us)
		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		if ipLayer != nil {
			fmt.Println("IPv4 layer detected.")
			ip, _ := ipLayer.(*layers.IPv4)
			fmt.Printf("From %s to %s\n", ip.SrcIP, ip.DstIP)
			fmt.Println("Protocol: ", ip.Protocol)
			fmt.Println("Payload: ", ip.Payload)
			fmt.Println()
		}

		// Iterate over all layers, printing out each layer type
		for _, layer := range packet.Layers() {
			fmt.Println("PACKET LAYER:", layer.LayerType())
		}
	}
	return nil
}

// Stop implements trigger.Trigger.Start
func (t *MyTrigger) Stop() error {
	// stop the trigger
	return nil
}

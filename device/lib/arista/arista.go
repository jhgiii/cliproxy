package arista

import (
	"encoding/json"
	"fmt"
	"net/netip"
)

type Arista struct {
	Hostname string
	IP       netip.Addr
}

func (a *Arista) DiscoverIPAddresses() {
	// TODO
}

func parseIPFromCli(b []byte) ([]netip.Addr, error) {
	var address []netip.Addr
	var Ints ResponseShowIntBrief
	err := json.Unmarshal(b, &Ints)
	if err != nil {
		fmt.Println(err)
	}

	for _, i := range Ints.Interfaces {

		a, err := netip.ParseAddr(i.InterfaceAddress.IPAddr.Address)
		if err != nil {
			continue
		}
		address = append(address, a)
	}
	return address, nil
}

type ResponseShowIntBrief struct {
	Interfaces map[string]Interfaces `json:"interfaces"`
}
type IPAddr struct {
	Address string `json:"address"`
	MaskLen int    `json:"maskLen"`
}
type InterfaceAddress struct {
	IPAddr IPAddr `json:"ipAddr"`
}
type Interfaces struct {
	Name                  string           `json:"name"`
	LineProtocolStatus    string           `json:"lineProtocolStatus"`
	InterfaceStatus       string           `json:"interfaceStatus"`
	Mtu                   int              `json:"mtu"`
	Ipv4Routable240       bool             `json:"ipv4Routable240"`
	Ipv4Routable0         bool             `json:"ipv4Routable0"`
	InterfaceAddress      InterfaceAddress `json:"interfaceAddress"`
	NonRoutableClassEIntf bool             `json:"nonRoutableClassEIntf"`
}

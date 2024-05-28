package arista

import (
	"encoding/json"
	"io"
	"net/netip"
)

var validConnectMethods = [...]string{"ssh", "rest", "cvp"}

type Arista struct {
	Hostname      string
	IP            netip.Addr
	ConnectMethod string
}

func NewArista(hostname, ip, connectmethod string) (Arista, error) {
	var a Arista
	a.Hostname = hostname
	var err error
	a.IP, err = netip.ParseAddr(ip)
	if err != nil {
		return a, err
	}
	a.ConnectMethod = connectmethod
	return a, nil
}
func (a *Arista) DiscoverIPAddresses() error {
	// TODO

	return nil
}

func parseIPFromCli(r io.Reader) ([]netip.Addr, error) {
	var address []netip.Addr
	var Ints ResponseShowIntBrief
	rd, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rd, &Ints)
	if err != nil {
		return nil, err
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

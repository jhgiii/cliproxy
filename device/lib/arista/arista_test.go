package arista

import (
	"bytes"
	"net/netip"
	"os"
	"reflect"
	"slices"
	"testing"
)

func Test_parseIPFromCli(t *testing.T) {
	test_input, err := os.ReadFile("testinput/show_ip_int_brief")
	if err != nil {
		t.Errorf("Unable to open file for testing: %v", err)
	}
	have, _ := parseIPFromCli(bytes.NewReader(test_input))
	want := make([]netip.Addr, 3)
	want[0], _ = netip.ParseAddr("0.0.0.0")
	want[1], _ = netip.ParseAddr("172.20.20.2")
	want[2], _ = netip.ParseAddr("1.1.1.1")
	if len(have) != len(want) {
		t.Errorf("Mismatch in expected length in response.  Have: %v.  Want: %v", len(have), len(want))
	}

	for _, j := range have {
		if !slices.Contains(want, j) {
			t.Errorf("Failed. %v is not in our \"Want\"", j)
		}
	}
}

func TestNewArista(t *testing.T) {
	type args struct {
		hostname      string
		ip            string
		connectmethod string
	}
	tests := []struct {
		name    string
		args    args
		want    Arista
		wantErr bool
	}{
		{name: "test1",
			args: args{hostname: "test1", ip: "1.2.2.2", connectmethod: "ssh"},
			want: Arista{Hostname: "test1", IP: netip.MustParseAddr("1.2.2.2"), ConnectMethod: "ssh"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewArista(tt.args.hostname, tt.args.ip, tt.args.connectmethod)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewArista() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewArista() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArista_DiscoverIPAddresses(t *testing.T) {
	type fields struct {
		Hostname      string
		IP            netip.Addr
		ConnectMethod string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Arista{
				Hostname:      tt.fields.Hostname,
				IP:            tt.fields.IP,
				ConnectMethod: tt.fields.ConnectMethod,
			}
			if err := a.DiscoverIPAddresses(); (err != nil) != tt.wantErr {
				t.Errorf("Arista.DiscoverIPAddresses() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

package arista

import (
	"net/netip"
	"os"
	"slices"
	"testing"
)

func Test_parseIPFromCli(t *testing.T) {
	test_input, err := os.ReadFile("testinput/show_ip_int_brief")
	if err != nil {
		t.Errorf("Unable to open file for testing: %v", err)
	}
	have, _ := parseIPFromCli(test_input)
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

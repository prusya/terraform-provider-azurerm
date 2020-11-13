package tests

import (
	"testing"

	"github.com/prusya/terraform-provider-azurerm/azurerm/internal/services/network"
)

func TestParsePointToSiteVPNGateway(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *network.PointToSiteVPNGatewayResourceID
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "No Point-to-Site VPN Gateways Segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo",
			Expected: nil,
		},
		{
			Name:     "No Point-to-Site VPN Gateways Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/p2sVpnGateways/",
			Expected: nil,
		},
		{
			Name:  "Completed",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/p2sVpnGateways/example",
			Expected: &network.PointToSiteVPNGatewayResourceID{
				Name:          "example",
				ResourceGroup: "foo",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := network.ParsePointToSiteVPNGatewayID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}

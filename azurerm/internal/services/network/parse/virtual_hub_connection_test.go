package parse

import (
	"testing"

	"github.com/prusya/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = VirtualHubConnectionId{}

func TestVirtualHubConnectionIDFormatter(t *testing.T) {
	subscriptionId := "12345678-1234-5678-1234-123456789012"
	vhubid := NewVirtualHubID("group1", "vhub1")
	actual := NewVirtualHubConnectionID(vhubid, "conn1").ID(subscriptionId)
	expected := "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Network/virtualHubs/vhub1/hubVirtualNetworkConnections/conn1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestVirtualHubConnectionID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *VirtualHubConnectionId
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "No Virtual Hubs Segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo",
			Expected: nil,
		},
		{
			Name:     "No Virtual Hubs Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/virtualHubs/",
			Expected: nil,
		},
		{
			Name:     "No Hub Network Connections Segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/virtualHubs/example",
			Expected: nil,
		},
		{
			Name:     "No Virtual Hubs Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/virtualHubs/hubVirtualNetworkConnections/",
			Expected: nil,
		},
		{
			Name:  "Completed",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/virtualHubs/example/hubVirtualNetworkConnections/connection1",
			Expected: &VirtualHubConnectionId{
				Name:           "connection1",
				VirtualHubName: "example",
				ResourceGroup:  "foo",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := VirtualHubConnectionID(v.Input)
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

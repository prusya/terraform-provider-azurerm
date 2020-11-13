package parse

import (
	"testing"

	"github.com/prusya/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = RoutingRuleId{}

func TestRoutingRuleIDFormatter(t *testing.T) {
	subscriptionId := "12345678-1234-5678-1234-123456789012"
	frontDoorId := NewFrontDoorID("group1", "frontdoor1")
	actual := NewRoutingRuleID(frontDoorId, "rule1").ID(subscriptionId)
	expected := "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Network/frontDoors/frontdoor1/routingRules/rule1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestRoutingRuleIDParser(t *testing.T) {
	testData := []struct {
		input    string
		expected *RoutingRuleId
	}{
		{
			// lower case
			input:    "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Network/frontdoors/frontDoor1/routingrules/rule1",
			expected: nil,
		},
		{
			// camel case
			input: "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Network/frontDoors/frontDoor1/routingRules/rule1",
			expected: &RoutingRuleId{
				ResourceGroup: "group1",
				FrontDoorName: "frontDoor1",
				Name:          "rule1",
			},
		},
		{
			// title case
			input: "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Network/Frontdoors/frontDoor1/RoutingRules/rule1",
			expected: &RoutingRuleId{
				ResourceGroup: "group1",
				FrontDoorName: "frontDoor1",
				Name:          "rule1",
			},
		},
		{
			// pascal case
			input:    "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Network/FrontDoors/frontDoor1/Routingrules/rule1",
			expected: nil,
		},
	}
	for _, test := range testData {
		t.Logf("Testing %q..", test.input)
		actual, err := RoutingRuleID(test.input)
		if err != nil && test.expected == nil {
			continue
		} else {
			if err == nil && test.expected == nil {
				t.Fatalf("Expected an error but didn't get one")
			} else if err != nil && test.expected != nil {
				t.Fatalf("Expected no error but got: %+v", err)
			}
		}

		if actual.ResourceGroup != test.expected.ResourceGroup {
			t.Fatalf("Expected ResourceGroup to be %q but was %q", test.expected.ResourceGroup, actual.ResourceGroup)
		}

		if actual.FrontDoorName != test.expected.FrontDoorName {
			t.Fatalf("Expected FrontDoorName to be %q but was %q", test.expected.FrontDoorName, actual.FrontDoorName)
		}

		if actual.Name != test.expected.Name {
			t.Fatalf("Expected name to be %q but was %q", test.expected.Name, actual.Name)
		}
	}
}

package parse

import (
	"fmt"

	"github.com/prusya/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DnsARecordId struct {
	ResourceGroup string
	ZoneName      string
	Name          string
}

func DnsARecordID(input string) (*DnsARecordId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse DNS A Record ID %q: %+v", input, err)
	}

	record := DnsARecordId{
		ResourceGroup: id.ResourceGroup,
	}

	if record.ZoneName, err = id.PopSegment("dnszones"); err != nil {
		return nil, err
	}

	if record.Name, err = id.PopSegment("A"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &record, nil
}

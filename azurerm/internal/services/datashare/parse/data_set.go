package parse

import (
	"fmt"

	"github.com/prusya/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DataShareDataSetId struct {
	ResourceGroup string
	AccountName   string
	ShareName     string
	Name          string
}

func DataShareDataSetID(input string) (*DataShareDataSetId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse DataShareDataSet ID %q: %+v", input, err)
	}

	dataShareDataSet := DataShareDataSetId{
		ResourceGroup: id.ResourceGroup,
	}
	if dataShareDataSet.AccountName, err = id.PopSegment("accounts"); err != nil {
		return nil, err
	}
	if dataShareDataSet.ShareName, err = id.PopSegment("shares"); err != nil {
		return nil, err
	}
	if dataShareDataSet.Name, err = id.PopSegment("dataSets"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &dataShareDataSet, nil
}

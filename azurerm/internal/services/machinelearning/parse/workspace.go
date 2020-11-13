package parse

import (
	"github.com/prusya/terraform-provider-azurerm/azurerm/helpers/azure"
	accountParser "github.com/prusya/terraform-provider-azurerm/azurerm/internal/services/storage/parsers"
)

type WorkspaceId struct {
	Name          string
	ResourceGroup string
}

func WorkspaceID(input string) (*WorkspaceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	workspace := WorkspaceId{
		ResourceGroup: id.ResourceGroup,
	}

	if workspace.Name, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &workspace, nil
}

// TODO -- use parse function "github.com/prusya/terraform-provider-azurerm/azurerm/internal/services/storage/parsers".ParseAccountID
// when issue https://github.com/Azure/azure-rest-api-specs/issues/8323 is addressed
func AccountIDCaseDiffSuppress(input string) (*accountParser.AccountID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	account := accountParser.AccountID{
		ResourceGroup: id.ResourceGroup,
	}

	if account.Name, err = id.PopSegment("storageAccounts"); err != nil {
		if account.Name, err = id.PopSegment("storageaccounts"); err != nil {
			return nil, err
		}
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &account, nil
}

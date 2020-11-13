package azurerm

import (
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/prusya/terraform-provider-azurerm/azurerm/internal/provider"
)

func Provider() terraform.ResourceProvider {
	return provider.AzureProvider()
}

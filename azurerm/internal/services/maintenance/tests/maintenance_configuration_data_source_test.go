package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/prusya/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMMaintenanceConfiguration_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_maintenance_configuration", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMaintenanceConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMaintenanceConfiguration_complete(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "scope", "Host"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.env", "TesT"),
				),
			},
		},
	})
}

func testAccDataSourceMaintenanceConfiguration_complete(data acceptance.TestData) string {
	template := testAccAzureRMMaintenanceConfiguration_complete(data)
	return fmt.Sprintf(`
%s

data "azurerm_maintenance_configuration" "test" {
  name                = azurerm_maintenance_configuration.test.name
  resource_group_name = azurerm_maintenance_configuration.test.resource_group_name
}
`, template)
}

package datashare

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/prusya/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/prusya/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/prusya/terraform-provider-azurerm/azurerm/internal/services/datashare/validate"
	"github.com/prusya/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/prusya/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/prusya/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceDataShareAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmDataShareAccountRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.DataShareAccountName(),
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"identity": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmDataShareAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.AccountClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("DataShare Account %q does not exist in Resource Group %q", name, resourceGroup)
		}
		return fmt.Errorf("retrieving DataShare Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("reading DataShare Account %q (Resource Group %q): ID is empty or nil", name, resourceGroup)
	}

	d.SetId(*resp.ID)
	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if err := d.Set("identity", flattenAzureRmDataShareAccountIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}
	return tags.FlattenAndSet(d, resp.Tags)
}

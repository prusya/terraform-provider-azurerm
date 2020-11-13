package eventhub

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/eventhub/mgmt/2018-01-01-preview/eventhub"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/prusya/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/prusya/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/prusya/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/prusya/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/prusya/terraform-provider-azurerm/azurerm/internal/services/eventhub/migration"
	"github.com/prusya/terraform-provider-azurerm/azurerm/internal/services/eventhub/parse"
	"github.com/prusya/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/prusya/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmEventHubNamespaceAuthorizationRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmEventHubNamespaceAuthorizationRuleCreateUpdate,
		Read:   resourceArmEventHubNamespaceAuthorizationRuleRead,
		Update: resourceArmEventHubNamespaceAuthorizationRuleCreateUpdate,
		Delete: resourceArmEventHubNamespaceAuthorizationRuleDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 2,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    migration.EventHubNamespaceAuthorizationRuleUpgradeV0Schema().CoreConfigSchema().ImpliedType(),
				Upgrade: migration.EventHubNamespaceAuthorizationRuleUpgradeV0ToV1,
				Version: 0,
			},
			{
				Type:    migration.EventHubNamespaceAuthorizationRuleUpgradeV1Schema().CoreConfigSchema().ImpliedType(),
				Upgrade: migration.EventHubNamespaceAuthorizationRuleUpgradeV1ToV2,
				Version: 1,
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: azure.EventHubAuthorizationRuleSchemaFrom(map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateEventHubAuthorizationRuleName(),
			},

			"namespace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateEventHubNamespaceName(),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),
		}),

		CustomizeDiff: azure.EventHubAuthorizationRuleCustomizeDiff,
	}
}

func resourceArmEventHubNamespaceAuthorizationRuleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.NamespacesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM EventHub Namespace Authorization Rule creation.")

	name := d.Get("name").(string)
	namespaceName := d.Get("namespace_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.GetAuthorizationRule(ctx, resourceGroup, namespaceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing EventHub Namespace Authorization Rule %q (Namespace %q / Resource Group %q): %s", name, namespaceName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_eventhub_namespace_authorization_rule", *existing.ID)
		}
	}

	locks.ByName(namespaceName, eventHubNamespaceResourceName)
	defer locks.UnlockByName(namespaceName, eventHubNamespaceResourceName)

	parameters := eventhub.AuthorizationRule{
		Name: &name,
		AuthorizationRuleProperties: &eventhub.AuthorizationRuleProperties{
			Rights: azure.ExpandEventHubAuthorizationRuleRights(d),
		},
	}

	if _, err := client.CreateOrUpdateAuthorizationRule(ctx, resourceGroup, namespaceName, name, parameters); err != nil {
		return fmt.Errorf("Error creating/updating EventHub Namespace Authorization Rule %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.GetAuthorizationRule(ctx, resourceGroup, namespaceName, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read EventHub Namespace Authorization Rule %s (resource group %s) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmEventHubNamespaceAuthorizationRuleRead(d, meta)
}

func resourceArmEventHubNamespaceAuthorizationRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.NamespacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NamespaceAuthorizationRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetAuthorizationRule(ctx, id.ResourceGroup, id.NamespaceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Authorization Rule %q (EventHub Namespace %q / Resource Group %q) : %+v", id.Name, id.NamespaceName, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("namespace_name", id.NamespaceName)
	d.Set("resource_group_name", id.ResourceGroup)

	if properties := resp.AuthorizationRuleProperties; properties != nil {
		listen, send, manage := azure.FlattenEventHubAuthorizationRuleRights(properties.Rights)
		d.Set("manage", manage)
		d.Set("listen", listen)
		d.Set("send", send)
	}

	keysResp, err := client.ListKeys(ctx, id.ResourceGroup, id.NamespaceName, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving Keys for Authorization Rule %q (EventHub Namespace %q / Resource Group %q): %+v", id.Name, id.NamespaceName, id.ResourceGroup, err)
	}

	d.Set("primary_key", keysResp.PrimaryKey)
	d.Set("secondary_key", keysResp.SecondaryKey)
	d.Set("primary_connection_string", keysResp.PrimaryConnectionString)
	d.Set("secondary_connection_string", keysResp.SecondaryConnectionString)
	d.Set("primary_connection_string_alias", keysResp.AliasPrimaryConnectionString)
	d.Set("secondary_connection_string_alias", keysResp.AliasSecondaryConnectionString)

	return nil
}

func resourceArmEventHubNamespaceAuthorizationRuleDelete(d *schema.ResourceData, meta interface{}) error {
	eventhubClient := meta.(*clients.Client).Eventhub.NamespacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NamespaceAuthorizationRuleID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.NamespaceName, eventHubNamespaceResourceName)
	defer locks.UnlockByName(id.NamespaceName, eventHubNamespaceResourceName)

	if _, err := eventhubClient.DeleteAuthorizationRule(ctx, id.ResourceGroup, id.NamespaceName, id.Name); err != nil {
		return fmt.Errorf("deleting Authorization Rule %q (EventHub Namespace %q / Resource Group %q): %+v", id.Name, id.NamespaceName, id.ResourceGroup, err)
	}

	return nil
}

package cosmos

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/cosmos-db/mgmt/2020-04-01-preview/documentdb"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/prusya/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/prusya/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/prusya/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/prusya/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/prusya/terraform-provider-azurerm/azurerm/internal/services/cosmos/common"
	"github.com/prusya/terraform-provider-azurerm/azurerm/internal/services/cosmos/migration"
	"github.com/prusya/terraform-provider-azurerm/azurerm/internal/services/cosmos/parse"
	"github.com/prusya/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/prusya/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmCosmosDbCassandraKeyspace() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmCosmosDbCassandraKeyspaceCreate,
		Read:   resourceArmCosmosDbCassandraKeyspaceRead,
		Update: resourceArmCosmosDbCassandraKeyspaceUpdate,
		Delete: resourceArmCosmosDbCassandraKeyspaceDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    migration.ResourceCassandraKeyspaceUpgradeV0Schema().CoreConfigSchema().ImpliedType(),
				Upgrade: migration.ResourceCassandraKeyspaceStateUpgradeV0ToV1,
				Version: 0,
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosEntityName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosAccountName,
			},

			"throughput": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.CosmosThroughput,
			},

			"autoscale_settings": common.DatabaseAutoscaleSettingsSchema(),
		},
	}
}

func resourceArmCosmosDbCassandraKeyspaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CassandraClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	account := d.Get("account_name").(string)

	existing, err := client.GetCassandraKeyspace(ctx, resourceGroup, account, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for presence of creating Cosmos Cassandra Keyspace %q (Account: %q): %+v", name, account, err)
		}
	} else {
		if existing.ID == nil && *existing.ID == "" {
			return fmt.Errorf("Error generating import ID for Cosmos Cassandra Keyspace %q (Account: %q)", name, account)
		}

		return tf.ImportAsExistsError("azurerm_cosmosdb_cassandra_keyspace", *existing.ID)
	}

	db := documentdb.CassandraKeyspaceCreateUpdateParameters{
		CassandraKeyspaceCreateUpdateProperties: &documentdb.CassandraKeyspaceCreateUpdateProperties{
			Resource: &documentdb.CassandraKeyspaceResource{
				ID: &name,
			},
			Options: &documentdb.CreateUpdateOptions{},
		},
	}

	if throughput, hasThroughput := d.GetOk("throughput"); hasThroughput {
		if throughput != 0 {
			db.CassandraKeyspaceCreateUpdateProperties.Options.Throughput = common.ConvertThroughputFromResourceData(throughput)
		}
	}

	if _, hasAutoscaleSettings := d.GetOk("autoscale_settings"); hasAutoscaleSettings {
		db.CassandraKeyspaceCreateUpdateProperties.Options.AutoscaleSettings = common.ExpandCosmosDbAutoscaleSettings(d)
	}

	future, err := client.CreateUpdateCassandraKeyspace(ctx, resourceGroup, account, name, db)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for Cosmos Cassandra Keyspace %q (Account: %q): %+v", name, account, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting on create/update future for Cosmos Cassandra Keyspace %q (Account: %q): %+v", name, account, err)
	}

	resp, err := client.GetCassandraKeyspace(ctx, resourceGroup, account, name)
	if err != nil {
		return fmt.Errorf("Error making get request for Cosmos Cassandra Keyspace %q (Account: %q): %+v", name, account, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Error getting ID from Cosmos Cassandra Keyspace %q (Account: %q)", name, account)
	}

	d.SetId(*resp.ID)

	return resourceArmCosmosDbCassandraKeyspaceRead(d, meta)
}

func resourceArmCosmosDbCassandraKeyspaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CassandraClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CassandraKeyspaceID(d.Id())
	if err != nil {
		return err
	}

	err = common.CheckForChangeFromAutoscaleAndManualThroughput(d)
	if err != nil {
		return fmt.Errorf("Error updating Cosmos Cassandra Keyspace %q (Account: %q) - %+v", id.Name, id.Account, err)
	}

	db := documentdb.CassandraKeyspaceCreateUpdateParameters{
		CassandraKeyspaceCreateUpdateProperties: &documentdb.CassandraKeyspaceCreateUpdateProperties{
			Resource: &documentdb.CassandraKeyspaceResource{
				ID: &id.Name,
			},
			Options: &documentdb.CreateUpdateOptions{},
		},
	}

	future, err := client.CreateUpdateCassandraKeyspace(ctx, id.ResourceGroup, id.Account, id.Name, db)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for Cosmos Cassandra Keyspace %q (Account: %q): %+v", id.ResourceGroup, id.Account, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting on create/update future for Cosmos Cassandra Keyspace %q (Account: %q): %+v", id.ResourceGroup, id.Account, err)
	}

	if common.HasThroughputChange(d) {
		throughputParameters := common.ExpandCosmosDBThroughputSettingsUpdateParameters(d)
		throughputFuture, err := client.UpdateCassandraKeyspaceThroughput(ctx, id.ResourceGroup, id.Account, id.Name, *throughputParameters)
		if err != nil {
			if response.WasNotFound(throughputFuture.Response()) {
				return fmt.Errorf("Error setting Throughput for Cosmos Cassandra Keyspace %q (Account: %q): %+v - "+
					"If the collection has not been created with an initial throughput, you cannot configure it later.", id.Name, id.Account, err)
			}
		}

		if err = throughputFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting on ThroughputUpdate future for Cosmos Cassandra Keyspace %q (Account: %q): %+v", id.Name, id.Account, err)
		}
	}

	return resourceArmCosmosDbCassandraKeyspaceRead(d, meta)
}

func resourceArmCosmosDbCassandraKeyspaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CassandraClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CassandraKeyspaceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetCassandraKeyspace(ctx, id.ResourceGroup, id.Account, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Cosmos Cassandra Keyspace %q (Account: %q) - removing from state", id.Name, id.Account)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Cosmos Cassandra Keyspace %q (Account: %q): %+v", id.Name, id.Account, err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("account_name", id.Account)
	if props := resp.CassandraKeyspaceGetProperties; props != nil {
		if res := props.Resource; res != nil {
			d.Set("name", res.ID)
		}
	}

	throughputResp, err := client.GetCassandraKeyspaceThroughput(ctx, id.ResourceGroup, id.Account, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(throughputResp.Response) {
			return fmt.Errorf("Error reading Throughput on Cosmos Cassandra Keyspace %q (Account: %q): %+v", id.Name, id.Account, err)
		} else {
			d.Set("throughput", nil)
			d.Set("autoscale_settings", nil)
		}
	} else {
		common.SetResourceDataThroughputFromResponse(throughputResp, d)
	}

	return nil
}

func resourceArmCosmosDbCassandraKeyspaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CassandraClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CassandraKeyspaceID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.DeleteCassandraKeyspace(ctx, id.ResourceGroup, id.Account, id.Name)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error deleting Cosmos Cassandra Keyspace %q (Account: %q): %+v", id.Name, id.Account, err)
		}
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting on delete future for Cosmos Cassandra Keyspace %q (Account: %q): %+v", id.Name, id.Account, err)
	}

	return nil
}

package subscription

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/prusya/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/prusya/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/prusya/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmSubscription() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmSubscriptionRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"subscription_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"location_placement_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"quota_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"spending_limit": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client)
	groupClient := client.Subscription.Client
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	subscriptionId := d.Get("subscription_id").(string)
	if subscriptionId == "" {
		subscriptionId = client.Account.SubscriptionId
	}

	resp, err := groupClient.Get(ctx, subscriptionId)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Subscription %q was not found", subscriptionId)
		}

		return fmt.Errorf("Error reading Subscription: %+v", err)
	}

	d.SetId(*resp.ID)
	d.Set("subscription_id", resp.SubscriptionID)
	d.Set("display_name", resp.DisplayName)
	d.Set("tenant_id", resp.TenantID)
	d.Set("state", resp.State)
	if resp.SubscriptionPolicies != nil {
		d.Set("location_placement_id", resp.SubscriptionPolicies.LocationPlacementID)
		d.Set("quota_id", resp.SubscriptionPolicies.QuotaID)
		d.Set("spending_limit", resp.SubscriptionPolicies.SpendingLimit)
	}

	return nil
}

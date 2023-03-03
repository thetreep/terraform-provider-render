package data_sources

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jackall3n/render-go"
	"github.com/jackall3n/terraform-provider-render/render/types"
)

func dataSourceOwner() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOwnerRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOwnerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*types.Context).Client

	var diags diag.Diagnostics

	email := d.Get("email").(string)

	tflog.Debug(ctx, fmt.Sprintf("getting owner: %s", email))

	response, err := c.GetOwnersWithResponse(ctx, &render.GetOwnersParams{
		Email: &[]string{email},
	})

	if err != nil {
		return diag.FromErr(err)
	}

	owner := (*response.JSON200)[0].Owner

	items := map[string]interface{}{
		"id":    owner.Id,
		"email": owner.Email,
		"name":  owner.Name,
		"type":  owner.Type,
	}

	for key, value := range items {
		if err := d.Set(key, value); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(owner.Id)

	return diags
}

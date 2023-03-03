package resources

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jackall3n/render-go"
	"github.com/jackall3n/terraform-provider-render/render/types"
	"github.com/jackall3n/terraform-provider-render/render/utils"
	"net/http"
)

func resourceService() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceCreate,
		ReadContext:   resourceServiceRead,
		UpdateContext: resourceServiceUpdate,
		DeleteContext: resourceServiceDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo": {
				Type:     schema.TypeString,
				Required: true,
			},
			"branch": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auto_deploy": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"web_service_details": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Describes the Service being deployed.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"env": {
							Type:     schema.TypeString,
							Required: true,
						},
						"region": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"instances": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"plan": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"health_check_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"pull_request_previews_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"native": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"build_command": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"start_command": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceServiceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*types.Context)

	ownerId, err := getOwnerId(c, d)

	if err != nil {
		return diag.FromErr(err)
	}

	service := render.ServicePOST{
		Name: d.Get("name").(string),
		Repo: d.Get("repo").(string),
		Type: utils.ParseServiceType(d.Get("type").(string)),
		//AutoDeploy: d.Get("auto_deploy"),
		OwnerId: *ownerId,
	}

	serviceDetails, err := transformServiceDetails(ctx, d)

	if err != nil {
		return diag.FromErr(err)
	}

	if serviceDetails != nil {
		service.ServiceDetails = serviceDetails
	}

	tflog.Debug(ctx, "creating service", utils.ToJson(service))

	response, err := c.Client.CreateServiceWithResponse(ctx, service)

	if err != nil {
		return diag.FromErr(err)
	}

	if response.StatusCode() != http.StatusCreated {
		return diag.Errorf("error creating service: %s", string(response.Body))
	}

	tflog.Debug(ctx, "Created service: "+response.Status())

	s := response.JSON201.Service

	d.SetId(*s.Id)

	return resourceServiceRead(ctx, d, m)
}

func resourceServiceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*types.Context)

	var diags diag.Diagnostics

	id := d.Id()

	s, err := c.Client.GetServiceWithResponse(ctx, id)

	service := s.JSON200

	if err != nil {
		return diag.FromErr(err)
	}

	properties := map[string]interface{}{
		"id":   service.Id,
		"name": service.Name,
		"type": service.Type,
		"repo": service.Repo,
	}

	for key, value := range properties {
		if err := d.Set(key, value); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceServiceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*types.Context)

	id := d.Id()

	name := d.Get("name").(string)

	service := render.ServicePATCH{
		Name: &name,
	}

	_, err := c.Client.UpdateService(ctx, id, service)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceServiceRead(ctx, d, m)
}

func resourceServiceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*types.Context)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	id := d.Id()

	_, err := c.Client.DeleteService(ctx, id)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func getOwnerId(c *types.Context, d *schema.ResourceData) (*string, error) {
	raw, ok := d.GetOk("owner")

	if !ok {
		if c.Owner == nil {
			return nil, fmt.Errorf("'owner' is required if a global email is not set")
		}

		return &c.Owner.Id, nil
	}

	id := raw.(string)

	return &id, nil

}

package resources

//
//import (
//	"context"
//	"fmt"
//	"github.com/hashicorp/terraform-plugin-log/tflog"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
//	"github.com/jackall3n/render-go"
//	"github.com/jackall3n/terraform-provider-render/render/types"
//
//	"github.com/hashicorp/terraform-plugin-framework/resource"
//	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
//	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
//	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
//	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
//	"net/http"
//	"time"
//)
//
//func resourceServiceEnvironment() *schema.Resource {
//	return &schema.Resource{
//		CreateContext: resourceServiceEnvironmentCreate,
//		ReadContext:   resourceServiceEnvironmentRead,
//		UpdateContext: resourceServiceEnvironmentUpdate,
//		DeleteContext: resourceServiceEnvironmentDelete,
//		Schema: map[string]*schema.Schema{
//			"service": {
//				Type:     schema.TypeString,
//				Required: true,
//			},
//			"variables": {
//				Type:     schema.TypeSet,
//				Required: true,
//				Elem: &schema.Resource{
//					Schema: map[string]*schema.Schema{
//						"key": {
//							Type:     schema.TypeString,
//							Required: true,
//						},
//						"generated": {
//							Type:     schema.TypeString,
//							Optional: true,
//							Default:  false,
//						},
//						"value": {
//							Type:     schema.TypeString,
//							Optional: true,
//						},
//					},
//				},
//			},
//		},
//	}
//}
//
//func resourceServiceEnvironmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
//	c := m.(*types.Context)
//
//	serviceId := d.Get("service").(string)
//
//	d.SetId(time.Now().String())
//
//	body := render.UpdateEnvVarsForServiceJSONRequestBody{}
//
//	for key, value := range d.Get("variables").(map[string]interface{}) {
//		tflog.Debug(ctx, fmt.Sprintf("%s %s", key, value))
//
//		item := render.EnvVarsPATCH_Item{}
//
//		v := value.(map[string]interface{})
//
//		if v["generated"] == true {
//			item.FromEnvVarKeyGenerateValue(render.EnvVarKeyGenerateValue{
//				Key: key,
//			})
//		} else {
//			item.FromEnvVarKeyValue(render.EnvVarKeyValue{
//				Key:   key,
//				Value: v["value"].(string),
//			})
//		}
//
//		body = append(body, item)
//	}
//
//	response, err := c.Client.UpdateEnvVarsForServiceWithResponse(ctx, serviceId, body)
//
//	if err != nil {
//		return diag.FromErr(err)
//	}
//
//	if response.StatusCode() != http.StatusNoContent {
//		return diag.Errorf("error updating service env vars: %s %s", response.StatusCode(), string(response.Body))
//	}
//
//	tflog.Debug(ctx, "updated service env vars: "+response.Status())
//
//	return resourceServiceEnvironmentRead(ctx, d, m)
//}
//
//func resourceServiceEnvironmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
//	c := m.(*types.Context)
//
//	var diags diag.Diagnostics
//
//	id := d.Id()
//
//	s, err := c.Client.GetServiceWithResponse(ctx, id)
//
//	service := s.JSON200
//
//	if err != nil {
//		return diag.FromErr(err)
//	}
//
//	properties := map[string]interface{}{
//		"id":   service.Id,
//		"name": service.Name,
//		"type": service.Type,
//		"repo": service.Repo,
//	}
//
//	for key, value := range properties {
//		if err := d.Set(key, value); err != nil {
//			return diag.FromErr(err)
//		}
//	}
//
//	return diags
//}
//
//func resourceServiceEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
//	c := m.(*types.Context)
//
//	serviceId := d.Get("service").(string)
//
//	body := render.UpdateEnvVarsForServiceJSONRequestBody{}
//
//	for key, value := range d.Get("variables").(map[string]interface{}) {
//		tflog.Debug(ctx, fmt.Sprintf("%s %s", key, value))
//
//		item := render.EnvVarsPATCH_Item{}
//
//		v := value.(map[string]interface{})
//
//		if v["generated"] == true {
//			item.FromEnvVarKeyGenerateValue(render.EnvVarKeyGenerateValue{
//				Key: key,
//			})
//		} else {
//			item.FromEnvVarKeyValue(render.EnvVarKeyValue{
//				Key:   key,
//				Value: v["value"].(string),
//			})
//		}
//
//		body = append(body, item)
//	}
//
//	response, err := c.Client.UpdateEnvVarsForServiceWithResponse(ctx, serviceId, body)
//
//	if err != nil {
//		return diag.FromErr(err)
//	}
//
//	if response.StatusCode() != http.StatusNoContent {
//		return diag.Errorf("error updating service env vars: %s %s", response.StatusCode(), string(response.Body))
//	}
//
//	tflog.Debug(ctx, "updated service env vars: "+response.Status())
//
//	return resourceServiceEnvironmentRead(ctx, d, m)
//}
//
//func resourceServiceEnvironmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
//	c := m.(*types.Context)
//
//	// Warning or errors can be collected in a slice type
//	var diags diag.Diagnostics
//
//	serviceId := d.Get("service").(string)
//
//	response, err := c.Client.UpdateEnvVarsForServiceWithResponse(ctx, serviceId, render.UpdateEnvVarsForServiceJSONRequestBody{})
//
//	if err != nil {
//		return diag.FromErr(err)
//	}
//
//	if response.StatusCode() != http.StatusNoContent {
//		return diag.Errorf("error deleting service env vars: %s %s", response.StatusCode(), string(response.Body))
//	}
//
//	tflog.Debug(ctx, "deleted service env vars: "+response.Status())
//
//	d.SetId("")
//
//	return diags
//}

package render

import (
	"context"
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jackall3n/render-go"
	"github.com/jackall3n/terraform-provider-render/render/types"
)

var host = "https://api.render.com/v1"

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiKey := d.Get("api_key").(string)
	email := getEmail(ctx, d)

	tflog.Debug(ctx, fmt.Sprintf("email: %s", email))

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	bearer, _ := securityprovider.NewSecurityProviderBearerToken(apiKey)
	client, _ := render.NewClientWithResponses(host, render.WithRequestEditorFn(bearer.Intercept))

	return createContext(ctx, client, email, diags)
}

func getEmail(ctx context.Context, d *schema.ResourceData) string {
	email, ok := d.GetOk("email")

	tflog.Debug(ctx, fmt.Sprintf("%s %t", email, ok))

	if ok {
		return email.(string)
	}

	return ""
}

func createContext(ctx context.Context, client *render.ClientWithResponses, email string, diags diag.Diagnostics) (interface{}, diag.Diagnostics) {
	c := &types.Context{Client: client}

	if email == "" {
		return c, diags
	}

	owner, diags := getOwner(ctx, client, email, diags)

	if owner == nil {
		return nil, diags
	}

	c.Owner = owner

	return c, diags
}

func getOwner(ctx context.Context, client *render.ClientWithResponses, email string, diags diag.Diagnostics) (*render.Owner, diag.Diagnostics) {
	tflog.Debug(ctx, fmt.Sprintf("getting owners with email: %s", email))

	response, err := client.GetOwnersWithResponse(ctx, &render.GetOwnersParams{
		Email: &[]string{email},
	})

	if err != nil {
		return nil, diag.FromErr(err)
	}

	owner := (*response.JSON200)[0].Owner

	return owner, diags
}

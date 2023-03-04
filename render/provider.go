package render

import (
	"context"
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jackall3n/render-go"
	"github.com/jackall3n/terraform-provider-render/render/datasources"
	"github.com/jackall3n/terraform-provider-render/render/resources"
	"os"
)

type renderProvider struct{}

// New instantiates a new instance of a render terraform provider.
func New() provider.Provider {
	return &renderProvider{}
}

func (p *renderProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "render"
}

func (p *renderProvider) Schema(_ context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `
The Render provider is used to interact with resources supported by Render.
        `,
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
			"email": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (p *renderProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.ServiceResource,
		resources.ServiceEnvironmentResource,
	}
}

func (p *renderProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.OwnerDataSource,
	}
}

type ProviderData struct {
	APIKey types.String `tfsdk:"api_key"`
	Email  types.String `tfsdk:"email"`
}

func (p *renderProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config ProviderData

	diags := req.Config.Get(ctx, &config)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// User must provide an api_token to the provider
	var apiKey string
	var email string

	if config.APIKey.IsUnknown() {
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as api_key",
		)
		return
	}

	if config.APIKey.IsNull() {
		apiKey = os.Getenv("RENDER_API_KEY")
	} else {
		apiKey = config.APIKey.ValueString()
	}

	if apiKey == "" {
		resp.Diagnostics.AddError(
			"Unable to find api_key",
			"api_key cannot be an empty string",
		)
		return
	}

	if config.Email.IsNull() {
		email = os.Getenv("RENDER_EMAIL")
	} else {
		email = config.Email.ValueString()
	}

	tflog.Debug(ctx, fmt.Sprintf("email: %s", email))

	bearer, _ := securityprovider.NewSecurityProviderBearerToken(apiKey)
	client, _ := render.NewClientWithResponses(host, render.WithRequestEditorFn(bearer.Intercept))

	c := createContext(ctx, client, email)

	resp.DataSourceData = c
	resp.ResourceData = c
}

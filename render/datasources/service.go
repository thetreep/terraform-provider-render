package datasources

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jackall3n/render-go"
	"github.com/jackall3n/terraform-provider-render/render/models"
	"github.com/jackall3n/terraform-provider-render/render/types"
	"net/http"
)

func ServiceDataSource() datasource.DataSource {
	return &serviceDataSource{}
}

type serviceDataSource struct {
	client  *render.ClientWithResponses
	context *types.Context
}

func (d *serviceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service"
}

func (d *serviceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	ctx, ok := req.ProviderData.(*types.Context)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *render.ClientWithResponses, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.context = ctx
	d.client = ctx.Client
}

// Schema returns the schema information for an service data source
func (_ *serviceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `Provides information about an existing Service resource.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"email": schema.StringAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"type": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *serviceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.Service

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := d.client.GetServiceWithResponse(ctx, data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("failed to get service", err.Error())
		return
	}

	if response.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("failed to get service", fmt.Sprintf("%s %s", response.Status(), string(response.Body)))
		return
	}

	service := response.JSON200

	result := data.FromResponse(*service)

	tflog.Trace(ctx, "read service", map[string]interface{}{
		"id":   result.ID.ValueString(),
		"name": result.Name.ValueString(),
	})

	diags = resp.State.Set(ctx, result)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

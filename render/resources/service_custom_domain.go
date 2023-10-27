package resources

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jackall3n/render-go"
	"github.com/jackall3n/terraform-provider-render/render/models"
	"github.com/jackall3n/terraform-provider-render/render/types"
	"github.com/jackall3n/terraform-provider-render/render/utils"
	"io"
)

func ServiceCustomDomainResource() resource.Resource {
	return &serviceCustomDomainResource{}
}

type serviceCustomDomainResource struct {
	client  *render.ClientWithResponses
	context *types.Context
}

func (r *serviceCustomDomainResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_custom_domain"
}

func (r *serviceCustomDomainResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ctx, ok := req.ProviderData.(*types.Context)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.ClientWithResponses, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.context = ctx
	r.client = ctx.Client
}

func (r *serviceCustomDomainResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `Provider for service custom domain resource`,
		Attributes: map[string]schema.Attribute{
			"service_id":  schema.StringAttribute{Required: true},
			"domain_name": schema.StringAttribute{Required: true},
		},
	}
}

func (r *serviceCustomDomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.ServiceCustomDomain

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	customDomainJSONBody := render.CreateCustomDomainJSONRequestBody{
		Name: plan.DomainName.ValueString(),
	}

	tflog.Debug(ctx, "creating custom domain", utils.ToJson(map[string]interface{}{
		"service_id":  plan.ServiceID.ValueString(),
		"domain_name": customDomainJSONBody,
	}))

	response, err := r.client.CreateCustomDomain(ctx, plan.ServiceID.ValueString(), customDomainJSONBody)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update service variables",
			fmt.Sprintf("Could not update service variables %s, unexpected error: %s",
				err.Error(),
				err,
			),
		)
		return
	}

	bytes, err := io.ReadAll(response.Body)

	if response.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Failed to create custom domain",
			fmt.Sprintf("Could not create custom domain %s, unexpected error: %s",
				response.Status,
				string(bytes),
			),
		)
		return
	}

	tflog.Debug(ctx, "Created custom domain "+response.Status, map[string]interface{}{
		"r": string(bytes),
	})

	resp.State.Set(ctx, plan)
}

func (r *serviceCustomDomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	//
}

func (r *serviceCustomDomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//
}

func (r *serviceCustomDomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	//
}

package modifiers

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jackall3n/render-go"
)

func OwnerDefault(owner *render.Owner) planmodifier.String {
	return &ownerDefaultPlanModifier{
		Owner: owner,
	}
}

type ownerDefaultPlanModifier struct {
	Owner *render.Owner
}

var _ planmodifier.String = (*ownerDefaultPlanModifier)(nil)

func (apm *ownerDefaultPlanModifier) Description(_ context.Context) string {
	return "Owner default modifier"
}

func (apm *ownerDefaultPlanModifier) MarkdownDescription(_ context.Context) string {
	return "Owner default modifier"
}

func (apm *ownerDefaultPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, res *planmodifier.StringResponse) {
	tflog.Debug(ctx, "PlanModifyString")

	// If the attribute configuration is not null, we are done here
	if !req.ConfigValue.IsNull() {
		return
	}

	// If the attribute plan is "known" and "not null", then a previous plan modifier in the sequence
	// has already been applied, and we don't want to interfere.
	if !req.PlanValue.IsUnknown() && !req.PlanValue.IsNull() {
		return
	}

	if apm.Owner == nil {
		res.Diagnostics.AddError("failed to set an owner", "'owner' is required if a global email is not set")
		return
	}

	res.PlanValue = types.StringValue(apm.Owner.Id)
}

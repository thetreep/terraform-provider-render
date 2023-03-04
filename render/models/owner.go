package models

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jackall3n/render-go"
)

type Owner struct {
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Email types.String `tfsdk:"email"`
	Type  types.String `tfsdk:"type"`
}

func (o Owner) FromResponse(response render.Owner) Owner {
	return Owner{
		ID:    types.StringValue(response.Id),
		Name:  types.StringValue(*response.Name),
		Type:  types.StringValue(fmt.Sprintf("%s", response.Type)),
		Email: o.Email,
	}
}

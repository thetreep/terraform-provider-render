package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ServiceCustomDomain struct {
	ServiceID  types.String `tfsdk:"service_id"`
	DomainName types.String `tfsdk:"domain_name"`
}

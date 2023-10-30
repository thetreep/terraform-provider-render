package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jackall3n/render-go"
)

type ServiceCustomDomain struct {
	ServiceID  types.String `tfsdk:"service_id"`
	DomainName types.String `tfsdk:"domain_name"`
}

func (s ServiceCustomDomain) FromResponse(response render.CustomDomain) ServiceCustomDomain {
	return ServiceCustomDomain{
		ServiceID:  s.ServiceID, //TODO: pass service ID from response
		DomainName: fromStringOptional(response.Name),
	}
}

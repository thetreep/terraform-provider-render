package resources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ProviderResource() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"render_service": resourceService(),
	}
}

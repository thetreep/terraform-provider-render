package data_sources

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ProviderDataSource() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"render_owner": dataSourceOwner(),
	}
}

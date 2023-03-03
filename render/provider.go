package render

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jackall3n/terraform-provider-render/render/data_sources"
	"github.com/jackall3n/terraform-provider-render/render/resources"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ConfigureContextFunc: providerConfigure,
		Schema:               providerSchema(),
		DataSourcesMap:       data_sources.ProviderDataSource(),
		ResourcesMap:         resources.ProviderResource(),
	}
}

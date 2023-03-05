package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jackall3n/render-go"
)

type ServiceEnvironment struct {
	Service   types.String                 `tfsdk:"service"`
	Variables []ServiceEnvironmentVariable `tfsdk:"variables"`
}

type ServiceEnvironmentVariable struct {
	Key       types.String `tfsdk:"key"`
	Value     types.String `tfsdk:"value"`
	Generated types.Bool   `tfsdk:"generated"`
}

func (s ServiceEnvironment) FromResponse(variables []ServiceEnvironmentVariable) ServiceEnvironment {
	return ServiceEnvironment{
		Service:   s.Service,
		Variables: variables,
	}
}

func (v ServiceEnvironmentVariable) FromResponse(response render.EnvVar) ServiceEnvironmentVariable {
	return ServiceEnvironmentVariable{
		Key:   types.StringValue(response.Key),
		Value: fromStringOptional(&response.Value),
	}
}

func (v ServiceEnvironmentVariable) ToEnvVarsItemPATCH() (*render.EnvVarsPATCH_Item, error) {
	item := render.EnvVarsPATCH_Item{}

	if v.Generated.ValueBool() {
		err := item.FromEnvVarKeyGenerateValue(render.EnvVarKeyGenerateValue{
			Key:           v.Key.ValueString(),
			GenerateValue: "true",
		})

		if err != nil {
			return nil, err
		}

		return &item, nil
	}

	err := item.FromEnvVarKeyValue(render.EnvVarKeyValue{
		Key:   v.Key.ValueString(),
		Value: v.Value.ValueString(),
	})

	if err != nil {
		return nil, err
	}

	return &item, nil
}

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jackall3n/render-go"
	"github.com/jackall3n/terraform-provider-render/render/utils"
)

type Service struct {
	ID                types.String       `tfsdk:"id"`
	Name              types.String       `tfsdk:"name"`
	Type              types.String       `tfsdk:"type"`
	Repo              types.String       `tfsdk:"repo"`
	Branch            types.String       `tfsdk:"branch"`
	Owner             types.String       `tfsdk:"owner"`
	AutoDeploy        types.Bool         `tfsdk:"auto_deploy"`
	WebServiceDetails *WebServiceDetails `tfsdk:"web_service_details"`
}

type WebServiceDetails struct {
	Env                        types.String             `tfsdk:"env"`
	Region                     types.String             `tfsdk:"region"`
	Plan                       types.String             `tfsdk:"plan"`
	PullRequestPreviewsEnabled types.Bool               `tfsdk:"pull_request_previews_enabled"`
	HealthCheckPath            types.String             `tfsdk:"health_check_path"`
	Native                     *WebServiceDetailsNative `tfsdk:"native"`
}

type WebServiceDetailsNative struct {
	BuildCommand types.String `tfsdk:"build_command"`
	StartCommand types.String `tfsdk:"start_command"`
}

func (s Service) FromResponse(response render.Service) Service {
	service := Service{
		ID:     fromStringOptional(response.Id),
		Name:   fromStringOptional(response.Name),
		Type:   fromServiceType(response.Type),
		Repo:   fromStringOptional(response.Repo),
		Branch: fromStringOptional(response.Branch),
		Owner:  fromStringOptional(response.OwnerId),
	}

	if *response.Type == render.WebService {
		webServiceDetails, _ := response.ServiceDetails.AsWebServiceDetails()

		service.WebServiceDetails = &WebServiceDetails{
			Region:          fromRegion(webServiceDetails.Region),
			Env:             fromServiceEnv(webServiceDetails.Env),
			Plan:            fromStringOptional(webServiceDetails.Plan),
			HealthCheckPath: fromStringOptional(webServiceDetails.HealthCheckPath),
		}

		native, err := webServiceDetails.EnvSpecificDetails.AsNativeEnvironmentDetails()

		if err == nil {
			service.WebServiceDetails.Native = &WebServiceDetailsNative{
				BuildCommand: fromStringOptional(native.BuildCommand),
				StartCommand: fromStringOptional(native.StartCommand),
			}
		}
	}

	return service
}

func (s Service) ToServicePOST(ownerId string) (*render.ServicePOST, error) {
	service := render.ServicePOST{
		Type:    render.ServiceType(s.Type.ValueString()),
		Name:    s.Name.ValueString(),
		Repo:    s.Repo.ValueString(),
		Branch:  stringOptional(s.Branch),
		OwnerId: ownerId,
	}

	if s.WebServiceDetails != nil {
		serviceDetails := render.WebServiceDetailsPOST{}

		err := utils.Struct(toWebServiceDetails(s.WebServiceDetails), &serviceDetails)

		if err != nil {
			return nil, err
		}

		details := render.ServicePOST_ServiceDetails{}
		details.FromWebServiceDetailsPOST(serviceDetails)

		service.ServiceDetails = &details
	}

	return &service, nil
}

func (s Service) ToServicePATCH() (*render.ServicePATCH, error) {
	service := render.ServicePATCH{
		Name:   stringOptional(s.Name),
		Branch: stringOptional(s.Branch),
	}

	if s.WebServiceDetails != nil {
		serviceDetails := render.WebServiceDetailsPATCH{}

		err := utils.Struct(toWebServiceDetails(s.WebServiceDetails), &serviceDetails)

		if err != nil {
			return nil, err
		}

		details := render.ServicePATCH_ServiceDetails{}
		details.FromWebServiceDetailsPATCH(serviceDetails)

		service.ServiceDetails = &details
	}

	return &service, nil
}

func toWebServiceDetails(webServiceDetails *WebServiceDetails) map[string]interface{} {
	details := map[string]interface{}{
		"region":          stringOptional(webServiceDetails.Region),
		"env":             stringOptional(webServiceDetails.Env),
		"plan":            stringOptional(webServiceDetails.Plan),
		"healthCheckPath": stringOptional(webServiceDetails.HealthCheckPath),
	}

	if webServiceDetails.Native != nil {
		native := map[string]interface{}{
			"buildCommand": webServiceDetails.Native.BuildCommand.ValueString(),
			"startCommand": webServiceDetails.Native.StartCommand.ValueString(),
		}

		details["envSpecificDetails"] = native
	}

	return details
}

func stringOptional(str types.String) *string {
	if str.IsNull() {
		return nil
	}

	value := str.ValueString()

	return &value
}

func int64Optional(num types.Int64) *int64 {
	if num.IsNull() {
		return nil
	}

	value := num.ValueInt64()

	return &value
}

func fromIntOptional(num *int) types.Int64 {
	if num == nil {
		return types.Int64Null()
	}

	return types.Int64Value(int64(*num))
}

func fromStringOptional(str *string) types.String {
	if str == nil {
		return types.StringNull()
	}

	return types.StringValue(*str)
}

func fromServiceType(t *render.ServiceType) types.String {
	if t == nil {
		return types.StringNull()
	}

	return types.StringValue(string(*t))
}

func fromServiceEnv(e *render.ServiceEnv) types.String {
	if e == nil {
		return types.StringNull()
	}

	return types.StringValue(string(*e))
}

func fromRegion(r *render.Region) types.String {
	if r == nil {
		return types.StringNull()
	}

	return types.StringValue(string(*r))
}

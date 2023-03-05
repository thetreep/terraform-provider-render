package models

import (
	"fmt"
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
	Env    types.String `tfsdk:"env"`
	Region types.String `tfsdk:"region"`
	//Instances                  types.Int64              `tfsdk:"instances"`
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
			Region: fromRegion(webServiceDetails.Region),
			Env:    fromServiceEnv(webServiceDetails.Env),
			//Instances:       fromIntOptional(webServiceDetails.NumInstances),
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
	serviceType := render.ServiceType(s.Type.ValueString())

	service := render.ServicePOST{
		Type:    serviceType,
		Name:    s.Name.ValueString(),
		Repo:    s.Repo.ValueString(),
		Branch:  stringOptionalNil(s.Branch),
		OwnerId: ownerId,
	}

	//js, _ := json.Marshal(service)
	//
	//return nil, fmt.Errorf("value: %s", js)

	if serviceType == render.WebService {
		if s.WebServiceDetails == nil {
			return nil, fmt.Errorf("'web_service_details' is required for service type 'web_service'")
		}

		serviceDetails := render.WebServiceDetailsPOST{}

		webServiceDetails, err := toWebServiceDetails(s.WebServiceDetails)

		if err != nil {
			return nil, err
		}

		err = utils.Struct(webServiceDetails, &serviceDetails)

		if err != nil {
			return nil, err
		}

		details := render.ServicePOST_ServiceDetails{}
		details.FromWebServiceDetailsPOST(serviceDetails)

		service.ServiceDetails = &details
	} else if s.WebServiceDetails != nil {
		return nil, fmt.Errorf("'web_service_details' can only be used with the service type 'web_service'")
	}

	return &service, nil
}

func (s Service) ToServicePATCH() (*render.ServicePATCH, error) {
	serviceType := render.ServiceType(s.Type.ValueString())

	service := render.ServicePATCH{
		Name:   stringOptional(s.Name),
		Branch: stringOptionalNil(s.Branch),
	}

	if serviceType == render.WebService {
		if s.WebServiceDetails == nil {
			return nil, fmt.Errorf("'web_service_details' is required for service type 'web_service'")
		}

		serviceDetails := render.WebServiceDetailsPATCH{}

		webServiceDetails, err := toWebServiceDetails(s.WebServiceDetails)

		if err != nil {
			return nil, err
		}

		err = utils.Struct(webServiceDetails, &serviceDetails)

		if err != nil {
			return nil, err
		}

		details := render.ServicePATCH_ServiceDetails{}
		details.FromWebServiceDetailsPATCH(serviceDetails)

		service.ServiceDetails = &details
	} else if s.WebServiceDetails != nil {
		return nil, fmt.Errorf("'web_service_details' can only be used with the service type 'web_service'")
	}

	return &service, nil
}

func toWebServiceDetails(webServiceDetails *WebServiceDetails) (map[string]interface{}, error) {
	details := map[string]interface{}{
		"region": stringOptionalNil(webServiceDetails.Region),
		"env":    stringOptional(webServiceDetails.Env),
		//"numInstances":    int64Optional(webServiceDetails.Instances),
		"plan":            stringOptionalNil(webServiceDetails.Plan),
		"healthCheckPath": stringOptional(webServiceDetails.HealthCheckPath),
	}

	if webServiceDetails.Native != nil {
		native := map[string]interface{}{
			"buildCommand": webServiceDetails.Native.BuildCommand.ValueString(),
			"startCommand": webServiceDetails.Native.StartCommand.ValueString(),
		}

		details["envSpecificDetails"] = native
	}

	return details, nil
}

func stringOptional(str types.String) *string {
	if str.IsNull() {
		return nil
	}

	value := str.ValueString()

	return &value
}

func stringOptionalNil(str types.String) *string {
	if str.IsNull() || str.ValueString() == "" {
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

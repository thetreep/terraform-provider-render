package resources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jackall3n/render-go"
	"github.com/jackall3n/terraform-provider-render/render/utils"
	"log"
	"reflect"
)

func transformServiceDetails(ctx context.Context, d *schema.ResourceData) (*render.ServicePOST_ServiceDetails, error) {
	if raw, ok := d.GetOk("web_service_details"); ok {
		return transformWebServiceDetails(ctx, utils.GetBlock(raw))
	}

	return nil, nil
}

func transformWebServiceDetails(ctx context.Context, serviceDetails map[string]interface{}) (*render.ServicePOST_ServiceDetails, error) {
	env := utils.ParseServiceEnv(serviceDetails["env"].(string))

	details := render.ServicePOST_ServiceDetails{}

	err := details.FromWebServiceDetailsPOST(render.WebServiceDetailsPOST{
		Env:                env,
		EnvSpecificDetails: transformWebServiceEnvSpecificDetails(ctx, serviceDetails),
	})

	if err != nil {
		return nil, err
	}

	return &details, err
}

func transformWebServiceEnvSpecificDetails(ctx context.Context, value map[string]interface{}) *render.WebServiceDetailsPOST_EnvSpecificDetails {
	tflog.Debug(ctx, "transformEnvSpecificDetails", value)

	envDetails := render.WebServiceDetailsPOST_EnvSpecificDetails{}

	if property := reflect.ValueOf(value["native"]); property.IsValid() {
		native := utils.GetBlock(value["native"])

		err := envDetails.FromNativeEnvironmentDetailsPOST(render.NativeEnvironmentDetailsPOST{
			BuildCommand: utils.TryString(native, "build_command"),
			StartCommand: utils.TryString(native, "start_command"),
		})

		if err != nil {
			log.Fatal(err)
		}
	}

	if property := reflect.ValueOf(value["docker"]); property.IsValid() {
		docker := utils.GetBlock(value["docker"])

		err := envDetails.FromDockerDetailsPOST(render.DockerDetailsPOST{
			DockerCommand:  utils.TryStringRef(docker, "docker_command"),
			DockerContext:  utils.TryStringRef(docker, "docker_context"),
			DockerfilePath: utils.TryStringRef(docker, "dockerfile_path"),
		})

		if err != nil {
			log.Fatal(err)
		}
	}

	return &envDetails
}

package utils

import (
	"encoding/json"
	"github.com/jackall3n/render-go"
	"reflect"
)

func GetBlock(value interface{}) map[string]interface{} {
	return value.([]interface{})[0].(map[string]interface{})
}

func TryString(value map[string]interface{}, name string) string {
	if property := reflect.ValueOf(value[name]); property.IsValid() {
		val := value[name].(string)
		return val
	}

	return ""
}

func TryStringRef(value map[string]interface{}, name string) *string {
	result := TryString(value, name)

	return &result
}

func ParseServiceEnv(env string) render.ServiceEnv {
	return serviceEnvMap[env]
}

func ParseServiceType(serviceType string) render.ServiceType {
	return serviceTypeMap[serviceType]
}

func ParseRegion(region string) render.Region {
	return regionMap[region]
}

var (
	serviceTypeMap = map[string]render.ServiceType{
		"web_service":       render.WebService,
		"private_service":   render.PrivateService,
		"background_worker": render.BackgroundWorker,
		"cron_job":          render.CronJob,
		"static_site":       render.StaticSite,
	}

	serviceEnvMap = map[string]render.ServiceEnv{
		"docker": render.Docker,
		"elixir": render.Elixir,
		"go":     render.Go,
		"node":   render.Node,
		"python": render.Python,
		"ruby":   render.Ruby,
		"rust":   render.Rust,
	}

	regionMap = map[string]render.Region{
		"frankfurt": render.Frankfurt,
		"oregon": render.Oregon,
	}
)

func ToJson(value interface{}) map[string]interface{} {
	b, _ := json.Marshal(&value)
	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)

	return m
}

type YesNo string

const (
	Yes YesNo = "yes"
	No  YesNo = "no"
)

func ToYesNo(value interface{}) *YesNo {
	if value == nil {
		return nil
	}

	b := value.(bool)

	var result YesNo

	if b {
		result = Yes
	} else {
		result = No
	}

	return &result
}

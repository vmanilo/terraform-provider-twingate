package config

import (
	"fmt"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/utils"
	"strings"
)

// Nprintf - this is a Printf sibling (Nprintf; Named Printf), which handles strings like
// Nprintf("Hello %{target}!", map[string]interface{}{"target":"world"}) == "Hello world!".
// This is particularly useful for generated tests, where we don't want to use Printf,
// since that would require us to generate a very particular ordering of arguments.
func Nprintf(format string, params map[string]interface{}) string {
	for key, val := range params {
		format = strings.ReplaceAll(format, fmt.Sprintf("${%s}", key), fmt.Sprintf("%v", val))
	}

	return format
}

type TerraformResource interface {
	TerraformResource() string
}

func optionalInt(val any) *int {
	if val == nil {
		return nil
	}

	switch t := val.(type) {
	case int:
		return &t
	case *int:
		return t
	default:
		return nil
	}
}

func optionalString(val any) *string {
	if val == nil {
		return nil
	}

	switch t := val.(type) {
	case string:
		return &t
	case *string:
		return t
	default:
		return nil
	}
}

func optionalBool(val any) *bool {
	if val == nil {
		return nil
	}

	switch t := val.(type) {
	case bool:
		return &t
	case *bool:
		return t
	default:
		return nil
	}
}

func GenUsers(count int) []Resource {
	users := make([]Resource, 0, count)

	for i := 0; i < count; i++ {
		users = append(users, NewResourceUser())
	}

	return users
}

func ResourceIDs(resources []Resource) []string {
	return utils.Map(resources, func(r Resource) string {
		return r.TerraformResourceID()
	})
}

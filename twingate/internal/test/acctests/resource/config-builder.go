package resource

import (
	"bytes"
	"fmt"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/utils"
)

type TerraformResource interface {
	TerraformResource() string
}

func collectResourceIDs[T TerraformResource](resources ...T) []string {
	ids := make([]string, 0, len(resources))

	for _, res := range resources {
		ids = append(ids, res.TerraformResource()+".id")
	}

	return ids
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

type wrapper struct {
	str string
}

func (w *wrapper) String() string {
	return w.str
}

func wrap(str string) fmt.Stringer {
	return &wrapper{str: str}
}

func configBuilder(resources ...any) string {
	var list []fmt.Stringer

	for _, r := range resources {
		switch t := r.(type) {
		case fmt.Stringer:
			list = append(list, t)
		case []*User:
			list = append(list, utils.Map(t, func(item *User) fmt.Stringer {
				return item
			})...)
		}
	}

	buff := bytes.NewBufferString("")
	for _, item := range list {
		buff.WriteString(item.String() + "\n")
	}

	return buff.String()
}

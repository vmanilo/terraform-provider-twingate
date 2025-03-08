package config

import (
	"bytes"
	"fmt"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/utils"
)

func Builder(resources ...any) string {
	var list []fmt.Stringer

	for _, r := range resources {
		switch t := r.(type) {
		case fmt.Stringer:
			list = append(list, t)
		case []*ResourceUser:
			list = append(list, utils.Map(t, func(item *ResourceUser) fmt.Stringer {
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

package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewResourceGroup(t *testing.T) {
	tests := []struct {
		name      string
		args      []any
		expRender string
	}{
		{
			name: "Group with name",
			args: []any{"name", "group_1"},
			expRender: `
	resource "twingate_group" "%s" {
	  name = "group_1"
	}
	`},
		{
			name: "Group with security policy",
			args: []any{"name", "group_2", "security_policy_id", "test-policy-id"},
			expRender: `
	resource "twingate_group" "%s" {
	  name = "group_2"
	  security_policy_id = "test-policy-id"
	}
	`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := NewResourceGroup(tt.args...)

			assert.NotEmpty(t, group.ResourceName())
			assert.Equal(t, fmt.Sprintf(tt.expRender, group.ResourceName()), group.String())
		})
	}
}

package client

import (
	"context"

	"github.com/Twingate/terraform-provider-twingate/twingate/internal/client/query"
)

func (client *Client) CreateSerialNumbers(ctx context.Context, serialNumbers []string) ([]string, error) {
	opr := resourceSerialNumbers.create()

	variables := newVars(
		gqlVar(serialNumbers, "serialNumbers"),
	)
	response := query.CreateSerialNumbers{}

	if err := client.mutate(ctx, &response, variables, opr); err != nil {
		return nil, err
	}

	return response.ToModel(), nil
}

func (client *Client) ReadSerialNumbers(ctx context.Context) ([]string, error) {
	opr := resourceSerialNumbers.read()

	response := query.ReadSerialNumbers{}

	if err := client.query(ctx, &response, nil, opr); err != nil {
		return nil, err
	}

	return response.ToModel(), nil
}

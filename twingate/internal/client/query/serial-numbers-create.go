package query

import (
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/utils"
	"github.com/hasura/go-graphql-client"
)

type CreateSerialNumbers struct {
	SerialNumbersEntityResponse `graphql:"serialNumbersCreate(serialNumbers: $serialNumbers)"`
}

type SerialNumbersEntityResponse struct {
	Entities []*gqlSerialNumber
	OkError
}

func (q CreateSerialNumbers) ToModel() []string {
	if q.Entities == nil {
		return nil
	}

	return utils.Map(q.Entities, func(item *gqlSerialNumber) string {
		return item.SerialNumber
	})
}

func (q CreateSerialNumbers) IsEmpty() bool {
	return q.Entities == nil
}

type gqlSerialNumber struct {
	ID           graphql.ID
	SerialNumber string
}

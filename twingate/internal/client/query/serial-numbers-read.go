package query

import "github.com/Twingate/terraform-provider-twingate/twingate/internal/utils"

type ReadSerialNumbers struct {
	SerialNumbersEntityResponse `graphql:"serialNumbers()"`
}

func (q ReadSerialNumbers) IsEmpty() bool {
	return len(q.SerialNumbersEntityResponse.Entities) == 0
}

func (q ReadSerialNumbers) ToModel() []string {
	if q.Entities == nil {
		return nil
	}

	return utils.Map(q.Entities, func(item *gqlSerialNumber) string {
		return item.SerialNumber
	})
}

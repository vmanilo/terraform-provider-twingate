package query

type ReadResourcesByName struct {
	Resources `graphql:"resources(filter: $filter, after: $resourcesEndCursor, first: $pageLimit)"`
}

func (q ReadResourcesByName) IsEmpty() bool {
	return len(q.Edges) == 0
}

type ResourceFilterInput struct {
	Name            *StringFilterOperationInput          `json:"name"`
	Tags            *TagsFilterOperatorInput             `json:"tags"`
	RemoteNetworkID *RemoteNetworkIdFilterOperationInput `json:"remoteNetworkId"`
}

func NewResourceFilterInput(name, filter string, tags map[string]string, remoteNetworkId *string) *ResourceFilterInput {
	return &ResourceFilterInput{
		Name:            NewStringFilterOperationInput(name, filter),
		Tags:            NewTagsFilterOperatorInput(tags),
		RemoteNetworkID: NewRemoteNetworkIdFilterOperationInput(remoteNetworkId),
	}
}

func NewTagsFilterOperatorInput(tags map[string]string) *TagsFilterOperatorInput {
	if len(tags) == 0 {
		return nil
	}

	filter := &TagsFilterOperatorInput{
		And: make([]TagKeyValueFilterInput, 0, len(tags)),
	}

	for key, value := range tags {
		filter.And = append(filter.And, TagKeyValueFilterInput{
			Key: key,
			Value: TagValueFilterInput{
				Eq: &value,
			},
		})
	}

	return filter
}

func NewRemoteNetworkIdFilterOperationInput(remoteNetworkId *string) *RemoteNetworkIdFilterOperationInput {
	if remoteNetworkId == nil {
		return nil
	}

	return &RemoteNetworkIdFilterOperationInput{
		Eq: remoteNetworkId,
	}
}

type TagsFilterOperatorInput struct {
	And []TagKeyValueFilterInput `json:"and"`
}

type TagKeyValueFilterInput struct {
	Key   string              `json:"key"`
	Value TagValueFilterInput `json:"value"`
}

type TagValueFilterInput struct {
	Eq *string `json:"eq"`
}

type RemoteNetworkIdFilterOperationInput struct {
	Eq *string  `json:"eq"`
	In []string `json:"in"`
}

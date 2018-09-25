package api

// Tag is struct for tags on tasks
type Tag struct {
	ID   string
	Name string
}

var tagsCache = make(map[string]string)

func (api *HabiticaAPI) getTags(t Todo) []string {
	tagNames := make([]string, len(t.Tags))

	for i, tagID := range t.Tags {
		tag := api.getTag(tagID)
		tagNames[i] = tag.Name
	}

	return tagNames
}

func (api *HabiticaAPI) getTag(id string) Tag {
	_, exists := tagsCache[id]

	tag := Tag{}
	if !exists {
		api.Get("/tags/"+id, &tag)
		tagsCache[tag.ID] = tag.Name
	}
	tag.ID = id
	tag.Name = tagsCache[id]
	return tag
}

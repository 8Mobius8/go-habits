package api

// Tag is struct for tags on tasks
type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (api *HabiticaAPI) getTagsByTask(t Task) []string {
	tagNames := make([]string, len(t.Tags))

	for i, tagID := range t.Tags {
		tag := api.GetTagByID(tagID)
		tagNames[i] = tag.Name
	}

	return tagNames
}

// GetTagByID will return a tag by querying the habitica server.
// Uses a cache to avoid multiple calls
func (api *HabiticaAPI) GetTagByID(id string) Tag {
	_, exists := fromCache(id)

	tag := Tag{}
	if !exists {
		api.Get("/tags/"+id, &tag)
		updateCache(tag)
	}
	tag, _ = fromCache(id)
	return tag
}

// AddTag creates a new tag on the server using the name given
func (api *HabiticaAPI) AddTag(name string) (Tag, error) {
	var t Tag
	tagName := struct {
		Name string `json:"name"`
	}{name}
	err := api.Post("/tags", &tagName, &t)
	if err != nil {
		return Tag{}, err
	}
	updateCache(t)
	return t, nil
}

// GetTag returns a Tag by name, uses a cache to avoid multiple
// calls to server
func (api *HabiticaAPI) GetTag(name string) (Tag, error) {
	_, exists := fromCache(name)

	tag := Tag{}
	if !exists {
		_, err := api.GetTags()
		if err != nil {
			return Tag{}, err
		}
	}
	tag, _ = fromCache(name)
	return tag, nil
}

// GetTags returns the users list of Tags and updates the tags cache.
func (api *HabiticaAPI) GetTags() ([]Tag, error) {
	tags := []Tag{}
	err := api.Get("/tags", &tags)
	if err != nil {
		return []Tag{}, err
	}

	for _, t := range tags {
		updateCache(t)
	}
	return tags, nil
}

var tagsCache = make(map[string]*Tag)

// ClearTagCache will clear out in-memory Tag cache for client
func (api *HabiticaAPI) ClearTagCache() {
	tagsCache = make(map[string]*Tag)
}

func updateCache(t Tag) {
	newTag := t
	tagsCache[newTag.ID] = &newTag
	tagsCache[newTag.Name] = &newTag
}

func fromCache(key string) (Tag, bool) {
	value, exists := tagsCache[key]
	if !exists {
		return Tag{}, exists
	}
	return *value, exists
}

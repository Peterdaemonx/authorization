package entity

type Metadata struct {
	CurrentPage int
	PageSize    int
	FirstPage   int
	LastPage    bool
}

func CalculateMetadata(fetchedRecords, page, pageSize int) Metadata {
	if fetchedRecords == 0 {
		return Metadata{LastPage: true}
	}

	return Metadata{
		CurrentPage: page,
		PageSize:    pageSize,
		FirstPage:   1,
		LastPage:    pageSize >= fetchedRecords,
	}
}

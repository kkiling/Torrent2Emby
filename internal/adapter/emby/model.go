package emby

type TypeCatalog string

const (
	UnknownTypeCatalog TypeCatalog = "unknown"
	SeasonTypeCatalog  TypeCatalog = "season"
	SeriesTypeCatalog  TypeCatalog = "series"
)

type CatalogInfo struct {
	Path     string
	Name     string
	ID       uint64
	IsFolder bool
	Type     TypeCatalog
}

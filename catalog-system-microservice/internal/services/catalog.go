package catalog

import "log/slog"

type Catalog struct {
	log             *slog.Logger
	catalogProvider CatalogProvider
	appProvider     AppProvider
}

type CatalogProvider interface {
}

type AppProvider interface {
}

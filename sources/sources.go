package blog

import (
	"context"

	"github.com/google/uuid"
	"github.com/maskarb/sources-integration-go/postgres"
)

type Source struct {
	tableName struct{} `pg:"api_sources"`

	ID               uint64    `pg:"source_id"` // Id is automatically detected as primary key
	SourceUUID       uuid.UUID `pg:"source_uuid"`
	Name             string
	AuthHeader       string
	Offset           int
	AccountID        int
	SourceType       string
	Authentication   struct{}
	BillingSource    struct{}
	KokuUUID         string `pg:"koku_uuid"`
	PendingDelete    bool
	PendingUpdate    bool
	OutOfOrderDelete bool
	Status           struct{}
}

func SelectSource(c context.Context, sourceID uint64) (*Source, error) {
	source := new(Source)
	if err := postgres.PGMain().ModelContext(c, source).
		Where("source_id = ?", sourceID).
		Select(); err != nil {
		return nil, err
	}
	return source, nil
}

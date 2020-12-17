package sources

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

type Source struct {
	tableName struct{} `pg:"api_sources"`

	ID               int       `pg:"source_id,pk"` // Id is automatically detected as primary key
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

func SelectSource(c context.Context, db *pg.DB, sourceID int) (*Source, error) {
	source := new(Source)
	if err := db.ModelContext(c, source).
		Where("source_id = ?", sourceID).
		Select(); err != nil {
		return nil, err
	}
	return source, nil
}

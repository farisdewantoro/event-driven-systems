package util

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Preload struct {
	AssociationTable string
	Argument         func(db *gorm.DB) *gorm.DB
}

type DbOptions struct {
	Transaction      *gorm.DB
	Preloads         []Preload
	SaveAssociations bool
	Clause           clause.Expression
}

func (o *DbOptions) Extract(ctx context.Context, db *gorm.DB) *gorm.DB {
	if o != nil && o.Transaction != nil {
		db = o.Transaction
	}

	db = db.WithContext(ctx)

	if o != nil && !o.SaveAssociations {
		db = db.Omit(clause.Associations)
	}

	for _, preload := range o.Preloads {
		if preload.Argument != nil {
			db = db.Preload(preload.AssociationTable, preload.Argument)
		} else {
			db = db.Preload(preload.AssociationTable)
		}
	}

	if o.Clause != nil {
		db = db.Clauses(o.Clause)
	}

	return db
}

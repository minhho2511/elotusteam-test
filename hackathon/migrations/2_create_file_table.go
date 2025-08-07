package migrations

import (
	"context"
	"github.com/minhho2511/elotusteam-test/internal/models"
	"github.com/uptrace/bun"
	"reflect"
	"time"
)

type CreateFileTable struct {
	Version int
}

func (m CreateFileTable) Up(db *bun.DB) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err = db.NewCreateTable().
		Model((*models.File)(nil)).
		IfNotExists().
		WithForeignKeys().
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (m CreateFileTable) Down(db *bun.DB) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err = db.NewDropTable().
		Model((*models.File)(nil)).
		IfExists().
		Cascade().
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (m CreateFileTable) GetStructName() string {
	if t := reflect.TypeOf(m); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

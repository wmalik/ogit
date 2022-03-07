package db

import (
	"context"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Database struct {
	DB *gorm.DB
}

func NewDB(dbPath string) (*Database, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Database{db}, nil
}

func (d *Database) Init() error {
	if err := d.DB.AutoMigrate(&Repository{}); err != nil {
		return err
	}

	return nil
}

func (d *Database) InsertRepository(ctx context.Context, repo *Repository) error {
	result := d.DB.WithContext(ctx).Create(repo)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (d *Database) UpsertRepositories(ctx context.Context, repos []Repository) error {
	result := d.DB.
		WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(&repos, 100)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (d *Database) SelectAllRepositories(ctx context.Context) ([]Repository, error) {
	var repos []Repository
	if result := d.DB.WithContext(ctx).Find(&repos); result.Error != nil {
		return nil, result.Error
	}

	return repos, nil
}

func (d *Database) SelectRepositories(ctx context.Context, org, filter string) ([]Repository, error) {
	var repos []Repository
	if result := d.DB.WithContext(ctx).
		Where("owner = ?", org).
		Where("name LIKE ?", "%"+filter+"%").
		Find(&repos); result.Error != nil {
		return nil, result.Error
	}

	return repos, nil
}

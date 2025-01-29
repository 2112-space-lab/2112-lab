package seeds

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var Seeds *gormigrate.Gormigrate
var SeedsList = []*gormigrate.Migration{}

func Init(db *gorm.DB) {
	Seeds = gormigrate.New(db, gormigrate.DefaultOptions, SeedsList)
}

func AddSeed(seed *gormigrate.Migration) {
	SeedsList = append(SeedsList, seed)
}

func Apply() error {
	return Seeds.Migrate()
}

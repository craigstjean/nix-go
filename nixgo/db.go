package nixgo

import (
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Name     string
	Path     string
	Packages []ProjectPackage
}

type ProjectPackage struct {
	gorm.Model
	Name      string
	ProjectID uint
}

func Start() *gorm.DB {
	homedir, _ := os.UserHomeDir()
	path := filepath.Join(homedir, ".config/nix-go/nixgo.db")
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Project{})
	db.AutoMigrate(&ProjectPackage{})

	return db
}

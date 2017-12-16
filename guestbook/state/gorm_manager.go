package state

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // using blank import cos that's how I know this to work
	"github.com/nikogura/guestbook/guestbook/config"
	"github.com/pkg/errors"
	"log"
)

// GORMStateManager  the thing that connects to the db to maintain state
type GORMStateManager struct {
	Config config.Config
	Logger *log.Logger
	db     *gorm.DB
}

// Visitor holds the information pertaining to an individual visitor
type Visitor struct {
	Name string
	IP   string
}

// NewGORMManager returns guess what?  A new GORMManager
func NewGORMManager(config config.Config, logger *log.Logger) (manager GORMStateManager, err error) {

	connectString, ok := config.Get("state.manager.connect_string")
	if !ok {
		return manager, errors.New("No db connection string found under [state.manager.connect_string] in config")
	}

	dialect := config.GetString("state.manager.dialect", "postgres")

	db, err := gorm.Open(dialect, connectString)
	if err != nil {
		return manager, err
	}

	db.LogMode(true)

	if !db.HasTable(&Visitor{}) {
		db.AutoMigrate(&Visitor{})
	}

	manager = GORMStateManager{
		Config: config,
		Logger: logger,
		db:     db,
	}

	return manager, err
}

// GetVisitor returns a visitor from the db, or an empty object if the visitor doesn't exist
func (gm *GORMStateManager) GetVisitor(ip string) (visitor Visitor, err error) {
	gm.db.Where("ip = ?", ip).First(&visitor)
	if visitor.Name != "" {
		return visitor, err
	}
	return visitor, err
}

// NewVisitor creates a new visitor in the db
func (gm *GORMStateManager) NewVisitor(visitor Visitor) (Visitor, error) {

	err := gm.db.Create(&visitor).Error

	return visitor, err
}

// RemoveVisitor removes a visitor from the db
func (gm *GORMStateManager) RemoveVisitor(visitor Visitor) (err error) {
	err = gm.db.Delete(&visitor).Error

	return err

}

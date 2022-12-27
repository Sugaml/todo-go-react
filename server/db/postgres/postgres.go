package postgres

import (
	"github.com/jinzhu/gorm"
	"github.com/sugam/golang-react-todo/models"
)

func Migration(db *gorm.DB) {
	db.AutoMigrate(
		models.ToDoList{},
	)
}

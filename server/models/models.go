package models

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type ToDoList struct {
	//ID     bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	gorm.Model

	Task   string `gorm:"not null" json:"task,omitempty"`
	Status bool   `gorm:"not null" json:"status,omitempty"`
}

func FindAll(db *gorm.DB) ([]ToDoList, error) {
	datas := []ToDoList{}
	err := db.Model(&ToDoList{}).Order("id desc").Find(&datas).Error
	if err != nil {
		return []ToDoList{}, err
	}
	return datas, nil
}

func CreateTodoList(db *gorm.DB, data ToDoList) (ToDoList, error) {
	err := db.Model(&ToDoList{}).Create(&data).Error
	fmt.Println("dat", data)
	if err != nil {
		return ToDoList{}, err
	}
	return data, nil
}

func FindByIdGateway(db *gorm.DB, pid uint) (ToDoList, error) {
	data := ToDoList{}
	err := db.Model(&ToDoList{}).Where("id = ?", pid).Take(&data).Error
	if err != nil {
		return ToDoList{}, err
	}
	return data, nil
}

func UpdateGateway(db *gorm.DB, data ToDoList) (ToDoList, error) {

	err := db.Model(&ToDoList{}).Where("id = ?", data.ID).Updates(data).Error
	if err != nil {
		return ToDoList{}, err
	}
	return data, nil
}

func DeleteTodoList(db *gorm.DB, pid uint) (int64, error) {
	result := db.Model(&ToDoList{}).Where("id = ?", pid).Take(&ToDoList{}).Delete(&ToDoList{})
	if result.Error != nil {
		if gorm.IsRecordNotFoundError(result.Error) {
			return 0, errors.New("gateway not found")
		}
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

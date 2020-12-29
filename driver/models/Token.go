package models

import (
	"errors"
	"html"
	"log"
	"os"
	"strings"
	"time"


	"github.com/jinzhu/gorm"
)


// PlaidIntegration Table that stores plaid access info needed for requests to linked bank accounts
type PlaidIntegration struct {
	ID          uint     `gorm:"primary_key;auto_increment" json:"id"`
	User		User	 `json:"user"`
	UserID      string   `json:"userid"`
	ItemID      string   `json:"itemid"`
	AccessToken string   `json:"accesstoken"`
	PaymentID   string   `json:"paymentid"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}


func (i *PlaidIntegration) SaveToken(db*gorm.DB) (*PlaidIntegration, error) {

	var err error
	err = db.Debug().Model(&PlaidIntegration).Create(&i).Error
	if err != nil {
		return &PlaidIntegration{}, err
	}
	if i.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", i.UserID).Take(&i.User).Error
		if err != nil {
			return &PlaidIntegration{}, err
		}
	}
	return i, nil

}

func (i *PlaidIntegration) FindUserIntegrations(db *gorm.DB, uid uint32) (*[]PlaidIntegration, error) {

	var err error
	integrations := []PlaidIntegration{}
	err = db.Debug().Model(&PlaidIntegration{}).Where("author_id = ?", uid).Limit(100).Order("created_at desc").Find(&integrations).Error
	if err != nil {
		return &[]PlaidIntegration{}, err
	}
	if len(integrations) > 0 {
		for i, _ := range integrations {
			err := db.Debug().Model(&User{}).Where("id = ?", integrations[i].UserID).Take(&integrations[i].User).Error
			if err != nil {
				return &[]PlaidIntegration{}, err
			}
		}
	}
	return &integrations, nil
}

func (i *PlaidIntegration) UpdateAItem(db *gorm.DB) (*PlaidIntegration, error) {

	var err error

	err = db.Debug().Model(&PlaidIntegration{}).Where("id = ?", i.ID).Updates(PlaidIntegration{ItemID: i.ItemID, AccessToken: i.AccessToken,  PaymentID: i.PaymentID, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Post{}, err
	}
	if i.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", i.UserID).Take(&i.User).Error
		if err != nil {
			return &PlaidIntegration{}, err
		}
	}
	return i, nil
}

//When a user is deleted, we also delete the post that the user had
func (i *PlaidIntegration) DeleteUserIntegrations(db *gorm.DB, uid uint32) (int64, error) {
	integrations := []PlaidIntegration{}
	db = db.Debug().Model(&PlaidIntegration{}).Where("author_id = ?", uid).Find(&integrations).Delete(&integrations)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil


func (i *PlaidIntegration) DeleteAIntegration(db *gorm.DB) (int64, error) {

	db = db.Debug().Model(&PlaidIntegration{}).Where("id = ?", i.ID).Take(&PlaidIntegration{}).Delete(&PlaidIntegration{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
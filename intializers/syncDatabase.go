package intializers

import "github.com/viru56/go-jwt/models"

func SyncDatabase () {
	DB.AutoMigrate(&models.User{})
}
package converter

import (
	"user_central/models"
)

func ConverttoEntity(boundary models.UserBoundary) models.UserEntity {
	return models.UserEntity{ID: "",
		Email:                 boundary.Email,
		Username:              boundary.Username,
		Password:              boundary.Password,
		Roles:                 boundary.Roles,
		RegistrationTimestamp: boundary.RegistrationTimestamp}
}

func ConverttoBoundary(entity models.UserEntity) models.UserBoundary {
	return models.UserBoundary{
		Email:                 entity.Email,
		Username:              entity.Username,
		Password:              entity.Password,
		Roles:                 entity.Roles,
		RegistrationTimestamp: entity.RegistrationTimestamp,
	}
}

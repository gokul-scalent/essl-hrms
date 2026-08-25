package converter

import (
	"time"

	coreAPIModel "github.com/scalent.io/scalent-hrms/apimodel/core"
	"github.com/scalent.io/scalent-hrms/entity"
)

func UserEntityToUserAPIModelResponse(e entity.User) coreAPIModel.UserResponse {
	var lastLoginAt *time.Time

	if !e.LastLoginAt.IsZero() {
		lastLoginAt = &e.LastLoginAt
	}
	list := coreAPIModel.UserResponse{

		ID:            e.ID,
		Email:         e.Email,
		Status:        e.Status,
		IsPasswordSet: e.IsPasswordSet,
		LastLoginAt:   lastLoginAt,
	}
	return list
}

//-----==-----==DO NOT ADD CODE BELOW THIS LINE------

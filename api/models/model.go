package models

import (
	"fmt"
	"io"
	"strconv"
)

type NewCategory struct {
	Category string `json:"category"`
}

type NewProject struct {
	Name         string                            `json:"name" validate:"required,gte=3,lte=50,notblank"`
	Description  string                            `json:"description" validate:"required,notblank"`
	ProfilePhoto *string                           `json:"profilePhoto,omitempty"`
	CategoryID   string                            `json:"categoryId" validate:"required,integer"`
	Rooms        *int                              `json:"rooms" validate:"required,min=1"`
	Floors       *int                              `json:"floors" validate:"required,min=1"`
	Price        *float64                          `json:"price,omitempty"`
	ScreenShot   []*NewScreenshotWithCreateproject `json:"screenShot,omitempty"`
}

type NewProjectMember struct {
	UserID    string `json:"userId" validate:"required,integer"`
	ProjectID string `json:"projectId" validate:"required,integer"`
	Role      Role   `json:"role" validate:"required"`
}

type NewScreenshot struct {
	ProjectID string `json:"projectId" validate:"required,integer"`
	ImageURL  string `json:"imageUrl" validate:"required"`
}

type NewScreenshotWithCreateproject struct {
	ImageURL string `json:"imageUrl" validate:"required"`
}

type UpdateProject struct {
	ID           string   `json:"id" validate:"required,notblank,integer"`
	Name         *string  `json:"name,omitempty"`
	Description  *string  `json:"description,omitempty"`
	ProfilePhoto *string  `json:"profilePhoto,omitempty"`
	CategoryID   *string  `json:"categoryId,omitempty" validate:"integer"`
	Rooms        *int     `json:"rooms,omitempty"`
	Floors       *int     `json:"floors,omitempty"`
	Price        *float64 `json:"price,omitempty"`
	IsHide       *bool    `json:"isHide,omitempty"`
}

type ProjectFilter struct {
	MinRooms  *int `json:"minRooms,omitempty"`
	MaxRooms  *int `json:"maxRooms,omitempty"`
	MinFloors *int `json:"minFloors,omitempty"`
	MaxFloors *int `json:"maxFloors,omitempty"`
	PriceMin  *int `json:"priceMin,omitempty"`
	PriceMax  *int `json:"priceMax,omitempty"`
}

type Role string

const (
	RoleOwner  Role = "OWNER"
	RoleAdmin  Role = "ADMIN"
	RoleEditor Role = "EDITOR"
	RoleViewer Role = "VIEWER"
)

var AllRole = []Role{
	RoleOwner,
	RoleAdmin,
	RoleEditor,
	RoleViewer,
}

func (e Role) IsValid() bool {
	switch e {
	case RoleOwner, RoleAdmin, RoleEditor, RoleViewer:
		return true
	}
	return false
}

func (e Role) String() string {
	return string(e)
}

func (e *Role) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Role(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Role", str)
	}
	return nil
}

func (e Role) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

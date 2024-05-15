package repository

import (
	"context"
	"database/sql"
	"project_management/api/constants"
	req "project_management/api/models"
	er "project_management/errors"

	"github.com/jmoiron/sqlx"
)

func AddScreenshot(ctx context.Context, input req.NewScreenshot) (string, error) {
	db := ctx.Value("db").(*sql.DB)
	userId := ctx.Value(constants.UserIDCtxKey).(string)

	role, err := getUserRole(ctx, userId, input.ProjectID)
	if err != nil {
		return "", err
	}

	if role == string(req.RoleOwner) {
		query := "INSERT INTO screenshots (project_id,image_url) VALUES ( ?, ?);"
		query = sqlx.Rebind(sqlx.DOLLAR, query)
		_, err := db.Exec(query, input.ProjectID, input.ImageURL)
		if err != nil {
			return "", er.DatabaseErrorHandling(err)
		}
	}
	return "ScreenShot Added Successfully", err
}

// func DeleteScreenshots(ctx context.Context, input []string) (string, error) {
// 	db := ctx.Value("db").(*sql.DB)
// 	userId := ctx.Value(constants.UserIDCtxKey).(string)

// 	role, err := getUserRole(ctx, userId, input.ProjectID)
// 	if err != nil {
// 		return "", err
// 	}
// }

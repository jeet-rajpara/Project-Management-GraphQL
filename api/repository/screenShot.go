package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
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

func DeleteScreenshots(ctx context.Context, ids []string, projectID string) (string, error) {
	db := ctx.Value("db").(*sql.DB)
	userId := ctx.Value(constants.UserIDCtxKey).(string)
	fmt.Println(ids)

	_, err := checkProjectExistence(db, projectID)
	if err != nil {
		log.Printf("Error in checking project existence: %v", err)
		return "", er.DatabaseErrorHandling(err)
	}

	role, err := getUserRole(ctx, userId, projectID)
	if err != nil {
		return "", err
	}

	if role == string(req.RoleOwner) {
		query := "DELETE FROM screenshots WHERE id IN (?)"
		query, arguments, err := sqlx.In(query, ids)
		if err != nil {
			log.Printf("Error in returning new argument list: %v", err)
			return "", er.InternalServerError
		}
		query = sqlx.Rebind(sqlx.DOLLAR, query)
		_, err = db.Query(query, arguments...)
		if err != nil {
			return "", er.DatabaseErrorHandling(err)
		}
	}
	return "Screenshots Deleted Successfully", nil
}

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"project_management/api/constants"

	// "project_management/api/models"
	req "project_management/api/models"
	er "project_management/errors"

	"github.com/jmoiron/sqlx"
)

func ShareProject(ctx context.Context, input req.NewProjectMember) (string, error) {
	db := ctx.Value("db").(*sql.DB)
	userID := ctx.Value(constants.UserIDCtxKey).(string)

	// get user role
	role, err := getUserRole(ctx, userID, input.ProjectID)
	if err != nil {
		return "", err
	}

	// if the user is either admin or owner then only share project
	if role == string(req.RoleAdmin) || role == string(req.RoleOwner) {
		query := "INSERT INTO project_member (project_id, user_id, role) VALUES (?, ?, ?);"
		query = sqlx.Rebind(sqlx.DOLLAR, query)
		_, err := db.Exec(query, input.ProjectID, input.UserID, input.Role)
		if err != nil {
			return "", er.DatabaseErrorHandling(err)
		}
		message := fmt.Sprintf("Project with ID %s is successfully shared with User ID %s", input.ProjectID, input.UserID)
		return message, nil
	}

	return "", er.UnauthorizedError
}

func getUserRole(ctx context.Context, userId string, projectId string) (string, error) {

	db := ctx.Value("db").(*sql.DB)
	var role string
	query := "SELECT role FROM project_member WHERE project_id= ? AND user_id = ?;"
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	row := db.QueryRow(query, projectId, userId)
	err := row.Scan(&role)
	fmt.Println(role)
	if err != nil {
		return "", er.DatabaseErrorHandling(err)
	}
	return role, nil
}

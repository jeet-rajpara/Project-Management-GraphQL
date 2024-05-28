package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	// "fmt"
	"log"
	"project_management/api/constants"
	req "project_management/api/models"
	er "project_management/errors"
	"project_management/graph/model"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// create new project
func CreateProject(ctx context.Context, newProject req.NewProject) (string, error) {
	db := ctx.Value("db").(*sql.DB)
	creatorId := ctx.Value(constants.UserIDCtxKey).(string)
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
		return "", er.InternalServerError
	}

	tx, err = manageInsertProject(tx, newProject, creatorId)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Printf("Error in rolling back transaction: %v", rbErr)
		}
		return "", er.InternalServerError
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error in committing transaction: %v", err)
		return "", er.InternalServerError
	}

	return "New Project Created Successfully.", nil
}

func manageInsertProject(tx *sql.Tx, newProject req.NewProject, creatorId string) (*sql.Tx, error) {
	var projectID int64
	query := "INSERT INTO projects (id,name,description,created_at,creator_id,category_id,profile_photo,rooms,floors,price) VALUES (DEFAULT,$1, $2, $3,$4,$5,$6,$7,$8,$9) RETURNING id;"
	row := tx.QueryRow(query, newProject.Name, newProject.Description, time.Now(), creatorId, newProject.CategoryID, newProject.ProfilePhoto, newProject.Rooms, newProject.Floors, newProject.Price)
	err := row.Scan(&projectID)
	if err != nil {
		log.Printf("Error inserting project: %v", err)
		return tx, er.InternalServerError
	}

	// store project screenshots
	if len(newProject.ScreenShot) > 0 {
		tx, err = manageProjectScreenshots(tx, projectID, newProject)
		if err != nil {
			log.Printf("Error in managing project screenshot: %v", err)
			return tx, er.InternalServerError
		}
	}

	// Add user with owner role to project_member table
	insertOwnerQuery := "INSERT INTO project_member (project_id, user_id, role) VALUES ($1, $2, $3)"
	_, err = tx.Exec(insertOwnerQuery, projectID, creatorId, req.RoleOwner)
	if err != nil {
		log.Printf("Error in inserting owner in project_member table: %v", err)
		return tx, er.InternalServerError
	}
	return tx, nil
}

func manageProjectScreenshots(tx *sql.Tx, projectID int64, newProject req.NewProject) (*sql.Tx, error) {
	query := "INSERT INTO screenshots (project_id, image_url) VALUES "

	placeholders := make([]string, len(newProject.ScreenShot))
	values := make([]interface{}, len(newProject.ScreenShot)*2)

	for i, screenshot := range newProject.ScreenShot {
		placeholders[i] = fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2)
		values[i*2] = projectID
		values[i*2+1] = screenshot.ImageURL
	}

	query += strings.Join(placeholders, ", ")
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	_, err := tx.Exec(query, values...)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Printf("Error in rolling back transaction: %v", rbErr)
			return tx, er.InternalServerError
		}
		return tx, er.DatabaseErrorHandling(err)
	}

	return tx, nil
}

func UpdateProject(ctx context.Context, updateProject req.UpdateProject) (model.Project, error) {

	db := ctx.Value("db").(*sql.DB)
	userID := ctx.Value(constants.UserIDCtxKey).(string)

	projectID, err := checkProjectExistence(db, updateProject.ID)
	if err != nil {
		fmt.Println(err)
		return model.Project{}, er.DatabaseErrorHandling(err)

	}

	role, err := getUserRole(ctx, userID, projectID)
	if err != nil {
		return model.Project{}, err
	}

	if role == string(req.RoleViewer) {
		return model.Project{}, er.UnauthorizedError
	}

	var project model.Project
	updateQuery := "UPDATE projects SET "
	var args []interface{}
	fmt.Println(len(args))
	if updateProject.Name != nil {
		updateQuery += " name = ?, "
		args = append(args, *updateProject.Name)
	}
	if updateProject.Description != nil {
		updateQuery += " description = ?, "
		args = append(args, *updateProject.Description)
	}
	if updateProject.CategoryID != nil {
		updateQuery += " category_id = ?, "
		args = append(args, *updateProject.CategoryID)
	}

	if updateProject.ProfilePhoto != nil {
		updateQuery += " profile_photo = ?, "
		args = append(args, *updateProject.ProfilePhoto)
	}
	if updateProject.Floors != nil {
		updateQuery += " floors = ?, "
		args = append(args, *updateProject.Floors)
	}
	if updateProject.Rooms != nil {
		updateQuery += " rooms = ?, "
		args = append(args, *updateProject.Rooms)
	}
	if updateProject.Price != nil {
		updateQuery += " price = ?, "
		args = append(args, *updateProject.Price)
	}
	if updateProject.IsHide != nil {
		updateQuery += " is_hide = ?, "
		args = append(args, *updateProject.IsHide)
	}

	updateQuery = strings.TrimSuffix(updateQuery, ", ")
	updateQuery += " WHERE id = ?  RETURNING id, name, description, category_id, created_at, creator_id, profile_photo, floors, rooms, price;"
	updateQuery = sqlx.Rebind(sqlx.DOLLAR, updateQuery)
	args = append(args, updateProject.ID)

	row := db.QueryRow(updateQuery, args...)
	err = row.Scan(&project.ID, &project.Name, &project.Description, &project.CategoryID, &project.CreatedAt, &project.CreatorID, &project.ProfilePhoto, &project.Floors, &project.Rooms, &project.Price)
	if err != nil {
		return model.Project{}, er.DatabaseErrorHandling(err)
	}
	return project, nil
}

func DeleteProject(ctx context.Context, projectID string) (string, error) {

	db := ctx.Value("db").(*sql.DB)
	userID := ctx.Value(constants.UserIDCtxKey).(string)
	var project model.Project

	getCreatorIdQuery := "SELECT creator_id FROM projects WHERE id = ?;"
	getCreatorIdQuery = sqlx.Rebind(sqlx.DOLLAR, getCreatorIdQuery)
	row := db.QueryRow(getCreatorIdQuery, projectID)
	err := row.Scan(&project.CreatorID)
	if err != nil {
		log.Printf("Error in fetching Creator ID: %v", err)
		return "", er.DatabaseErrorHandling(err)
	}

	if project.CreatorID == userID {
		query := "DELETE FROM projects WHERE id = ?;"
		query = sqlx.Rebind(sqlx.DOLLAR, query)
		_, err := db.Exec(query, projectID)
		if err != nil {
			log.Printf("Error in deleting project: %v", err)
			return "", er.DatabaseErrorHandling(err)
		}

		message := fmt.Sprintf("Project with ID %s deleted successfully", projectID)
		return message, nil
	}
	return "User is not authorized to delete project", nil
}

func Projects(ctx context.Context, limit *int, filter *req.ProjectFilter, sortBy *model.ProjectSort) ([]*model.Project, error) {
	db := ctx.Value("db").(*sql.DB)

	var projects []*model.Project

	var args []interface{}
	var where []string
	var joins []string
	orderBy := "ORDER BY created_at DESC"
	where = append(where, "WHERE is_hide = 'false'")

	if filter.MinRooms != nil {
		args = append(args, *filter.MinRooms)
		where = append(where, " AND rooms >= ?")
	}
	if filter.MaxRooms != nil {
		args = append(args, *filter.MaxRooms)
		where = append(where, " AND rooms <= ?")
	}
	if filter.MinFloors != nil {
		args = append(args, *filter.MinFloors)
		where = append(where, " AND floors >= ?")
	}
	if filter.MaxFloors != nil {
		args = append(args, *filter.MaxFloors)
		where = append(where, " AND floors <= ?")
	}
	if filter.PriceMin != nil {
		args = append(args, *filter.PriceMin)
		where = append(where, " AND price >= ?")
	}
	if filter.PriceMax != nil {
		args = append(args, *filter.PriceMax)
		where = append(where, " AND price <= ?")
	}

	switch sortBy.String() {
	case "NEWEST":
		orderBy = " ORDER BY created_at DESC"
	case "OLDEST":
		orderBy = " ORDER BY created_at ASC"
	case "NAME_ASC":
		orderBy = " ORDER BY LOWER(name) ASC"
	case "NAME_DESC":
		orderBy = " ORDER BY LOWER(name) DESC"
	}

	if limit != nil && *limit > 0 {
		args = append(args, *limit)
	}

	query := fmt.Sprintf("SELECT id,name,description,created_at,creator_id,category_id,profile_photo,rooms,floors,price,is_hide FROM projects %v %v %v LIMIT ?", strings.Join(joins, " "), strings.Join(where, " "), orderBy)
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Printf("Error in getting projects data: %v", err)
		return nil, er.DatabaseErrorHandling(err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			return
		}
	}()
	for rows.Next() {
		var project model.Project
		if err := rows.Scan(&project.ID, &project.Name, &project.Description, &project.CreatedAt, &project.CreatorID, &project.CategoryID, &project.ProfilePhoto, &project.Rooms, &project.Floors, &project.Price, &project.IsHide); err != nil {
			log.Printf("Error in Scanning the data: %v", err)
			return nil, err
		}
		projects = append(projects, &project)
	}
	return projects, nil
}

func Project(ctx context.Context, id string) (*model.GetProjectDetail, error) {

	db := ctx.Value("db").(*sql.DB)

	_, err := checkProjectExistence(db, id)
	if err != nil {
		log.Printf("Error in checking project existence: %v", err)
		return nil, er.DatabaseErrorHandling(err)
	}

	query := "SELECT id,name,description,created_at,creator_id,category_id,profile_photo,rooms,floors,price,is_hide FROM projects WHERE id= ?"
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	row := db.QueryRow(query, id)

	var project model.GetProjectDetail
	if err := row.Scan(&project.ID, &project.Name, &project.Description, &project.CreatedAt, &project.CreatorID, &project.CategoryID, &project.ProfilePhoto, &project.Rooms, &project.Floors, &project.Price, &project.IsHide); err != nil {
		return nil, er.DatabaseErrorHandling(err)
	}
	return &project, nil
}

func checkProjectExistence(db *sql.DB, projectID string) (string, error) {
	var tempProjectID string
	row := db.QueryRow("SELECT id FROM projects WHERE id = $1", projectID)
	err := row.Scan(&tempProjectID)
	if err != nil {
		fmt.Println(err)
		log.Printf("Error in Scanning the data: %v", err)
		return "", er.DatabaseErrorHandling(err)
	}
	return tempProjectID, nil
}

func ProjectsByUserIDs(ctx context.Context, ids []string, limit *int, filter *req.ProjectFilter, sortBy *model.ProjectSort) ([]*model.UserProjectDetail, error) {
	db := ctx.Value("db").(*sql.DB)

	var projectsByUserIDs []*model.UserProjectDetail

	var args []interface{}
	args = append(args, pq.Array(ids))
	var where []string
	var joins []string
	orderBy := "ORDER BY created_at DESC"

	if filter.MinRooms != nil {
		args = append(args, *filter.MinRooms)
		where = append(where, " AND p.rooms >= ?")
	}
	if filter.MaxRooms != nil {
		args = append(args, *filter.MaxRooms)
		where = append(where, " AND p.rooms <= ?")
	}
	if filter.MinFloors != nil {
		args = append(args, *filter.MinFloors)
		where = append(where, " AND p.floors >= ?")
	}
	if filter.MaxFloors != nil {
		args = append(args, *filter.MaxFloors)
		where = append(where, " AND p.floors <= ?")
	}
	if filter.PriceMin != nil {
		args = append(args, *filter.PriceMin)
		where = append(where, " AND p.price >= ?")
	}
	if filter.PriceMax != nil {
		args = append(args, *filter.PriceMax)
		where = append(where, " AND p.price <= ?")
	}

	switch sortBy.String() {
	case "NEWEST":
		orderBy = " ORDER BY p.created_at DESC"
	case "OLDEST":
		orderBy = " ORDER BY p.created_at ASC"
	case "NAME_ASC":
		orderBy = " ORDER BY p.name ASC"
	case "NAME_DESC":
		orderBy = " ORDER BY p.name DESC"
	}

	if limit != nil && *limit > 0 {
		args = append(args, *limit)
	}

	query := fmt.Sprintf("SELECT p.id, p.name, p.description, p.profile_photo, p.category_id, p.rooms, p.floors, p.price, p.created_at, p.creator_id, u.name AS creator_name, u.email AS creator_email FROM projects p LEFT JOIN users u ON p.creator_id = u.id WHERE p.creator_id = ANY ( ? ) %v %v %v LIMIT ?", strings.Join(joins, " "), strings.Join(where, " "), orderBy)
	query = sqlx.Rebind(sqlx.DOLLAR, query)
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Printf("Error in getting data: %v", err)
		return nil, er.DatabaseErrorHandling(err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Error in closing rows: %v", err)
		}
	}()

	for rows.Next() {
		var project model.UserProjectDetail
		var creator model.User
		err := rows.Scan(&project.ID, &project.Name, &project.Description, &project.ProfilePhoto, &project.CategoryID, &project.Rooms, &project.Floors, &project.Price, &project.CreatedAt, &project.CreatorID, &creator.Name, &creator.Email)
		if err != nil {
			log.Printf("Error in Scanning the data: %v", err)
			return nil, er.DatabaseErrorHandling(err)
		}

		project.Creator = &creator
		projectsByUserIDs = append(projectsByUserIDs, &project)
	}

	return projectsByUserIDs, nil
}

package dataloaders

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"project_management/graph/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type ctxKeyType struct{ name string }

var CtxKey = ctxKeyType{"dataloaderctx"}

type Loaders struct {
	Screenshots              *ScreenshotsSliceLoader
	UserByID                 *UserLoader
	ProjectMemberByProjectID *ProjectMemberLoader
}

func LoaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		db := ctx.Value("db").(*sql.DB)
		loader := Loaders{}

		wait := 350 * time.Microsecond
		//maxBatch := 200

		loader.Screenshots = &ScreenshotsSliceLoader{
			wait:     wait,
			maxBatch: 100,
			fetch: func(ids []string) ([][]*model.Screenshot, []error) {
				var sqlQuery string
				if len(ids) == 1 {
					sqlQuery = "SELECT id, project_id, image_url FROM screenshots WHERE project_id = ?"
				} else {
					sqlQuery = "SELECT id, project_id, image_url FROM screenshots WHERE project_id IN (?)"
				}
				sqlQuery, arguments, err := sqlx.In(sqlQuery, ids)
				if err != nil {
					log.Printf("Error in returning new query and args : %v", err)
				}

				sqlQuery = sqlx.Rebind(sqlx.DOLLAR, sqlQuery)
				rows, err := db.Query(sqlQuery, arguments...)
				if err != nil {
					log.Printf("Error in getting data: %v", err)
				}
				screenshotsByProjectID := map[string][]*model.Screenshot{}
				for rows.Next() {
					var screenshot model.Screenshot
					if err := rows.Scan(&screenshot.ID, &screenshot.ProjectID, &screenshot.ImageURL); err != nil {
						log.Printf("Error in Scanning rows: %v", err)
					}

					screenshotPairs := screenshotsByProjectID[screenshot.ProjectID]

					if screenshotPairs == nil {
						screenshotArr := []*model.Screenshot{}
						screenshotArr = append(screenshotArr, &screenshot)
						screenshotsByProjectID[screenshot.ProjectID] = screenshotArr
					} else {
						screenshotsByProjectID[screenshot.ProjectID] = append(screenshotsByProjectID[screenshot.ProjectID], &screenshot)
					}
				}

				if err := rows.Close(); err != nil {
					log.Printf("Error in closing rows: %v", err)
				}
				screenshots := make([][]*model.Screenshot, len(ids))
				for i, id := range ids {
					screenshots[i] = screenshotsByProjectID[id]
					i++
				}
				return screenshots, nil
			},
		}

		loader.UserByID = &UserLoader{
			wait:     wait,
			maxBatch: 100,
			fetch: func(ids []string) ([]*model.User, []error) {
				var sqlQuery string
				if len(ids) == 1 {
					sqlQuery = "SELECT id, name,email  FROM users WHERE id = ?"
				} else {
					sqlQuery = "SELECT id, name,email FROM users WHERE id IN (?)"
				}
				sqlQuery, arguments, err := sqlx.In(sqlQuery, ids)
				if err != nil {
					log.Printf("Error in returning new query and args : %v", err)
				}

				sqlQuery = sqlx.Rebind(sqlx.DOLLAR, sqlQuery)
				rows, err := db.Query(sqlQuery, arguments...)
				if err != nil {
					log.Printf("Error in getting data: %v", err)
				}
				var creators []*model.User
				defer rows.Close()

				for rows.Next() {
					var creator model.User
					if err := rows.Scan(&creator.ID, &creator.Name, &creator.Email); err != nil {
						log.Printf("Error in Scanning rows: %v", err)
						return nil, []error{err}
					}
					creators = append(creators, &creator)
				}
				return creators, nil
			},
		}

		loader.ProjectMemberByProjectID = &ProjectMemberLoader{
			wait:     wait,
			maxBatch: 100,
			fetch: func(ids []string) ([][]*model.ProjectMember, []error) {

				var sqlQuery string
				if len(ids) == 1 {
					sqlQuery = "SELECT id, project_id, user_id,role FROM project_member WHERE project_id = ?"
				} else {
					sqlQuery = "SELECT id, project_id, user_id,role FROM project_member WHERE project_id IN (?)"
				}
				sqlQuery, arguments, err := sqlx.In(sqlQuery, ids)
				if err != nil {
					log.Printf("Error in returning new query and args : %v", err)

				}

				sqlQuery = sqlx.Rebind(sqlx.DOLLAR, sqlQuery)
				rows, err := db.Query(sqlQuery, arguments...)
				if err != nil {
					log.Printf("Error in getting data: %v", err)
				}
				projectMemberByProjectID := map[string][]*model.ProjectMember{}
				for rows.Next() {
					var projectMember model.ProjectMember

					if err := rows.Scan(&projectMember.ID, &projectMember.ProjectID, &projectMember.UserID, &projectMember.Role); err != nil {
						log.Printf("Error in Scanning rows: %v", err)
					}

					projectMemberPairs := projectMemberByProjectID[projectMember.ProjectID]

					if projectMemberPairs == nil {
						projectMemberArr := []*model.ProjectMember{}
						projectMemberArr = append(projectMemberArr, &projectMember)
						projectMemberByProjectID[projectMember.ProjectID] = projectMemberArr
					} else {
						projectMemberByProjectID[projectMember.ProjectID] = append(projectMemberByProjectID[projectMember.ProjectID], &projectMember)
					}
				}

				if err := rows.Close(); err != nil {
					log.Fatal(err)
				}
				projectMembers := make([][]*model.ProjectMember, len(ids))
				for i, id := range ids {
					projectMembers[i] = projectMemberByProjectID[id]
					i++
				}
				return projectMembers, nil
			},
		}
		ctx = context.WithValue(ctx, CtxKey, loader)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func CtxLoaders(ctx context.Context) Loaders {
	return ctx.Value(CtxKey).(Loaders)
}

-- migrate:up

CREATE TABLE IF NOT EXISTS users (
	id INT NOT NULL DEFAULT unique_rowid(),
	"name" STRING NOT NULL,
	email STRING NOT NULL,
	password STRING NULL,
	CONSTRAINT "primary" PRIMARY KEY (id ASC)
);

CREATE TABLE IF NOT EXISTS categories (
	id INT NOT NULL DEFAULT unique_rowid(),
	"name" STRING NOT NULL,
	CONSTRAINT pk_category_id PRIMARY KEY (id ASC)
);

CREATE TABLE IF NOT EXISTS projects (
	id INT NOT NULL DEFAULT unique_rowid(),
    "name" STRING(100) NOT NULL,
	description STRING NOT NULL,
	created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	updated_at TIMESTAMP WITHOUT TIME ZONE NULL,
	creator_id INT NOT NULL,
	category_id INT NOT NULL,
    profile_photo STRING NULL,
    rooms INT NOT NULL,
    floors INT NOT NULL,
    price FLOAT NULL,
    is_hide BOOL NOT NULL DEFAULT false,
	CONSTRAINT pk_projects PRIMARY KEY (id ASC),
	CONSTRAINT fk_projects_category_id FOREIGN KEY (category_id) REFERENCES categories (id),
    CONSTRAINT fk_projects_owner_id FOREIGN KEY (creator_id) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS screenshots(
    id INT NOT NULL DEFAULT unique_rowid(),
    project_id INT NOT NULL,
    image_url STRING NOT NULL,
    CONSTRAINT pk_screenshots PRIMARY KEY (id ASC),
    CONSTRAINT fk_project_id FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS project_member(
    id INT NOT NULL DEFAULT unique_rowid(),
    project_id INT NOT NULL,
    user_id INT NOT NULL,
    "role" STRING NOT NULL,
    CONSTRAINT pk_project_member_id PRIMARY KEY (id ASC),
    CONSTRAINT fk_project_id FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);

-- migrate:down

DROP TABLE IF EXISTS screenshots;
DROP TABLE IF EXISTS project_member;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;

 
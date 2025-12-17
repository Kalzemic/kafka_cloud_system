package storage

import (
	"database/sql"
	"user_central/models"

	"github.com/lib/pq"
)

type UserRepo interface {
	Init() error
	CreateUser(user *models.UserEntity) error
	UpdateUser(user *models.UserEntity) error
	GetbyEmailDomain(domain string, page int, size int) ([]models.UserEntity, error)
	FindUser(email string, password string) (*models.UserEntity, error)
	GetAllUsers(page int, size int) ([]models.UserEntity, error)
	GetUsersbyRoles(role string, page int, limit int) ([]models.UserEntity, error)
	GetUsersbyRegistrationToday(page int, size int) ([]models.UserEntity, error)
	DeleteUsers() error
}

type PostgresUserRepo struct {
	Database *sql.DB
}

func (repo *PostgresUserRepo) Init() error {
	query := `
    CREATE EXTENSION IF NOT EXISTS "pgcrypto";

    CREATE TABLE IF NOT EXISTS users (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        username TEXT UNIQUE NOT NULL,
        email TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL,
        roles TEXT[] NOT NULL DEFAULT ARRAY['user'],
        registrationTimestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
`
	_, err := repo.Database.Exec(query)
	return err
}

func (repo *PostgresUserRepo) CreateUser(user *models.UserEntity) error {
	query := ` INSERT INTO users (username, email, password, roles)
	VALUES ($1, $2, $3, $4)
	RETURNING id, registrationTimestamp
	`

	return repo.Database.QueryRow(query, user.Username, user.Email, user.Password, pq.Array(user.Roles)).Scan(&user.ID, &user.RegistrationTimestamp)
}

func (repo *PostgresUserRepo) FindUser(email string, password string) (*models.UserEntity, error) {
	query := ` SELECT id, username, email, password, roles, registrationTimestamp
	FROM users
	WHERE email = $1 AND password = $2
	`
	row := repo.Database.QueryRow(query, email, password)
	user := &models.UserEntity{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, pq.Array(&user.Roles), &user.RegistrationTimestamp)

	return user, err

}

func (repo *PostgresUserRepo) UpdateUser(user *models.UserEntity) error {
	query := `UPDATE users
	SET username=$2, password=$3, 
	WHERE email=$1 
	`

	_, err := repo.Database.Exec(query, user.Email, user.Username, user.Password)
	return err
}

func (repo *PostgresUserRepo) GetAllUsers(page int, size int) ([]models.UserEntity, error) {

	offset := (page - 1) * size
	query := `SELECT id, username, email, password, roles, registrationTimestamp 
	FROM users 
	ORDER BY id
	LIMIT $1 OFFSET $2
	`
	rows, err := repo.Database.Query(query, size, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.UserEntity
	for rows.Next() {
		var user models.UserEntity
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, pq.Array(&user.Roles), &user.RegistrationTimestamp); err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	return users, nil
}

func (repo *PostgresUserRepo) GetUsersbyRoles(role string, page int, size int) ([]models.UserEntity, error) {

	query := `SELECT id, username, email, password, roles, registrationTimestamp 
	FROM users 
	WHERE roles[1] = $3 OR roles[2] = $3
	ORDER BY id
	LIMIT $1 OFFSET $2
	`
	offset := (page - 1) * size

	rows, err := repo.Database.Query(query, size, offset, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.UserEntity
	for rows.Next() {
		var user models.UserEntity
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, pq.Array(&user.Roles), &user.RegistrationTimestamp); err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	return users, nil
}

func (repo *PostgresUserRepo) GetbyEmailDomain(domain string, page int, size int) ([]models.UserEntity, error) {

	offset := (page - 1) * size
	query := ` SELECT id, username, email, password, roles, registrationTimestamp
	FROM users
	WHERE email ~ $1 
	ORDER BY id
	LIMIT $2 OFFSET $3;
	`
	rows, err := repo.Database.Query(query, domain, size, offset)
	if err != nil {
		return nil, err
	}

	var users []models.UserEntity
	for rows.Next() {
		var user models.UserEntity
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, pq.Array(&user.Roles), &user.RegistrationTimestamp); err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	return users, nil
}

func (repo *PostgresUserRepo) GetUsersbyRegistrationToday(page int, size int) ([]models.UserEntity, error) {

	query := ` SELECT id, email, username, password, roles, registrationTimestamp
	FROM users
	WHERE registrationTimestamp >= NOW() - INTERVAL '24 hours'
	ORDER BY registrationTimestamp DESC
	LIMIT $1 OFFSET $2;
	`
	offset := (page - 1) * size
	rows, err := repo.Database.Query(query, size, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.UserEntity
	for rows.Next() {
		var user models.UserEntity
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, pq.Array(&user.Roles), &user.RegistrationTimestamp); err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	return users, nil
}

func (repo *PostgresUserRepo) DeleteUsers() error {
	_, err := repo.Database.Exec("DELETE FROM users")
	return err
}

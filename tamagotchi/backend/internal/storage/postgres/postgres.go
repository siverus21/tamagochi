package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // init postgres driver
)

type Storage struct {
	db *sql.DB
}

// Типы данных
type Entitlement struct {
	ID        int    `json:"id"`
	Privilege string `json:"privilege"`
}

type User struct {
	ID               int          `json:"id"`
	UserName         string       `json:"userName"`
	Email            string       `json:"email"`
	Password         string       `json:"password"`
	RegistrationDate string       `json:"registrationDate"`
	LastLoginDate    *string      `json:"lastLoginDate,omitempty"`
	Entitlement      *Entitlement `json:"entitlement,omitempty"`
}

func New(conSettings string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("postgres", conSettings)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Проверяем подключение к базе
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) GetAllUsers() ([]User, error) {
	const query = `
		SELECT 
			u.id, u.userName, u.email, u.password, u.registrationDate, u.lastLoginDate,
			e.id AS entitlementId, e.privilege
		FROM users u
		LEFT JOIN entitlement e ON u.entitlementId = e.id
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("GetAllUsers: query error: %w", err)
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		var entID sql.NullInt64
		var privilege sql.NullString
		var lastLogin sql.NullString

		err := rows.Scan(
			&user.ID,
			&user.UserName,
			&user.Email,
			&user.Password,
			&user.RegistrationDate,
			&lastLogin,
			&entID,
			&privilege,
		)
		if err != nil {
			return nil, fmt.Errorf("GetAllUsers: row scan error: %w", err)
		}

		if lastLogin.Valid {
			user.LastLoginDate = &lastLogin.String
		}

		if entID.Valid && privilege.Valid {
			user.Entitlement = &Entitlement{
				ID:        int(entID.Int64),
				Privilege: privilege.String,
			}
		}

		users = append(users, user)
	}

	return users, nil
}

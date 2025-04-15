package postgres

import (
	"database/sql"
	"fmt"
	"log"

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

// Метод для аутентификации пользователя.
// Если пользователь не найден, возвращает (false, nil).
// В случае технической ошибки возвращает (false, err).
// Прямое сравнение паролей используется только для примера.
func (s *Storage) AuthenticateUser(login, password string) (bool, error) {
	var storedPassword string
	query := "SELECT password FROM users WHERE username=$1"
	err := s.db.QueryRow(query, login).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Пользователь не найден:", login)
			return false, nil
		}
		log.Println("Ошибка при получении данных пользователя:", err)
		return false, err
	}
	log.Println("Пользователь найден!", login)
	return password == storedPassword, nil
}

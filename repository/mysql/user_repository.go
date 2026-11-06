package mysql

import (
	"beverages-cli/domain"
	"database/sql"
	"errors"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByEmailAndPassword(email, password string) (*domain.User, error) {
	query := "SELECT id, name, email, password, role, balance, created_at FROM user WHERE email = ? AND password = ? AND role = 'admin' LIMIT 1"
	var u domain.User
	err := r.db.QueryRow(query, email, password).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.Balance, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) GetByID(id int) (*domain.User, error) {
	query := "SELECT id, name, email, password, role, balance, created_at FROM user WHERE id = ?"
	var u domain.User
	err := r.db.QueryRow(query, id).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.Balance, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) CheckEmailExists(email string, excludeID int) (bool, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM user WHERE email = ? AND id != ?", email, excludeID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userRepository) Update(id int, name, email, password string) error {
	query := "UPDATE user SET name = ?, email = ?, password = ? WHERE id = ? AND role = 'admin'"
	res, err := r.db.Exec(query, name, email, password, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("tidak ada baris yang diperbarui")
	}
	return nil
}

func (r *userRepository) Delete(id int) error {
	query := "DELETE FROM user WHERE id = ? AND role = 'admin'"
	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("akun tidak ditemukan")
	}
	return nil
}

func (r *userRepository) Create(u *domain.User) error {
	query := "INSERT INTO user (name, email, password, role, balance) VALUES (?, ?, ?, 'admin', 0.00)"
	res, err := r.db.Exec(query, u.Name, u.Email, u.Password)
	if err != nil {
		return err
	}
	lastID, _ := res.LastInsertId()
	u.ID = int(lastID)
	return nil
}

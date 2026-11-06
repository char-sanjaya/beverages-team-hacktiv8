package domain

import "time"

// struct user
type User struct {
	ID        int
	Name      string
	Email     string
	Password  string
	Role      string
	Balance   float64
	CreatedAt time.Time
}

// interface/Kontrak fungsi database untuk user/admin
type UserRepository interface {
	GetByEmailAndPassword(email, password string) (*User, error)
	GetByID(id int) (*User, error)
	CheckEmailExists(email string, excludeID int) (bool, error)
	Update(id int, name, email, password string) error
	Delete(id int) error
	Create(user *User) error
}



// interface/Kontrak logika bisnis akun admin
type UserUsecase interface{
	LoginAdmin(email, password string) (*User, error)
	GetProfile(adminId int) (*User, error)
	UpdateProfile(adminId int, name, email, password string) error
	DeleteAccount(adminId int) error
	RegisterAdmin(name, email, password string) error
}

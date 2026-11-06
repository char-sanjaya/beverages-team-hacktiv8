package usecase

import (
	"errors"
	"fmt"
	"strings"

	"beverages-cli/domain" // sesuaikan nama modul di go.mod Anda
)

type userUsecase struct {
	userRepo domain.UserRepository
}

// Constructor: Menerima UserRepository lewat Dependency Injection
func NewUserUsecase(repo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{userRepo: repo}
}

// 1. Logika Login Admin
func (u *userUsecase) LoginAdmin(email, password string) (*domain.User, error) {
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)

	if email == "" || password == "" {
		return nil, errors.New("email dan password tidak boleh kosong")
	}

	user, err := u.userRepo.GetByEmailAndPassword(email, password)
	if err != nil {
		// CETAK ERROR ASLI DARI DATABASE:
		return nil, fmt.Errorf("error database asli: %v", err)
	}

	// ATURAN BISNIS KUNCI: Wajib cek role admin!
	if user.Role != "admin" {
		return nil, errors.New("akses ditolak! Akun Anda bukan akun admin")
	}

	return user, nil
}

// 2. Ambil Profil Admin
func (u *userUsecase) GetProfile(adminID int) (*domain.User, error) {
	return u.userRepo.GetByID(adminID)
}

// 3. Update Profil Admin Sendiri
func (u *userUsecase) UpdateProfile(adminID int, name, email, password string) error {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)

	if name == "" || email == "" || password == "" {
		return errors.New("nama, email, dan password wajib diisi")
	}

	// Cek apakah email sudah dipakai akun lain
	exists, err := u.userRepo.CheckEmailExists(email, adminID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email sudah digunakan oleh akun lain, gunakan email berbeda")
	}

	return u.userRepo.Update(adminID, name, email, password)
}

// 4. Hapus Akun Admin Sendiri
func (u *userUsecase) DeleteAccount(adminID int) error {
	return u.userRepo.Delete(adminID)
}

// 5. Registrasi Admin Baru
func (u *userUsecase) RegisterAdmin(name, email, password string) error {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)

	if name == "" || email == "" || password == "" {
		return errors.New("semua kolom pendaftaran wajib diisi")
	}

	exists, err := u.userRepo.CheckEmailExists(email, 0)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email sudah terdaftar di sistem")
	}

	newAdmin := &domain.User{
		Name:     name,
		Email:    email,
		Password: password,
		Role:     "admin", // Otomatis diset sebagai admin
	}
	return u.userRepo.Create(newAdmin)
}

package handler

import (
	"beverages-cli/entity"
	"database/sql"
	"fmt"
)

func AdminLandingPage(db *sql.DB, admin *entity.User) {
	fmt.Printf("Halo Admin %s, halaman admin belum diimplementasikan.\n", admin.Name)
}

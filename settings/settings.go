package settings

import (
	"os"
	"strings"
)

// Admins struct for export
type Admins struct {
	ProductionAdmins []string
	StageAdmins      []string
}

var (
	// AdminLists for exporting
	AdminLists = &Admins{parseStringList(os.Getenv("PRODUCTION_ADMIN_LIST")), []string{}}
)

func parseStringList(str string) []string {
	return strings.Split(strings.ReplaceAll(str, " ", ""), ",")
}

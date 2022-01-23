package xormcmd

var (
	genXorm                   = 0
	genJson                   = 0
	ignoreColumnsJSON         []string
	created, updated, deleted = []string{"created_at"}, []string{"updated_at"}, []string{"deleted_at"}
)

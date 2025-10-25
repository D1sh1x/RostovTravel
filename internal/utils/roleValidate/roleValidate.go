package rolevalidate

var allowedRoles = map[string]struct{}{
	"user":  {},
	"admin": {},
}

func IsValidRole(r string) bool {
	if r == "" {
		return false
	}
	_, ok := allowedRoles[r]
	return ok
}

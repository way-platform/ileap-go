package demo

// User represents a user.
type User struct {
	// Username is the user's username.
	Username string
	// Password is the user's password.
	Password string
}

// Users returns the list of demo users.
func Users() []User {
	return []User{
		{
			Username: "hello",
			Password: "pathfinder",
		},

		{
			Username: "transport_service_user",
			Password: "ileap",
		},

		{
			Username: "transport_service_organizer",
			Password: "ileap",
		},
	}
}

package domain

type Account struct {
	ID           int64
	Organization Organization `json:"organization"`
	Origin       Origin
	Contact      Contact
	Profiles     []Profile
}

type Profile struct {
	Role        string
	Permissions []string
}

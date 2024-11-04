package contacts

type Contact struct {
	ID           int      `json:"id"`
	FirstName    string   `json:"first_name"`
	LastName     string   `json:"last_name"`
	PhoneNumbers []string `json:"phone_numbers"`
}

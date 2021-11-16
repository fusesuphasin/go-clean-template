package domain

type PhoneBook struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	Company  string `json:"company"`
	Position string `json:"position"`
}

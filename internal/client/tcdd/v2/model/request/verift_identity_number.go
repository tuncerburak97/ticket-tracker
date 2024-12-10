package request

type VerifyIdentityNumberRequest struct {
	IdentityNumber string `json:"tcKimlikNo"`
	Name           string `json:"ad"`
	LastName       string `json:"soyad"`
	BirthDate      string `json:"dogumYili"`
}

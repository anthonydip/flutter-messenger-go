package dtos

type User struct {
	Id       string `firestore:"id,omitempty" json:"id,omitempty"`
	Email    string `firestore:"email,omitempty" json:"email,omitempty"`
	Provider string `firestore:"provider,omitempty" json:"provider,omitempty"`
	Password string `firestore:"password,omitempty" json:"password,omitempty"`
}

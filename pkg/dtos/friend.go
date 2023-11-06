package dtos

import (
	"fmt"
)

type Friend struct {
	Id    string `firestore:"id,omitempty" json:"id,omitempty"`
	Email string `firestore:"email,omitempty" json:"email,omitempty"`
}

func (friend Friend) String() string {
	return fmt.Sprintf("Friend{Id: %s, Email: %s}", friend.Id, friend.Email)
}

package storefront

import (
	"github.com/anthonydip/flutter-messenger-go/pkg/dtos"
)

func (bkr Broker) GetUser(id string) (dtos.User, error) {

	return dtos.User{
		Name: "Test User",
	}, nil
}

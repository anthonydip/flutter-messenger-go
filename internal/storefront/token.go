package storefront

import (
	"fmt"
)

func (bkr Broker) AddAccessToken(token string) error {
	fmt.Printf("token: %s\n", token)

	return nil
}

func (bkr Broker) InvalidateToken(token string) error {
	fmt.Printf("token: %s\n", token)

	return nil
}

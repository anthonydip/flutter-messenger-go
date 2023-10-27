package mock

import (
	"github.com/anthonydip/flutter-messenger-go/pkg/dtos"
)

type Result func(c *mockConfig)

type mockConfig struct {
	getUser           error
	getUserByEmail    error
	signIn            error
	postUser          error
	addAccessToken    error
	deleteAccessToken error
}

// Mock for mocking Storefront service
type Mock struct {
	cfg mockConfig
}

// Function to create a new, mock storefront service
func New(opts ...Result) Mock {
	r := Mock{}

	for _, o := range opts {
		if o != nil {
			o(&r.cfg)
		}
	}

	return r
}

// GetUser mocks Storefront GetUser() call
func (m Mock) GetUser(string) (dtos.User, error) {
	if m.cfg.getUser != nil {
		return dtos.User{}, m.cfg.getUser
	}

	return dtos.User{
		Id:       "8ae84a23-fa49-45eb-8000-bdc9b9fe074a",
		Email:    "mock@storefront-mock.com",
		Provider: "Flutter",
		Password: "*****",
	}, nil
}

// GetUserResult sets the result of the mock GetUser()
func GetUserResult(e error) Result {
	return func(c *mockConfig) {
		c.getUser = e
	}
}

// GetUserByEmail mocks Storefront GetUserByEmail() call
func (m Mock) GetUserByEmail(string) (dtos.User, error) {
	if m.cfg.getUserByEmail != nil {
		return dtos.User{}, m.cfg.getUserByEmail
	}

	return dtos.User{
		Id:       "8ae84a23-fa49-45eb-8000-bdc9b9fe074a",
		Email:    "mock@storefront-mock.com",
		Provider: "Flutter",
		Password: "*****",
	}, nil
}

// GetUserByEmailResult sets the result of the mock GetUserGetUserByEmail()
func GetUserByEmailResult(e error) Result {
	return func(c *mockConfig) {
		c.getUserByEmail = e
	}
}

// SignIn mocks Storefront SignIn() call
func (m Mock) SignIn(dtos.User) error {
	if m.cfg.signIn != nil {
		return m.cfg.signIn
	}

	return nil
}

// SignInResult sets the result of the mock SignIn()
func SignInResult(e error) Result {
	return func(c *mockConfig) {
		c.signIn = e
	}
}

// PostUser mocks Storefront PostUser() call
func (m Mock) PostUser(dtos.User) (dtos.User, error) {
	if m.cfg.postUser != nil {
		return dtos.User{}, m.cfg.postUser
	}

	return dtos.User{
		Id:       "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTg0NDc4NTEsImlhdCI6MTY5ODM2MTQ1MSwiVG9rZW5UeXBlIjoidXNlciIsIkVtYWlsIjoibW9ja0BzdG9yZWZyb250LW1vY2suY29tIiwiUHJvdmlkZXIiOiJGbHV0dGVyIn0.xly0f3j9-BEqdfFNToY1IwS5-1H34HbwR1qklhaLjDiWeBIhmpshOqeiq-yKBrQM7RExLVgtMqZYW6oxeagLD5dOANVKttSH207BqFo1V1T1US4JMOjMPpzYDUleqWbYdB-muW7tggH2wxR-t9aepyBYcoWK1hKiNEaAuedhfz4",
		Email:    "mock@storefront-mock.com",
		Provider: "Flutter",
		Password: "*****",
	}, nil
}

// PostUserResult sets the result of the mock PostUser()
func PostUserResult(e error) Result {
	return func(c *mockConfig) {
		c.postUser = e
	}
}

// AddAccessToken mocks Storefront AddAccessToken() call
func (m Mock) AddAccessToken(string, dtos.User) error {
	if m.cfg.addAccessToken != nil {
		return m.cfg.addAccessToken
	}

	return nil
}

// AddAccessTokenResult sets the result of the mock AddAccessToken()
func AddAccessTokenResult(e error) Result {
	return func(c *mockConfig) {
		c.addAccessToken = e
	}
}

// DeleteAccessToken mocks Storefront DeleteAccessToken() call
func (m Mock) DeleteAccessToken(string) error {
	if m.cfg.deleteAccessToken != nil {
		return m.cfg.deleteAccessToken
	}

	return nil
}

// DeleteAccessTokenResult sets the result of the mock DeleteAccessToken()
func DeleteAccessTokenResult(e error) Result {
	return func(c *mockConfig) {
		c.deleteAccessToken = e
	}
}

// TODO
func (m Mock) AccessTokenExists(string) error {
	return nil
}

// TODO
func (m Mock) PostFriend(string, dtos.User) error {
	return nil
}

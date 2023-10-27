package mock

type Result func(c *mockConfig)

type mockConfig struct {
}

// Mock the Authorization agent
type Mock struct {
	cfg mockConfig
}

// Function to create a new Mock Authorization agent
func New(opts ...Result) *Mock {
	r := &Mock{}

	for _, o := range opts {
		if o != nil {
			o(&r.cfg)
		}
	}

	return r
}

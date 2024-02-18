package ticketing

import (
	"net/http"
)

func NewHTTPClient(jwtToken string) *http.Client {
	// Create a new HTTP client
	client := &http.Client{}

	// Create a new request modifier function to add the JWT token to request headers
	modifyRequest := func(req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+jwtToken) // Set the JWT token in Authorization header
		return nil
	}

	// Wrap the client Transport with a RoundTripper that applies the request modifier
	client.Transport = &roundTripperFunc{fn: modifyRequest}

	return client
}

type roundTripperFunc struct {
	fn func(*http.Request) error
}

func (rt *roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	err := rt.fn(req)
	if err != nil {
		return nil, err
	}
	return http.DefaultTransport.RoundTrip(req)
}

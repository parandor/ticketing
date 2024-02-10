package ticketing_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	connect "connectrpc.com/connect"

	server "github.com/parandor/ticketing"
	ticketingv1 "github.com/parandor/ticketing/internal/gen/proto/train_ticketing/v1/train_ticketingv1connect"

	v1 "github.com/parandor/ticketing/internal/gen/proto/train_ticketing/v1"
)

func TestPurchaseTicket(t *testing.T) {
	_, httpHandler := server.NewMyTicketingServiceHandler()

	// Create a new HTTP server with the handler
	server := httptest.NewServer(httpHandler)
	defer server.Close()

	// Create a new TicketingSystemClient
	jwtToken := "auth_token"
	client := ticketingv1.NewTrainTicketingServiceClient(newHTTPClient(jwtToken), server.URL)

	// Test scenario 1: Purchase ticket successfully
	testPurchaseTicketSuccess(t, client)

	// Add additional assertions if needed
}

func newHTTPClient(jwtToken string) *http.Client {
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

func testPurchaseTicketSuccess(t *testing.T, client ticketingv1.TrainTicketingServiceClient) {

	response, err := client.PurchaseTicket(context.Background(), &connect.Request[v1.PurchaseTicketRequest]{
		Msg: &v1.PurchaseTicketRequest{
			Ticket: &v1.Ticket{
				From:      "City A",
				To:        "City B",
				User:      &v1.User{FirstName: "John", LastName: "Doe", Email: "john@example.com"},
				PricePaid: 25.0,
			},
		},
	})

	if err != nil {
		t.Fatalf("PurchaseTicket failed: %v", err)
	}

	// Check if the response contains the receipt
	if response == nil || response.Msg.Receipt == nil {
		t.Fatalf("PurchaseTicket response is nil or missing receipt")
	}

}

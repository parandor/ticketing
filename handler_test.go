package ticketing_test

import (
	"context"
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
	client := ticketingv1.NewTrainTicketingServiceClient(server.Client(), server.URL)

	// Test scenario 1: Purchase ticket successfully
	testPurchaseTicketSuccess(t, client)

	// Add additional assertions if needed
}

func testPurchaseTicketSuccess(t *testing.T, client ticketingv1.TrainTicketingServiceClient) {
	//client.PurchaseTicket(context.Context, *connect.Request[v1.PurchaseTicketRequest]) (*connect.Response[v1.PurchaseTicketResponse], error)

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

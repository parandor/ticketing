package ticketing_test

import (
	"context"
	"net/http/httptest"
	"testing"

	connect "connectrpc.com/connect"
	ticketing "github.com/parandor/ticketing"
	v1 "github.com/parandor/ticketing/internal/gen/proto/library/v1"
	library_v1 "github.com/parandor/ticketing/internal/gen/proto/library/v1/libraryv1connect"
	"github.com/stretchr/testify/assert"
)

func TestBorrowBook(t *testing.T) {
	_, httpHandler := ticketing.NewLibraryService()

	// Create a new HTTP server with the handler
	server := httptest.NewServer(httpHandler)
	defer server.Close()

	// Create a new TicketingSystemClient
	jwtToken := "auth_token"
	client := library_v1.NewLibraryServiceClient(ticketing.NewHTTPClient(jwtToken), server.URL)

	// Test scenario 1: Purchase ticket successfully
	testBorrowBook(t, client)


	// Add additional assertions if needed
}

func TestReturnBook(t *testing.T) {
	_, httpHandler := ticketing.NewLibraryService()

	// Create a new HTTP server with the handler
	server := httptest.NewServer(httpHandler)
	defer server.Close()

	// Create a new TicketingSystemClient
	jwtToken := "auth_token"
	client := library_v1.NewLibraryServiceClient(ticketing.NewHTTPClient(jwtToken), server.URL)

	testBorrowBook(t, client)

	// Test scenario 1: Purchase ticket successfully
	testReturnBook(t, client)


	// Add additional assertions if needed
}

func testReturnBook(t *testing.T, client library_v1.LibraryServiceClient) {
    // Mock the request and response
    request := &connect.Request[v1.ReturnBookRequest]{
        Msg: &v1.ReturnBookRequest{
            BookId: "123", // Provide the ID of the book to be returned
        },
    }
    expectedResponse := &connect.Response[v1.ReturnBookResponse]{
        Msg: &v1.ReturnBookResponse{
            Success: true,
            Message: "Book returned successfully",
        },
    }

    // Call the method being tested
    response, err := client.ReturnBook(context.Background(), request)

    // Assert that no error occurred
    assert.NoError(t, err)

    // Assert the response matches the expected response
    assert.Equal(t, expectedResponse.Msg.Message, response.Msg.Message)
}

func testBorrowBook(t *testing.T, client library_v1.LibraryServiceClient) {
	// Mock the request and response
	request := &connect.Request[v1.BorrowBookRequest]{
		Msg: &v1.BorrowBookRequest{
			BookId: "123", // Provide a book ID
			UserId: "456", // Provide a user ID
		},
	}
	expectedResponse := &connect.Response[v1.BorrowBookResponse]{
		Msg: &v1.BorrowBookResponse{
			Success: true,
			Message: "Book borrowed successfully",
			Book: &v1.Book{
				Id:     "123",
				Title:  "Sample Book",
				Author: "Sample Author",
			},
		},
	}

	// Call the method being tested
	response, err := client.BorrowBook(context.Background(), request)

	// Assert that no error occurred
	assert.NoError(t, err)

	assert.NotNil(t, response.Msg)
	assert.NotNil(t, response.Msg.Book)

	// Assert the response matches the expected response
	assert.Equal(t, expectedResponse.Msg.Success, response.Msg.Success)
	assert.Equal(t, expectedResponse.Msg.Book, response.Msg.Book)
}

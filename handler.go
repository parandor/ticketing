package ticketing

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	connect "connectrpc.com/connect"
	v1 "github.com/parandor/ticketing/internal/gen/proto/train_ticketing/v1"
	ticketingv1 "github.com/parandor/ticketing/internal/gen/proto/train_ticketing/v1/train_ticketingv1connect"
)

// MyTrainTicketingServiceHandler is an implementation of the TrainTicketingServiceHandler interface.
type MyTrainTicketingServiceHandler struct {
	users map[string]*v1.User // Map to store users by ID
	seats map[string]*v1.Seat // Map to store seats by ID
	mu    sync.Mutex          // Mutex to ensure safe access to the maps
}


func NewMyTicketingServiceHandler() (string, http.Handler) {
	handler := &MyTrainTicketingServiceHandler{
		users: make(map[string]*v1.User),
		seats: make(map[string]*v1.Seat)}

	for i := 1; i <= 20; i++ {
		seatID := fmt.Sprintf("%d", i)
		seat := &v1.Seat{SeatNumber: 0}
		handler.seats[seatID] = seat
	}
	
	// Use NewTicketingServiceHandler to create the HTTP handler
	path, httpHandler := ticketingv1.NewTrainTicketingServiceHandler(handler)

	// Apply middleware to intercept JWT tokens
	httpHandler = withJWTInterceptor(httpHandler)

	// Optionally, you can add middleware or modify the http.Handler here

	return path, httpHandler
}

func withJWTInterceptor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Your JWT token verification logic here
		// For example, you can check the "Authorization" header for the token
		token := r.Header.Get("Authorization")

		if token != "Bearer auth_token" {
			http.Error(w, "Unauthorized: No JWT token provided", http.StatusUnauthorized)
			return
		}

		// If the token is valid, you can proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

// PurchaseTicket implements the PurchaseTicket method of TrainTicketingServiceHandler.
func (h *MyTrainTicketingServiceHandler) PurchaseTicket(ctx context.Context, req *connect.Request[v1.PurchaseTicketRequest]) (*connect.Response[v1.PurchaseTicketResponse], error) {
	// Extract the PurchaseTicketRequest parameters from the req
	ticket := req.Msg.GetTicket()

	// Extract ticket information from the request
	user := ticket.GetUser()

	// Lock the mutex to ensure safe access to the maps
	h.mu.Lock()
	defer h.mu.Unlock()

	// Validate user
	if user == nil || user.GetFirstName() == "" || user.GetLastName() == "" || user.GetEmail() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("user information is invalid"))
	}

	var availableSeatID string
	for seatID, seat := range h.seats {
		if seat.GetSeatNumber() == 0 {
			availableSeatID = seatID
			break
		}
	}

	// Check if a seat is available
	if availableSeatID == "" {
		return nil, connect.NewError(connect.CodeResourceExhausted, errors.New("no available seats"))
	}

	seatNumber, err := strconv.Atoi(availableSeatID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to parse available seat ID"))
	}

	// Assign the seat to the user
	assignedSeat := h.seats[availableSeatID]
	assignedSeat.SeatNumber = int32(seatNumber)
	assignedSeat.User = user // Associate the user with the seat

	// Generate a receipt
	receipt := &v1.Receipt{
		Ticket: ticket,
	}

	// Return the response containing the receipt
	response := &v1.PurchaseTicketResponse{
		Receipt: receipt,
	}

	return connect.NewResponse(response), nil
}

// ViewReceipt implements the ViewReceipt method of TrainTicketingServiceHandler.
func (h *MyTrainTicketingServiceHandler) ViewReceipt(ctx context.Context, req *connect.Request[v1.ViewReceiptRequest]) (*connect.Response[v1.ViewReceiptResponse], error) {
	// Extract the ViewReceiptRequest parameters from the req
	ticket := req.Msg.GetTicket()

	// Lock the mutex to ensure safe access to the maps
	h.mu.Lock()
	defer h.mu.Unlock()

	// Retrieve the receipt for the provided ticket
	receipt, err := h.retrieveReceipt(ticket)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err) // Assuming receipt not found is treated as a not found error
	}

	// Create the ViewReceiptResponse containing the retrieved receipt
	response := &v1.ViewReceiptResponse{
		Receipt: receipt,
	}

	// Return the response
	return connect.NewResponse(response), nil
}

func (h *MyTrainTicketingServiceHandler) retrieveReceipt(ticket *v1.Ticket) (*v1.Receipt, error) {
	// Iterate over the seats to find the one with the matching ticket information
	for _, seat := range h.seats {
		if seat.GetSeatNumber() != 0 && seat.GetUser() != nil {
			// If the seat is occupied and the user matches the ticket information, return the receipt
			if seat.GetUser().GetFirstName() == ticket.GetUser().GetFirstName() &&
				seat.GetUser().GetLastName() == ticket.GetUser().GetLastName() &&
				seat.GetUser().GetEmail() == ticket.GetUser().GetEmail() {
				return &v1.Receipt{
					Ticket: ticket,
				}, nil
			}
		}
	}

	// If no matching receipt is found, return an error
	return nil, errors.New("receipt not found")
}

// ViewAdminDetails implements the ViewAdminDetails method of TrainTicketingServiceHandler.
func (h *MyTrainTicketingServiceHandler) ViewAdminDetails(ctx context.Context, req *connect.Request[v1.ViewAdminDetailsRequest]) (*connect.Response[v1.ViewAdminDetailsResponse], error) {
	// Lock the mutex to ensure safe access to the maps
	h.mu.Lock()
	defer h.mu.Unlock()

	// Initialize slices to store all users and seats
	var allUsers []*v1.User
	var allSeats []*v1.Seat

	// Iterate through the seats to collect all users and seats information
	for _, seat := range h.seats {
		if seat.GetUser() != nil {
			allSeats = append(allSeats, seat)
			allUsers = append(allUsers, seat.GetUser())
		}
	}

	// Create a response containing all users and seats information
	adminView := &v1.AdminView{
		Users: allUsers,
		Seats: allSeats,
	}

	// Return the response
	response := &v1.ViewAdminDetailsResponse{
		AdminView: adminView,
	}

	return connect.NewResponse(response), nil
}

// RemoveUser implements the RemoveUser method of TrainTicketingServiceHandler.
func (h *MyTrainTicketingServiceHandler) RemoveUser(ctx context.Context, req *connect.Request[v1.RemoveUserRequest]) (*connect.Response[v1.RemoveUserResponse], error) {
	// Extract the user's first name to be removed from the request
	firstName := req.Msg.GetUser().GetFirstName() // Assuming you have a FirstName field in the User message

	// Lock the mutex to ensure safe access to the maps
	h.mu.Lock()
	defer h.mu.Unlock()

	// Check if the user to be removed exists
	var found bool
	for _, user := range h.users {
		if user.GetFirstName() == firstName {
			found = true
			break
		}
	}
	if !found {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("user to be removed not found"))
	}

	// Remove the user
	delete(h.users, firstName)

	// Remove the association of the user with any seat
	for _, seat := range h.seats {
		if seat.GetUser() != nil && seat.GetUser().GetFirstName() == firstName {
			seat.SeatNumber = 0
			seat.User = &v1.User{
				FirstName: "",
				LastName:  "",
				Email:     "",
			}
		}
	}

	// Return a success response
	response := &v1.RemoveUserResponse{}
	return connect.NewResponse(response), nil
}

// ModifySeat implements the ModifySeat method of TrainTicketingServiceHandler.
func (h *MyTrainTicketingServiceHandler) ModifySeat(ctx context.Context, req *connect.Request[v1.ModifySeatRequest]) (*connect.Response[v1.ModifySeatResponse], error) {
	// Extract ModifySeatRequest parameters from the request
	modifyReq := req.Msg

	// Extract user and new seat information from the request
	user := modifyReq.GetUser()

	// Lock the mutex to ensure safe access to the maps
	h.mu.Lock()
	defer h.mu.Unlock()

	var availableSeatID string
	for seatID, seat := range h.seats {
		if seat.GetSeatNumber() == 0 {
			availableSeatID = seatID
			break
		}
	}
	seatNumber, err := strconv.ParseInt(availableSeatID, 10, 32)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to parse available seat ID"))
	}

	// Find the user by first name in the seats map
	found := false
	for _, seat := range h.seats {
		// Check if admin, for sake of time, replace admin ID check with first name
		if seat.User != nil && seat.User.GetFirstName() == user.GetFirstName() {
			// Update the seat information for the user
			seat.User = user
			seat.SeatNumber = int32(seatNumber)

			// Set found flag to true
			found = true
			break
		}
	}

	// If user is not found, return an error
	if !found {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("user not found in seats"))
	}

	// Return a success response
	return connect.NewResponse(&v1.ModifySeatResponse{}), nil
}

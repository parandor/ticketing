# Requirements:

1. Code must be published in Github with a link we can access (use public repo).

2. Code must compile with some effort on unit tests, doesn't have to be 100%, but it shouldn't be 0%.

3. Please code this with Golang and gRPC

4. No persistence layer is required, just store the data in the current session/in memory.

5. The results can be in the console output from your grpc-server and grpc-client 

6. Depending on the level of authentication, take different actions



App to be coded

Note: All APIs referenced are gRPC APIs, not REST ones.

I want to board a train from London to France. The train ticket will cost $20, regardless of section or seat.

1. Authenticated APIs should be able to parse a JWT, formatted as if from an OAuth2 server, from the metadata to authenticate a request. No signature validation is required.

2. Create a public API where you can submit a purchase for a ticket. Details included in the receipt are:

a. From, To, User , Price Paid.

i. User should include first name, last name, email address

b. The user is allocated a seat in the train as a result of the purchase. Assume the train has only 2 sections, section A and section B and each section has 10 seats.

3. An authenticated API that shows the details of the receipt for the user

4. An authenticated API that lets an admin view all the users and seats they are allocated by the requested section

5. An authenticated API to allow an admin or the user to remove the user from the train

6. An authenticated API to allow an admin or the user to modify the user's seat

# Test

Run tests with: 
```
go test -v ./...
```

## CICD

Test action script lives [here](\.github/workflows/test.yaml)
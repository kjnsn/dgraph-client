package dgraph

import (
	"context"

	"github.com/dgraph-io/dgraph/client"
	"github.com/dgraph-io/dgraph/protos/graphp"
	"google.golang.org/grpc"
)

// Connection represents a connection to Dgraph
type Connection struct {
	dclient  graphp.DgraphClient
	grpcConn *grpc.ClientConn
}

// Connect creates a new connection to the Dgraph instance
// given by address.
func Connect(address string) (*Connection, error) {
	cc, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	conn := graphp.NewDgraphClient(cc)
	return &Connection{
		dclient:  conn,
		grpcConn: cc,
	}, nil
}

// Close the connection
func (c *Connection) Close() {
	c.grpcConn.Close()
}

// Query sends a query string to dgraph and returns the response
func (c *Connection) Query(query string) (*Response, error) {
	return c.QueryVariables(query, make(map[string]string))
}

// QueryVariables sends the given query string to dgraph, with the variables given.
func (c *Connection) QueryVariables(query string, variables map[string]string) (*Response, error) {
	req := client.Req{}
	req.SetQuery(query, variables)
	resp, err := c.dclient.Run(context.Background(), req.Request())
	if err != nil {
		return nil, err
	}
	return ReadResponse(resp), nil
}

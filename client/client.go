package client

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	// DefaultName prefix client name
	DefaultName = "client"
)

// A Client defines parameters for making calls to an HTTP server.
// The zero value for Client is a valid configuration.
type Client struct {
	cfg    Config
	health *health.Server
	logger zerolog.Logger
}

// New create a new client
func New(opts ...Option) *Client {
	return newClient(opts...)
}

// newClient will instantiate a new Client with the given config
func newClient(opts ...Option) *Client {
	c := &Client{
		cfg:    newConfig(),
		health: health.NewServer(),
	}
	c.logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	if len(opts) > 0 {
		c = c.WithOptions(opts...)
	}
	return c
}

// WithOptions clones the current Client, applies the supplied Options, and
// returns the resulting Client. It's safe to use concurrently.
func (c *Client) WithOptions(opts ...Option) *Client {
	for _, opt := range opts {
		opt.apply(c)
	}
	return c
}

func (c *Client) applyOptions(opts ...Option) {
	for _, opt := range opts {
		opt.apply(c)
	}
}

// Init initialize logger and interceptors
func (c *Client) Init(opts ...Option) {
	c.applyOptions(opts...)
	c.logger = c.logger.With().Dict("client", zerolog.Dict().
		Str("id", c.cfg.ID).
		Str("name", c.cfg.Name).
		Str("format", c.cfg.Format).
		Str("target", c.cfg.Target),
	).Logger()
	c.logger.Info().Msg(fmt.Sprintf("starting %s %s, connecting to %s", c.cfg.Format, c.Name(), c.cfg.Target))
	// append dial interceptor to grpc client
	c.cfg.mu.Lock()
	defer c.cfg.mu.Unlock()
	c.cfg.UnaryClientInterceptors = append(c.cfg.UnaryClientInterceptors, dialUnaryClientInterceptor(c))
	c.cfg.StreamClientInterceptors = append(c.cfg.StreamClientInterceptors, dialStreamClientInterceptor(c))
}

// Dial create a gRPC channel to communicate with the server
func (c *Client) Dial(opts ...Option) (*grpc.ClientConn, error) {
	c.applyOptions(opts...)
	// TODO use TLS
	c.cfg.mu.Lock()
	defer c.cfg.mu.Unlock()
	conn, err := grpc.Dial(
		c.cfg.Target,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(c.cfg.Timeout),
		grpc.WithUnaryInterceptor(WrapperUnaryClient(c.cfg.UnaryClientInterceptors...)),
		grpc.WithStreamInterceptor(WrapperStreamClient(c.cfg.StreamClientInterceptors...)),
	)
	switch grpc.Code(err) {
	case codes.OK:
		break
	default:
		return nil, err
	}
	return conn, nil
}

// Dial create a gRPC channel to communicate with the server
func (c *Client) DialWithCredentials(token string, opts ...Option) (*grpc.ClientConn, error) {
	c.applyOptions(opts...)
	// TODO use TLS
	c.cfg.mu.Lock()
	defer c.cfg.mu.Unlock()
	conn, err := grpc.Dial(
		c.cfg.Target,
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(TokenAuth{
			token: token,
		}),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(c.cfg.Timeout),
		grpc.WithUnaryInterceptor(WrapperUnaryClient(c.cfg.UnaryClientInterceptors...)),
		grpc.WithStreamInterceptor(WrapperStreamClient(c.cfg.StreamClientInterceptors...)),
	)
	switch grpc.Code(err) {
	case codes.OK:
		break
	default:
		return nil, err
	}
	return conn, nil
}

// Stub to perform RPCs
func (c *Client) Stub(conn *grpc.ClientConn) interface{} {
	return c.cfg.GRPCRegister(conn)
}

// Close not implemented
func (c *Client) Close() error {
	return nil
}

// Name returns client name
func (c *Client) Name() string {
	return c.cfg.Name
}

func (c *Client) healthRPC() {
	conn, err := c.Dial()
	if err != nil {
		c.logger.Error().Msg(err.Error())
		c.health.SetServingStatus(c.cfg.ID, healthpb.HealthCheckResponse_NOT_SERVING)
		return
	}
	defer conn.Close()

	h := healthpb.NewHealthClient(conn)
	hcr, err := h.Check(context.Background(), &healthpb.HealthCheckRequest{})
	if err != nil {
		c.logger.Error().Msg(err.Error())
		c.health.SetServingStatus(c.cfg.ID, healthpb.HealthCheckResponse_NOT_SERVING)
		return
	}
	c.health.SetServingStatus(c.cfg.ID, hcr.Status)
}

// Health health check on client take in consideration
// a round trip to a server
func (c *Client) Health() *healthpb.HealthCheckResponse {
	// perform health check RPC
	c.healthRPC()
	// check client status
	hcr, err := c.health.Check(
		context.Background(),
		&healthpb.HealthCheckRequest{Service: c.cfg.ID},
	)
	if err != nil {
		c.logger.Error().Msg(err.Error())
		return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_NOT_SERVING}
	}
	return hcr
}

// Token based authentication
type TokenAuth struct {
	token string
}

// Return value is mapped to request headers.
func (t TokenAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + t.token,
	}, nil
}

func (TokenAuth) RequireTransportSecurity() bool {
	return false
}

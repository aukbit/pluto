package client

import (
	"go.uber.org/zap"

	g "github.com/aukbit/pluto/client/grpc"
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
	logger *zap.Logger // client logger
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
	c.logger, _ = zap.NewProduction()
	if len(opts) > 0 {
		c = c.WithOptions(opts...)
	}
	return c
}

// WithOptions clones the current Client, applies the supplied Options, and
// returns the resulting Client. It's safe to use concurrently.
func (c *Client) WithOptions(opts ...Option) *Client {
	d := c.clone()
	for _, opt := range opts {
		opt.apply(d)
	}
	return d
}

func (c *Client) applyOptions(opts ...Option) {
	if len(opts) > 0 {
		for _, opt := range opts {
			opt.apply(c)
		}
	}
}

// clone creates a shallow copy client
func (c *Client) clone() *Client {
	copy := *c
	return &copy
}

// Init initialize logger and interceptors
func (c *Client) Init(opts ...Option) {
	c.applyOptions(opts...)
	// set logger
	c.logger = c.logger.With(
		zap.String("id", c.cfg.ID),
		zap.String("name", c.cfg.Name),
		zap.String("format", c.cfg.Format),
	)
	// append dial interceptor to grpc client
	c.cfg.UnaryClientInterceptors = append(c.cfg.UnaryClientInterceptors, dialUnaryClientInterceptor(c))
}

// Dial create a gRPC channel to communicate with the server
func (c *Client) Dial(opts ...Option) (*grpc.ClientConn, error) {
	c.applyOptions(opts...)
	// TODO use TLS
	conn, err := grpc.Dial(
		c.cfg.Target,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(c.cfg.Timeout),
		grpc.WithUnaryInterceptor(g.WrapperUnaryClient(c.cfg.UnaryClientInterceptors...)),
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
		c.logger.Error("healthRPC", zap.String("err", err.Error()))
		c.health.SetServingStatus(c.cfg.ID, healthpb.HealthCheckResponse_NOT_SERVING)
		return
	}
	defer conn.Close()

	h := healthpb.NewHealthClient(conn)
	hcr, err := h.Check(context.Background(), &healthpb.HealthCheckRequest{})
	if err != nil {
		c.logger.Error("healthRPC", zap.String("err", err.Error()))
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
		c.logger.Error("Health", zap.String("err", err.Error()))
		return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_NOT_SERVING}
	}
	return hcr
}

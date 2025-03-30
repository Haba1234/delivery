package geo

import (
	"context"
	"log"
	"time"

	"github.com/Haba1234/delivery/internal/core/domain/model/kernel"
	"github.com/Haba1234/delivery/internal/core/ports"
	"github.com/Haba1234/delivery/internal/pkg/errs"
	pb "github.com/Haba1234/delivery/pkg/clients/geo/geosrv/geopb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _ ports.IGeoClient = &Client{}

type Client struct {
	conn     *grpc.ClientConn
	pbClient pb.GeoClient
	timeout  time.Duration
}

func NewClient(host string) (*Client, error) {
	if host == "" {
		return nil, errs.NewValueIsRequiredError("host")
	}

	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	pbClient := pb.NewGeoClient(conn)

	return &Client{
		conn:     conn,
		pbClient: pbClient,
		timeout:  5 * time.Second,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetGeolocation(ctx context.Context, street string) (kernel.Location, error) {
	// Формируем запрос
	req := &pb.GetGeolocationRequest{
		Street: street,
	}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.pbClient.GetGeolocation(ctx, req)
	if err != nil {
		return kernel.Location{}, err
	}

	location, err := kernel.CreateLocation(int(resp.Location.X), int(resp.Location.Y))
	if err != nil {
		return kernel.Location{}, err
	}
	return location, nil
}

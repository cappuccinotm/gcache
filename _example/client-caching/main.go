package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"

	"github.com/cappuccinotm/gcache"
	"github.com/cappuccinotm/gcache/_example/client-caching/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var modeFlag = flag.String("mode", "server", "run mode")

func main() {
	flag.Parse()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	switch *modeFlag {
	case "server":
		if err := runServer(ctx); err != nil {
			log.Fatalf("failed to run server: %v", err)
		}
	case "client":
		if err := runClient(ctx); err != nil {
			log.Fatalf("failed to run client: %v", err)
		}
	default:
		log.Fatalf("unknown mode: %s", *modeFlag)
	}
}

func runServer(ctx context.Context) error {
	srv := grpc.NewServer()
	order.RegisterOrderServiceServer(srv, &server{})

	log.Printf("starting server")
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	go func() {
		<-ctx.Done()
		log.Printf("shutting down server")
		srv.GracefulStop()
	}()

	if err = srv.Serve(l); err != nil {
		return fmt.Errorf("serve: %w", err)
	}

	return nil
}

type server struct {
	order.UnimplementedOrderServiceServer
}

const timeStamp = "2021-01-01T00:00:00Z"

func (s *server) GetOrder(ctx context.Context, req *order.GetOrderRequest) (*order.Order, error) {
	etag := gcache.ETag(ctx)
	log.Printf("requested order, etag: %s", etag)
	if etag == timeStamp {
		log.Printf("skipping, not modified")
		return nil, gcache.NotChanged(ctx, timeStamp)
	}

	log.Printf("delivering actual order")

	if err := grpc.SendHeader(ctx, metadata.Pairs("ETag", timeStamp)); err != nil {
		return nil, fmt.Errorf("send header: %w", err)
	}

	return &order.Order{
		Id:         "order-1",
		CustomerId: "customer-1",
		Items: []*order.Order_OrderItem{{
			Id:       "item-1",
			Name:     "Pizza",
			Quantity: 1,
		}},
	}, nil
}

func runClient(ctx context.Context) error {
	icptr := gcache.NewInterceptor(gcache.WithLogger(slog.Default()))
	conn, err := grpc.NewClient("localhost:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(icptr.UnaryClientInterceptor()),
	)
	if err != nil {
		return fmt.Errorf("dial localhost:8080: %w", err)
	}

	client := order.NewOrderServiceClient(conn)
	resp, err := client.GetOrder(ctx, &order.GetOrderRequest{})
	if err != nil {
		return fmt.Errorf("get order: %w", err)
	}

	log.Printf("order: %v", resp)

	if resp, err = client.GetOrder(ctx, &order.GetOrderRequest{}); err != nil {
		return fmt.Errorf("get order: %w", err)
	}

	log.Printf("order: %v", resp)

	return nil
}

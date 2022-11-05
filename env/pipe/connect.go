package pipe

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
)

func ConnectToServer(ctx context.Context, s *grpc.Server) (*grpc.ClientConn, func()) {

	pipe := ListenPipe()

	go func() {
		fmt.Printf("Serving on '%s'\n", pipe.Addr().String())
		if err := s.Serve(pipe); err != nil {
			if err != ErrPipeListenerClosed {
				log.Fatalf("failed to serve: %v", err)
			}
		}
	}()

	fmt.Printf("Dialing '%s'\n", pipe.Addr())
	clientConn, err := grpc.DialContext(ctx, `sim`,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(c context.Context, s string) (net.Conn, error) {
			return pipe.DialContext(c, `sim`, s)
		}),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return clientConn, func() {
		e := pipe.Close()
		if e != nil {
			log.Fatalf("failed to close: %s", e)
		}
	}

}

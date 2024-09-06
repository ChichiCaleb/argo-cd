package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"runtime/debug"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"

	"github.com/argoproj/argo-cd/v2/common"
)

// PanicLoggerUnaryServerInterceptor returns a new unary server interceptor for recovering from panics and returning error
func PanicLoggerUnaryServerInterceptor(log *logrus.Entry) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("Recovered from panic: %+v\n%s", r, debug.Stack())
				err = status.Errorf(codes.Internal, "%s", r)
			}
		}()
		return handler(ctx, req)
	}
}

// PanicLoggerStreamServerInterceptor returns a new streaming server interceptor for recovering from panics and returning error
func PanicLoggerStreamServerInterceptor(log *logrus.Entry) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("Recovered from panic: %+v\n%s", r, debug.Stack())
				err = status.Errorf(codes.Internal, "%s", r)
			}
		}()
		return handler(srv, stream)
	}
}

// BlockingDial is a helper method to create a gRPC "channel" using the new gRPC `NewClient` function.
// The connection will be established as needed when the returned `ClientConn` is used for RPCs.
// Adapted from: https://github.com/fullstorydev/grpcurl/blob/master/grpcurl.go
func BlockingDial(ctx context.Context, network, address string, creds credentials.TransportCredentials, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
    // Custom dialer to handle TLS and provide better error messages
    dialer := func(ctx context.Context, address string) (net.Conn, error) {
        conn, err := proxy.Dial(ctx, network, address)
        if err != nil {
            return nil, fmt.Errorf("error dial proxy: %w", err)
        }
        if creds != nil {
            conn, _, err = creds.ClientHandshake(ctx, address, conn)
            if err != nil {
                return nil, fmt.Errorf("error creating connection: %w", err)
            }
        }
        return conn, nil
    }

    // Configure the gRPC dial options
    opts = append(opts,
        grpc.WithContextDialer(dialer),
        grpc.WithTransportCredentials(insecure.NewCredentials()), // We are handling TLS, so tell grpc not to
        grpc.WithKeepaliveParams(keepalive.ClientParameters{Time: common.GetGRPCKeepAliveTime()}),
    )

    // Use `NewClient` to create a gRPC "channel"
    conn, err := grpc.NewClient(address, opts...)
    if err != nil {
        return nil, fmt.Errorf("failed to create gRPC client: %w", err)
    }

    // Check if the context is done before proceeding
    select {
    case <-ctx.Done():
        conn.Close()
        return nil, ctx.Err()
    default:
    }

    // Attempt to establish the connection manually 
    conn.Connect() 

    return conn, nil
}

type TLSTestResult struct {
	TLS         bool
	InsecureErr error
}

func TestTLS(address string, dialTime time.Duration) (*TLSTestResult, error) {
	if parts := strings.Split(address, ":"); len(parts) == 1 {
		// If port is unspecified, assume the most likely port
		address += ":443"
	}
	var testResult TLSTestResult
	var tlsConfig tls.Config
	tlsConfig.InsecureSkipVerify = true
	creds := credentials.NewTLS(&tlsConfig)

	// Set timeout when dialing to the server
	// fix: https://github.com/argoproj/argo-cd/issues/9679
	ctx, cancel := context.WithTimeout(context.Background(), dialTime)
	defer cancel()

	conn, err := BlockingDial(ctx, "tcp", address, creds)
	if err == nil {
		_ = conn.Close()
		testResult.TLS = true
		creds := credentials.NewTLS(&tls.Config{})
		ctx, cancel := context.WithTimeout(context.Background(), dialTime)
		defer cancel()

		conn, err := BlockingDial(ctx, "tcp", address, creds)
		if err == nil {
			_ = conn.Close()
		} else {
			// if connection was successful with InsecureSkipVerify true, but unsuccessful with
			// InsecureSkipVerify false, it means server is not configured securely
			testResult.InsecureErr = err
		}
		return &testResult, nil
	}
	// If we get here, we were unable to connect via TLS (even with InsecureSkipVerify: true)
	// It may be because server is running without TLS, or because of real issues (e.g. connection
	// refused). Test if server accepts plain-text connections
	ctx, cancel = context.WithTimeout(context.Background(), dialTime)
	defer cancel()
	conn, err = BlockingDial(ctx, "tcp", address, nil)
	if err == nil {
		_ = conn.Close()
		testResult.TLS = false
		return &testResult, nil
	}
	return nil, err
}

func WithTimeout(duration time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		clientDeadline := time.Now().Add(duration)
		ctx, cancel := context.WithDeadline(ctx, clientDeadline)
		defer cancel()
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

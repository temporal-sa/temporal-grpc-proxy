package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net"
	"os"

	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// command line flags
var bindAddress string
var upstreamAddress string
var clientCert string
var clientKey string

func init() {
	flag.StringVar(&bindAddress, "bindAddress", "localhost:7233", "Address to bind to")
	flag.StringVar(&upstreamAddress, "upstreamAddress", os.Getenv("TEMPORAL_ADDRESS"), "host:port for upstream Temporal frontend service [$TEMPORAL_ADDRESS]")
	flag.StringVar(&clientCert, "tls_cert_path", os.Getenv("TEMPORAL_TLS_CERT"), "Path to client x509 certificate [$TEMPORAL_TLS_CERT]")
	flag.StringVar(&clientKey, "tls_key_path", os.Getenv("TEMPORAL_TLS_KEY"), "Path to client certificate private key [$TEMPORAL_TLS_KEY]")
}

func main() {
	flag.Parse()

	// parse flags, load certs and create transport credentials for grpc client
	transportCredentials := insecure.NewCredentials()
	if clientCert != "" && clientKey != "" {
		cert, err := tls.LoadX509KeyPair(clientCert, clientKey)
		if err != nil {
			log.Fatalf("failed to load client cert and key: %v", err)
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}

		transportCredentials = credentials.NewTLS(tlsConfig)
	}

	// create grpc client
	grpcClient, err := grpc.Dial(
		upstreamAddress,
		grpc.WithTransportCredentials(transportCredentials),
	)
	if err != nil {
		log.Fatalf("failed to create grpc client: %v", err)
	}
	defer grpcClient.Close()

	// create workflow client from grpc client
	workflowClient := workflowservice.NewWorkflowServiceClient(grpcClient)
	if err != nil {
		log.Fatalf("failed to create workflow client: %v", err)
	}

	// create tcp listener on bind port
	listener, err := net.Listen("tcp", bindAddress)
	if err != nil {
		log.Fatalf("failed to create tcp listener: %v", err)
	}

	// create grpc server and register workflow service proxy handler
	server := grpc.NewServer()
	handler, err := client.NewWorkflowServiceProxyServer(
		client.WorkflowServiceProxyOptions{Client: workflowClient},
	)
	if err != nil {
		log.Fatalf("failed to create service proxy: %v", err)
	}
	workflowservice.RegisterWorkflowServiceServer(server, handler)

	// create server for health checks
	healthServer := health.NewServer()
	healthServer.SetServingStatus("temporal.api.workflowservice.v1.WorkflowService", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(server, healthServer)

	// start grpc server
	log.Println("proxy server listening on", listener.Addr().String())
	err = server.Serve(listener)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

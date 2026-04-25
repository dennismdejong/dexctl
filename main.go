package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/dexidp/dex/api/v2"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Version is set via ldflags during build
var version string = "dev"

var (
	serverAddr string
	certFile   string
	keyFile    string
	caFile     string
	insecure   bool
	outputJSON bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "dexctl",
		Short: "Dex client CLI",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	rootCmd.PersistentFlags().StringVar(&serverAddr, "server", "localhost:5556", "Dex gRPC server address")
	rootCmd.PersistentFlags().StringVar(&certFile, "cert", "", "TLS certificate file")
	rootCmd.PersistentFlags().StringVar(&keyFile, "key", "", "TLS key file")
	rootCmd.PersistentFlags().StringVar(&caFile, "ca", "", "TLS CA certificate file")
	rootCmd.PersistentFlags().BoolVar(&insecure, "insecure", false, "Skip TLS certificate verification")
	rootCmd.PersistentFlags().BoolVar(&outputJSON, "json", false, "Output as JSON")

	// Version command
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of dexctl",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("dexctl version %s\n", version)
		},
	}
	rootCmd.AddCommand(versionCmd)

	// Client commands
	clientCmd := &cobra.Command{
		Use:   "client",
		Short: "Manage Dex clients",
	}
	clientCmd.AddCommand(newClientCreateCmd())
	clientCmd.AddCommand(newClientGetCmd())
	clientCmd.AddCommand(newClientUpdateCmd())
	clientCmd.AddCommand(newClientDeleteCmd())
	clientCmd.AddCommand(newClientListCmd())

	rootCmd.AddCommand(clientCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newGRPCConn() (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	if insecure {
		opts = append(opts, grpc.WithInsecure())
	} else {
		var creds credentials.TransportCredentials
		if certFile != "" && keyFile != "" && caFile != "" {
			// Load client certs
			cert, err := tls.LoadX509KeyPair(certFile, keyFile)
			if err != nil {
				return nil, fmt.Errorf("failed to load client certificate: %w", err)
			}
			// Load CA
			caCert, err := os.ReadFile(caFile)
			if err != nil {
				return nil, fmt.Errorf("failed to read CA file: %w", err)
			}
			caPool := x509.NewCertPool()
			if ok := caPool.AppendCertsFromPEM(caCert); !ok {
				return nil, fmt.Errorf("failed to append CA certs")
			}
			tlsConfig := &tls.Config{
				Certificates: []tls.Certificate{cert},
				RootCAs:      caPool,
			}
			creds = credentials.NewTLS(tlsConfig)
		} else if caFile != "" {
			// Only CA cert (for server verification)
			caCert, err := os.ReadFile(caFile)
			if err != nil {
				return nil, fmt.Errorf("failed to read CA file: %w", err)
			}
			caPool := x509.NewCertPool()
			if ok := caPool.AppendCertsFromPEM(caCert); !ok {
				return nil, fmt.Errorf("failed to append CA certs")
			}
			creds = credentials.NewClientTLSFromCert(caPool, "")
		} else {
			// No TLS provided, treat as insecure (plaintext)
			opts = append(opts, grpc.WithInsecure())
		}
		if creds != nil {
			opts = append(opts, grpc.WithTransportCredentials(creds))
		}
	}
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", serverAddr, err)
	}
	return conn, nil
}

func printOutput(v interface{}) error {
	if outputJSON {
		data, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}
	// Default: print in a simple format
	switch v := v.(type) {
	case *api.Client:
		fmt.Printf("ID: %s\n", v.GetId())
		fmt.Printf("Name: %s\n", v.GetName())
		fmt.Printf("Public: %v\n", v.GetPublic())
		fmt.Printf("RedirectURIs: %v\n", v.GetRedirectUris())
		fmt.Printf("Secret: %s\n", v.GetSecret())
		fmt.Printf("TrustedPeers: %v\n", v.GetTrustedPeers())
		fmt.Printf("LogoURL: %s\n", v.GetLogoUrl())
	case []*api.ClientInfo:
		if len(v) == 0 {
			fmt.Println("No clients found")
			return nil
		}
		for i, c := range v {
			if i > 0 {
				fmt.Println("---")
			}
			fmt.Printf("ID: %s\n", c.GetId())
			fmt.Printf("Name: %s\n", c.GetName())
			fmt.Printf("Public: %v\n", c.GetPublic())
			fmt.Printf("RedirectURIs: %v\n", c.GetRedirectUris())
			fmt.Printf("TrustedPeers: %v\n", c.GetTrustedPeers())
			fmt.Printf("LogoURL: %s\n", c.GetLogoUrl())
		}
	default:
		fmt.Println(v)
	}
	return nil
}

// Client command constructors

func newClientCreateCmd() *cobra.Command {
	var (
		name        string
		public      bool
		redirectURIs []string
		secret      string
		trustedPeers []string
		logoURL     string
	)
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new client",
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := newGRPCConn()
			if err != nil {
				return err
			}
			defer conn.Close()
			client := api.NewDexClient(conn)

			req := &api.CreateClientReq{
				Client: &api.Client{
					Id:           "", // Let server generate
					Name:         name,
					Public:       public,
					RedirectUris: redirectURIs,
					Secret:       secret,
					TrustedPeers: trustedPeers,
					LogoUrl:      logoURL,
				},
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			resp, err := client.CreateClient(ctx, req)
			if err != nil {
				return fmt.Errorf("create client failed: %w", err)
			}
			return printOutput(resp.GetClient())
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "Client name (required)")
	cmd.Flags().BoolVar(&public, "public", false, "Public client")
	cmd.Flags().StringSliceVar(&redirectURIs, "redirect-uri", []string{}, "Redirect URIs (can be specified multiple times)")
	cmd.Flags().StringVar(&secret, "secret", "", "Client secret")
	cmd.Flags().StringSliceVar(&trustedPeers, "trusted-peer", []string{}, "Trusted peers (can be specified multiple times)")
	cmd.Flags().StringVar(&logoURL, "logo-url", "", "Logo URL")
	cmd.MarkFlagRequired("name")
	return cmd
}

func newClientGetCmd() *cobra.Command {
	var id string
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a client by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := newGRPCConn()
			if err != nil {
				return err
			}
			defer conn.Close()
			client := api.NewDexClient(conn)

			req := &api.GetClientReq{
				Id: id,
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			resp, err := client.GetClient(ctx, req)
			if err != nil {
				return fmt.Errorf("get client failed: %w", err)
			}
			return printOutput(resp.GetClient())
		},
	}
	cmd.Flags().StringVar(&id, "id", "", "Client ID (required)")
	cmd.MarkFlagRequired("id")
	return cmd
}

func newClientUpdateCmd() *cobra.Command {
	var (
		id          string
		name        string
		redirectURIs []string
		trustedPeers []string
		logoURL     string
	)
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a client by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := newGRPCConn()
			if err != nil {
				return err
			}
			defer conn.Close()
			client := api.NewDexClient(conn)

			req := &api.UpdateClientReq{
				Id: id,
				// Note: The UpdateClientReq in the dex api does not have a Secret field.
				// To update the secret, one would need to use a different method (if available) or recreate.
				Name:        name,
				RedirectUris: redirectURIs,
				TrustedPeers: trustedPeers,
				LogoUrl:      logoURL,
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			resp, err := client.UpdateClient(ctx, req)
			if err != nil {
				return fmt.Errorf("update client failed: %w", err)
			}
			if resp.GetNotFound() {
				return fmt.Errorf("client not found")
			}
			// Fetch and display the updated client
			getReq := &api.GetClientReq{Id: id}
			getResp, err := client.GetClient(context.Background(), getReq)
			if err != nil {
				return fmt.Errorf("failed to fetch updated client: %w", err)
			}
			return printOutput(getResp.GetClient())
		},
	}
	cmd.Flags().StringVar(&id, "id", "", "Client ID (required)")
	cmd.Flags().StringVar(&name, "name", "", "New client name")
	cmd.Flags().StringSliceVar(&redirectURIs, "redirect-uri", []string{}, "New redirect URIs (can be specified multiple times)")
	cmd.Flags().StringSliceVar(&trustedPeers, "trusted-peer", []string{}, "New trusted peers (can be specified multiple times)")
	cmd.Flags().StringVar(&logoURL, "logo-url", "", "New logo URL")
	cmd.MarkFlagRequired("id")
	return cmd
}

func newClientDeleteCmd() *cobra.Command {
	var id string
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a client by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := newGRPCConn()
			if err != nil {
				return err
			}
			defer conn.Close()
			client := api.NewDexClient(conn)

			req := &api.DeleteClientReq{
				Id: id,
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			resp, err := client.DeleteClient(ctx, req)
			if err != nil {
				return fmt.Errorf("delete client failed: %w", err)
			}
			if resp.GetNotFound() {
				return fmt.Errorf("client not found")
			}
			fmt.Println("Client deleted successfully")
			return nil
		},
	}
	cmd.Flags().StringVar(&id, "id", "", "Client ID (required)")
	cmd.MarkFlagRequired("id")
	return cmd
}

func newClientListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all clients",
		RunE: func(cmd *cobra.Command, args []string) error {
			conn, err := newGRPCConn()
			if err != nil {
				return err
			}
			defer conn.Close()
			client := api.NewDexClient(conn)

			req := &api.ListClientReq{}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			resp, err := client.ListClients(ctx, req)
			if err != nil {
				return fmt.Errorf("list clients failed: %w", err)
			}
			return printOutput(resp.GetClients())
		},
	}
	return cmd
}
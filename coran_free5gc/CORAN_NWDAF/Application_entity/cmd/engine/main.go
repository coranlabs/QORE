package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/coranlabs/CORAN_NWDAF/Application_entity/logger"
	engine "github.com/coranlabs/CORAN_NWDAF/Application_entity/server/engine"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type MainConfig struct {
	Server struct {
		// Addr string `envconfig:"SERVER_ADDR"`
	}
}

func main() {
	logger.InitializeLogger(logrus.InfoLevel)

	// Initialize internal engine package
	interfaceName := "eth0"
	err := godotenv.Load("config/nwdaf.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var config MainConfig
	err = envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err.Error())
	}
	engine.InitConfig()

	// Create router for engine
	engineRouter := engine.NewRouter()
	ip, err := getIPAddress(interfaceName)
	if err != nil {
		log.Fatalf("Failed to get IP address from interface %s: %v", interfaceName, err)
	}

	ip += ":8000"
	// Start server for engine
	startServer(ip, engineRouter)
}

func getIPAddress(interfaceName string) (string, error) {
	iface, err := net.InterfaceByName(interfaceName) // Get interface by name (e.g., "eth0")
	if err != nil {
		return "", err
	}

	addrs, err := iface.Addrs() // Get a list of addresses associated with the interface
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		// Check if the address is an IP address and is not a loopback address
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		if ip != nil && !ip.IsLoopback() && ip.To4() != nil {
			// Return the first valid IPv4 address
			return ip.String(), nil
		}
	}
	return "", fmt.Errorf("no valid IP address found for interface %s", interfaceName)
}

// Helper function to start a server
func startServer(addr string, handler http.Handler) {
	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.EngineLog.Infof("Engine Server listening at %s", addr)
	log.Fatal(server.ListenAndServe())
}

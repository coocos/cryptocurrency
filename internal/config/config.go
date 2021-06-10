package config

import "os"

// SeedHost returns the seed host for the node, i.e. the node to synchronize blockchain state from
func SeedHost() (string, bool) {
	return os.LookupEnv("NODE_SEED_HOST")
}

// BindHost returns the address to bind the API to listen to
func BindHost() string {
	if host, ok := os.LookupEnv("NODE_BIND_HOST"); ok {
		return host
	}
	return "localhost:8000"
}

// AdvertisedHost returns the address used by other nodes to communicate with this node
func AdvertisedHost() string {
	if host, ok := os.LookupEnv("NODE_ADVERTISED_HOST"); ok {
		return host
	}
	return BindHost()
}

package Nnrf_Bootstrapping

// APIClient manages communication with the NRF Bootstrapping API
type APIClient struct {
	cfg    *Configuration
	common service // Reuse a single struct instead of allocating one for each service on the heap.

	// API Services
	BootstrappingApi *BootstrappingRequestApiService
}


type service struct {
	client *APIClient
} 
// NewAPIClient creates a new API client
func NewAPIClient(cfg *Configuration) *APIClient {
	c := &APIClient{}
	c.cfg = cfg
	c.common.client = c

	// API Services
	c.BootstrappingApi = (*BootstrappingRequestApiService)(&c.common)

	return c
}

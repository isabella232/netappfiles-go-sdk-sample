package models

// AzureAuthInfo object definition
type AzureAuthInfo struct {
	ClientID                       *string
	ClientSecret                   *string
	SubscriptionID                 *string
	TenantID                       *string
	ActiveDirectoryEndpointURL     *string
	ResourceManagerEndpointURL     *string
	ActiveDirectoryGraphResourceID *string
	SqlManagementEndpointURL       *string
	GalleryEndpointURL             *string
	ManagementEndpointURL          *string
}

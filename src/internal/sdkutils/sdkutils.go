// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

// This package centralizes any function that directly
// using any of the Azure's (with exception of authentication related ones)
// available SDK packages.

package sdkutils

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure-Samples/netappfiles-go-sdk-sample/internal/iam"
	"github.com/Azure-Samples/netappfiles-go-sdk-sample/internal/uri"
	"github.com/Azure-Samples/netappfiles-go-sdk-sample/internal/utils"

	"github.com/Azure/azure-sdk-for-go/services/netapp/mgmt/2019-08-01/netapp"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-05-01/resources"
	"github.com/Azure/go-autorest/autorest/to"
)

const (
	userAgent = "anf-sdk-sample-agent"
	nfsv3     = "NFSv3"
	nfsv41    = "NFSv4.1"
	cifs      = "CIFS"
)

var (
	validProtocols = []string{nfsv3, nfsv41, "CIFS"}
)

func getResourcesClient() (resources.Client, error) {

	authorizer, subscriptionID, err := iam.GetAuthorizer()
	if err != nil {
		return resources.Client{}, err
	}

	client := resources.NewClient(subscriptionID)
	client.Authorizer = authorizer
	client.AddToUserAgent(userAgent)

	return client, nil
}

func getAccountsClient() (netapp.AccountsClient, error) {

	authorizer, subscriptionID, err := iam.GetAuthorizer()
	if err != nil {
		return netapp.AccountsClient{}, err
	}

	client := netapp.NewAccountsClient(subscriptionID)
	client.Authorizer = authorizer
	client.AddToUserAgent(userAgent)

	return client, nil
}

func getPoolsClient() (netapp.PoolsClient, error) {

	authorizer, subscriptionID, err := iam.GetAuthorizer()
	if err != nil {
		return netapp.PoolsClient{}, err
	}

	client := netapp.NewPoolsClient(subscriptionID)
	client.Authorizer = authorizer
	client.AddToUserAgent(userAgent)

	return client, nil
}

func getVolumesClient() (netapp.VolumesClient, error) {

	authorizer, subscriptionID, err := iam.GetAuthorizer()
	if err != nil {
		return netapp.VolumesClient{}, err
	}

	client := netapp.NewVolumesClient(subscriptionID)
	client.Authorizer = authorizer
	client.AddToUserAgent(userAgent)

	return client, nil
}

// GetResourceByID gets a generic resource
func GetResourceByID(ctx context.Context, resourceID, APIVersion string) (resources.GenericResource, error) {

	resourcesClient, err := getResourcesClient()
	if err != nil {
		return resources.GenericResource{}, err
	}

	parentResource := ""
	resourceGroup, _ := uri.GetResourceGroup(resourceID)
	resourceProvider, _ := uri.GetResourceValue(resourceID, "providers")
	resourceName, _ := uri.GetResourceName(resourceID)
	resourceType, _ := uri.GetResourceValue(resourceID, resourceProvider)

	if strings.Contains(resourceID, "/subnets/") {
		parentResourceName, _ := uri.GetResourceValue(resourceID, resourceType)
		parentResource = fmt.Sprintf("%v/%v", resourceType, parentResourceName)
		resourceType = "subnets"
	}

	return resourcesClient.Get(
		ctx,
		resourceGroup,
		resourceProvider,
		parentResource,
		resourceType,
		resourceName,
		APIVersion,
	)
}

// CreateAnfAccount creates an ANF Account resource
func CreateAnfAccount(ctx context.Context, location, resourceGroupName, accountName string, tags map[string]*string) (netapp.Account, error) {

	accountClient, err := getAccountsClient()
	if err != nil {
		return netapp.Account{}, err
	}

	future, err := accountClient.CreateOrUpdate(
		ctx,
		netapp.Account{
			Location: to.StringPtr(location),
			Tags:     tags,
		},
		resourceGroupName,
		accountName,
	)
	if err != nil {
		return netapp.Account{}, fmt.Errorf("cannot create account: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, accountClient.Client)
	if err != nil {
		return netapp.Account{}, fmt.Errorf("cannot get the account create or update future response: %v", err)
	}

	return future.Result(accountClient)
}

// CreateAnfCapacityPool creates an ANF Capacity Pool within ANF Account
func CreateAnfCapacityPool(ctx context.Context, location, resourceGroupName, accountName, poolName, serviceLevel string, sizeBytes int64, tags map[string]*string) (netapp.CapacityPool, error) {

	poolClient, err := getPoolsClient()
	if err != nil {
		return netapp.CapacityPool{}, err
	}

	svcLevel, err := validateAnfServiceLevel(serviceLevel)
	if err != nil {
		return netapp.CapacityPool{}, err
	}

	future, err := poolClient.CreateOrUpdate(
		ctx,
		netapp.CapacityPool{
			Location: to.StringPtr(location),
			Tags:     tags,
			PoolProperties: &netapp.PoolProperties{
				ServiceLevel: svcLevel,
				Size:         to.Int64Ptr(sizeBytes),
			},
		},
		resourceGroupName,
		accountName,
		poolName,
	)

	if err != nil {
		return netapp.CapacityPool{}, fmt.Errorf("cannot create pool: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, poolClient.Client)
	if err != nil {
		return netapp.CapacityPool{}, fmt.Errorf("cannot get the pool create or update future response: %v", err)
	}

	return future.Result(poolClient)
}

// CreateAnfVolume creates an ANF volume within a Capacity Pool
func CreateAnfVolume(ctx context.Context, location, resourceGroupName, accountName, poolName, volumeName, serviceLevel, subnetID string, protocolTypes []string, volumeUsageQuota int64, unixReadOnly, unixReadWrite bool, tags map[string]*string) (netapp.Volume, error) {

	if len(protocolTypes) > 1 {
		return netapp.Volume{}, fmt.Errorf("only one protocol type is supported at this time")
	}

	_, found := utils.FindInSlice(validProtocols, protocolTypes[0])
	if !found {
		return netapp.Volume{}, fmt.Errorf("invalid protocol type, valid protocol types are: %v", validProtocols)
	}

	svcLevel, err := validateAnfServiceLevel(serviceLevel)
	if err != nil {
		return netapp.Volume{}, err
	}

	volumeClient, err := getVolumesClient()
	if err != nil {
		return netapp.Volume{}, err
	}

	svcLevel, err = validateAnfServiceLevel(serviceLevel)
	if err != nil {
		return netapp.Volume{}, err
	}

	future, err := volumeClient.CreateOrUpdate(
		ctx,
		netapp.Volume{
			Location: to.StringPtr(location),
			Tags:     tags,
			VolumeProperties: &netapp.VolumeProperties{
				ExportPolicy: &netapp.VolumePropertiesExportPolicy{
					Rules: &[]netapp.ExportPolicyRule{
						{
							AllowedClients: to.StringPtr("0.0.0.0/0"),
							Cifs:           to.BoolPtr(map[bool]bool{true: true, false: false}[protocolTypes[0] == cifs]),
							Nfsv3:          to.BoolPtr(map[bool]bool{true: true, false: false}[protocolTypes[0] == nfsv3]),
							Nfsv41:         to.BoolPtr(map[bool]bool{true: true, false: false}[protocolTypes[0] == nfsv41]),
							RuleIndex:      to.Int32Ptr(1),
							UnixReadOnly:   to.BoolPtr(unixReadOnly),
							UnixReadWrite:  to.BoolPtr(unixReadWrite),
						},
					},
				},
				ProtocolTypes:  &protocolTypes,
				ServiceLevel:   svcLevel,
				SubnetID:       to.StringPtr(subnetID),
				UsageThreshold: to.Int64Ptr(volumeUsageQuota),
				CreationToken:  to.StringPtr(volumeName),
			},
		},
		resourceGroupName,
		accountName,
		poolName,
		volumeName,
	)

	if err != nil {
		return netapp.Volume{}, fmt.Errorf("cannot create volume: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, volumeClient.Client)
	if err != nil {
		return netapp.Volume{}, fmt.Errorf("cannot get the volume create or update future response: %v", err)
	}

	return future.Result(volumeClient)
}

func validateAnfServiceLevel(serviceLevel string) (validatedServiceLevel netapp.ServiceLevel, err error) {

	var svcLevel netapp.ServiceLevel

	switch strings.ToLower(serviceLevel) {
	case "ultra":
		svcLevel = netapp.Ultra
	case "premium":
		svcLevel = netapp.Premium
	case "standard":
		svcLevel = netapp.Standard
	default:
		return "", fmt.Errorf("invalid service level, supported service levels are: %v", netapp.PossibleServiceLevelValues())
	}

	return svcLevel, nil
}

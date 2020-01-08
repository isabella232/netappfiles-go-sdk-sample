// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

// This sample code creates an Azure Netapp Files Account, a Capacity Pool,
// and two volumes, one NFSv3 and one NFSv4.1, then it takes a snapshot
// of the first volume (NFSv3) and performs clean up if the variable
// shouldCleanUp is changed to true.
//
// This package uses go-haikunator package (https://github.com/yelinaung/go-haikunator)
// port from Python's haikunator module and therefore used here just for sample simplification,
// this doesn't mean that it is endorsed/thouroughly tested by any means, use at own risk.
// Feel free to provide your own names on variables using it.

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Azure-Samples/netappfiles-go-sdk-sample/internal/sdkutils"
	"github.com/Azure-Samples/netappfiles-go-sdk-sample/internal/utils"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/yelinaung/go-haikunator"
)

const (
	virtualNetworksApiVersion string = "2019-09-01"
)

var (
	shouldCleanUp         bool = false
	exitCode              int
	location              string = "westus2"
	resourceGroupName     string = "anf02-rg"
	vnetResourceGroupName string = "anf02-rg"
	vnetName              string = "vnet-03"
	subnetName            string = "anf-sn"
	anfAccountName        string = haikunator.New(time.Now().UTC().UnixNano()).Haikunate()
	capacityPoolName      string = "Pool01"
	serviceLevel          string = "Standard"    // Valid service levels are Standard, Premium and Ultra
	capacityPoolSizeBytes int64  = 4398046511104 // 4TiB (minimum size)
	nfsv3VolumeName       string = fmt.Sprintf("NFSv3-Vol-%v-%v", anfAccountName, capacityPoolName)
	nfsv41VolumeName      string = fmt.Sprintf("NFSv41-Vol-%v-%v", anfAccountName, capacityPoolName)
	sampleTags                   = map[string]*string{
		"Author":  to.StringPtr("ANF Go SDK Sample"),
		"Service": to.StringPtr("Azure Netapp Files"),
	}
)

func main() {

	cntx := context.Background()

	// Cleanup and exit handling
	defer func() { exit(); os.Exit(exitCode) }()

	utils.PrintHeader("Azure NetAppFiles Go SDK Sample - sample application that performs CRUD management operations (deploys NFSv3 and NFSv4.1 Volumes)")

	// Getting subscription ID from authentication file
	config, err := utils.ReadAzureBasicInfoJSON(os.Getenv("AZURE_AUTH_LOCATION"))
	if err != nil {
		utils.ConsoleOutput(fmt.Sprintf("an error ocurred getting non-sensitive info from AzureAuthFile: %v", err))
		exitCode = 1
		return
	}

	// Checking if subnet exists before any other operation starts
	subnetID := fmt.Sprintf("/subscriptions/%v/resourceGroups/%v/providers/Microsoft.Network/virtualNetworks/%v/subnets/%v",
		*config.SubscriptionID,
		vnetResourceGroupName,
		vnetName,
		subnetName,
	)

	utils.ConsoleOutput(fmt.Sprintf("Checking if subnet %v exists.", subnetID))

	_, err = sdkutils.GetResourceByID(cntx, subnetID, virtualNetworksApiVersion)
	if err != nil {
		if string(err.Error()) == "NotFound" {
			utils.ConsoleOutput(fmt.Sprintf("error: subnet %v not found: %v", subnetID, err))
		} else {
			utils.ConsoleOutput(fmt.Sprintf("error: an error ocurred trying to check if %v exists: %v", subnetID, err))
		}

		exitCode = 1
		return
	}

	// Adding

	// Azure NetApp Files Account creation
	utils.ConsoleOutput("Creating Azure NetApp Files account...")
	account, err := sdkutils.CreateAnfAccount(cntx, location, resourceGroupName, anfAccountName, sampleTags)
	if err != nil {
		utils.ConsoleOutput(fmt.Sprintf("an error ocurred while creating account: %v", err))
		exitCode = 1
		return
	}
	utils.ConsoleOutput(fmt.Sprintf("Account successfully created, resource id: %v", *account.ID))

	// Capacity pool creation
	utils.ConsoleOutput("Creating Capacity Pool...")
	capacityPool, err := sdkutils.CreateAnfCapacityPool(
		cntx,
		location,
		resourceGroupName,
		*account.Name,
		capacityPoolName,
		serviceLevel,
		capacityPoolSizeBytes,
		sampleTags,
	)
	if err != nil {
		utils.ConsoleOutput(fmt.Sprintf("an error ocurred while creating capacity pool: %v", err))
		exitCode = 1
		return
	}
	utils.ConsoleOutput(fmt.Sprintf("Capacity Pool successfully created, resource id: %v", *capacityPool.ID))

	// NFS v3 volume creation

	// Check this to set true/false on nfsv3 or nfsv41 properties of export rule
	// c := map[bool]int{true: a, false: b}[a > b]

	// NFS v4.1 volume creation

	// NFS v3 snapshot creation

}

func exit() {
	utils.ConsoleOutput("Exiting")

	if shouldCleanUp {
		utils.ConsoleOutput("\tPerforming clean up")

		// Snapshot Cleanup

		// Volumes Cleanup

		// Pool Cleanup

		// Account Cleanup

	}
}

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
	"github.com/Azure/azure-sdk-for-go/services/netapp/mgmt/2019-10-01/netapp"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/yelinaung/go-haikunator"
)

const (
	virtualNetworksApiVersion string = "2019-09-01"
)

var (
	shouldCleanUp           bool   = true
	location                string = "westus2"
	resourceGroupName       string = "anf02-rg"
	vnetResourceGroupName   string = "anf02-rg"
	vnetName                string = "vnet-03"
	subnetName              string = "anf-sn"
	anfAccountName          string = haikunator.New(time.Now().UTC().UnixNano()).Haikunate()
	capacityPoolName        string = "Pool01"
	serviceLevel            string = "Standard"          // Valid service levels are Standard, Premium and Ultra
	capacityPoolSizeBytes   int64  = 4398046511104       // 4TiB (minimum capacity pool size)
	volumeSizeBytes         int64  = 107374182400        // 100GiB (minimum volume size)
	nfsv3ProtocolTypes             = []string{"NFSv3"}   // Multiple NFS protocol types are not supported at the moment this sample was written
	nfsv41ProtocolTypes            = []string{"NFSv4.1"} // Multiple NFS protocol types are not supported at the moment this sample was written
	nfsv3VolumeName         string = fmt.Sprintf("NFSv3-Vol-%v-%v", anfAccountName, capacityPoolName)
	nfsv3SnapshotName       string = fmt.Sprintf("Snapshot-NFSv3-Vol-%v-%v", anfAccountName, capacityPoolName)
	nfsv3VolumeNameFromSnap string = fmt.Sprintf("NFSv3-FromSnapshot-Vol-%v-%v", anfAccountName, capacityPoolName)
	nfsv41VolumeName        string = fmt.Sprintf("NFSv41-Vol-%v-%v", anfAccountName, capacityPoolName)
	sampleTags                     = map[string]*string{
		"Author":  to.StringPtr("ANF Go SDK Sample"),
		"Service": to.StringPtr("Azure Netapp Files"),
	}
	exitCode                  int
	snapshotID                string = ""
	nfsv3VolumeID             string = ""
	nfsv41VolumeID            string = ""
	nfsv3VolumeFromSnapshotID string = ""
	capacityPoolID            string = ""
	acccountID                string = ""
)

func main() {

	cntx := context.Background()

	// Cleanup and exit handling
	defer func() { exit(cntx); os.Exit(exitCode) }()

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

	// Azure NetApp Files Account creation
	utils.ConsoleOutput("Creating Azure NetApp Files account...")
	account, err := sdkutils.CreateAnfAccount(cntx, location, resourceGroupName, anfAccountName, sampleTags)
	if err != nil {
		utils.ConsoleOutput(fmt.Sprintf("an error ocurred while creating account: %v", err))
		exitCode = 1
		return
	}
	acccountID = *account.ID
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
	capacityPoolID = *capacityPool.ID
	utils.ConsoleOutput(fmt.Sprintf("Capacity Pool successfully created, resource id: %v", *capacityPool.ID))

	// NFS v3 volume creation
	utils.ConsoleOutput("Creating NFSv3 Volume...")
	nfsv3Volume, err := sdkutils.CreateAnfVolume(
		cntx,
		location,
		resourceGroupName,
		*account.Name,
		capacityPoolName,
		nfsv3VolumeName,
		serviceLevel,
		subnetID,
		"",
		nfsv3ProtocolTypes,
		volumeSizeBytes,
		false,
		true,
		sampleTags,
	)
	if err != nil {
		utils.ConsoleOutput(fmt.Sprintf("an error ocurred while creating NFSv3 volume: %v", err))
		exitCode = 1
		return
	}
	nfsv3VolumeID = *nfsv3Volume.ID
	utils.ConsoleOutput(fmt.Sprintf("NFSv3 volume successfully created, resource id: %v", *nfsv3Volume.ID))

	// NFS v4.1 volume creation
	utils.ConsoleOutput("Creating NFSv4.1 Volume...")
	nfsv41Volume, err := sdkutils.CreateAnfVolume(
		cntx,
		location,
		resourceGroupName,
		*account.Name,
		capacityPoolName,
		nfsv41VolumeName,
		serviceLevel,
		subnetID,
		"",
		nfsv41ProtocolTypes,
		volumeSizeBytes,
		false,
		true,
		sampleTags,
	)
	if err != nil {
		utils.ConsoleOutput(fmt.Sprintf("an error ocurred while creating NFSv4.1 volume: %v", err))
		exitCode = 1
		return
	}
	nfsv41VolumeID = *nfsv41Volume.ID
	utils.ConsoleOutput(fmt.Sprintf("NFSv4.1 volume successfully created, resource id: %v", *nfsv41Volume.ID))

	// NFS v3 snapshot creation
	// Note: there is no difference between protocol types when creating a snapshot
	//       we're taking it from NFSv3 in this example just for convenience
	utils.ConsoleOutput("Creating Snapshot from NFSv3 Volume...")
	snapshot, err := sdkutils.CreateAnfSnapshot(
		cntx,
		location,
		resourceGroupName,
		*account.Name,
		capacityPoolName,
		nfsv3VolumeName,
		nfsv3SnapshotName,
		sampleTags,
	)
	if err != nil {
		utils.ConsoleOutput(fmt.Sprintf("an error ocurred while creating snapshot from NFSv3 volume: %v", err))
		exitCode = 1
		return
	}
	snapshotID = *snapshot.ID
	utils.ConsoleOutput(fmt.Sprintf("Snapshot successfully created, resource id: %v", *snapshot.ID))

	// Creating new volume (NFSv3) from Snapshot
	// Note: At the time when this sample code was written, creating a volume from snapshot with a different protocol
	//       other than the protocol from the source volume is not supported.
	utils.ConsoleOutput("Creating new NFSv3 Volume from Snapshot...")
	newNFSv3Volume, err := sdkutils.CreateAnfVolume(
		cntx,
		location,
		resourceGroupName,
		*account.Name,
		capacityPoolName,
		nfsv3VolumeNameFromSnap,
		serviceLevel,
		subnetID,
		*snapshot.SnapshotID,
		nfsv3ProtocolTypes,
		volumeSizeBytes,
		false,
		true,
		sampleTags,
	)
	if err != nil {
		utils.ConsoleOutput(fmt.Sprintf("an error ocurred while creating NFSv3 volume from snapshot: %v", err))
		exitCode = 1
		return
	}
	nfsv3VolumeFromSnapshotID = *newNFSv3Volume.ID
	utils.ConsoleOutput(fmt.Sprintf("NFSv3 volume from snapshot successfully created, resource id: %v", *newNFSv3Volume.ID))

	// Update NFS v4 volume size to double its size (200GiB in this example)
	utils.ConsoleOutput("Updating NFSv4.1 volume size...")

	newVolumeSize := volumeSizeBytes * int64(2)
	volumeChanges := netapp.VolumePatchProperties{
		UsageThreshold: &newVolumeSize,
	}

	updatedNFS41Volume, err := sdkutils.UpdateAnfVolume(
		cntx,
		location,
		resourceGroupName,
		*account.Name,
		capacityPoolName,
		nfsv41VolumeName,
		volumeChanges,
		sampleTags,
	)
	if err != nil {
		utils.ConsoleOutput(fmt.Sprintf("an error ocurred while updating NFSv4.1 volume: %v", err))
		exitCode = 1
		return
	}
	utils.ConsoleOutput(fmt.Sprintf("NFSv4.1 volume successfully update with new size %v, resource id: %v", newVolumeSize, *updatedNFS41Volume.ID))

}

func exit(cntx context.Context) {
	utils.ConsoleOutput("Exiting")

	if shouldCleanUp {
		utils.ConsoleOutput("\tPerforming clean up")

		// Volume restored from Snaphost cleanup
		utils.ConsoleOutput("\tCleaning up NFSv3 Volume Restored from Snapshot ...")
		time.Sleep(3 * time.Second)
		err := sdkutils.DeleteAnfVolume(
			cntx,
			resourceGroupName,
			anfAccountName,
			capacityPoolName,
			nfsv3VolumeNameFromSnap,
		)
		if err != nil {
			utils.ConsoleOutput(fmt.Sprintf("an error ocurred while deleting volume: %v", err))
			exitCode = 1
			return
		}
		sdkutils.WaitForNoANFResource(cntx, nfsv3VolumeFromSnapshotID, 60, 60)
		utils.ConsoleOutput("\tVolume successfully deleted")

		// Snapshot Cleanup
		utils.ConsoleOutput("\tCleaning up NFSv3 Volume Snapshot ...")
		err = sdkutils.DeleteAnfSnapshot(
			cntx,
			resourceGroupName,
			anfAccountName,
			capacityPoolName,
			nfsv3VolumeName,
			nfsv3SnapshotName,
		)
		if err != nil {
			utils.ConsoleOutput(fmt.Sprintf("an error ocurred while deleting NFSv3 volume snapshot: %v", err))
			exitCode = 1
			return
		}
		sdkutils.WaitForNoANFResource(cntx, snapshotID, 60, 60)
		utils.ConsoleOutput("\tSnapshot successfully deleted")

		// Other Volumes Cleanup
		utils.ConsoleOutput("\tCleaning up other volumes...")
		volumes := map[string]string{
			nfsv3VolumeName:  nfsv3VolumeID,
			nfsv41VolumeName: nfsv41VolumeID,
		}
		for volumeName, resourceID := range volumes {
			utils.ConsoleOutput(fmt.Sprintf("\tCleaning up volume %v", volumeName))
			err := sdkutils.DeleteAnfVolume(
				cntx,
				resourceGroupName,
				anfAccountName,
				capacityPoolName,
				volumeName,
			)
			if err != nil {
				utils.ConsoleOutput(fmt.Sprintf("an error ocurred while deleting volume: %v", err))
				exitCode = 1
				return
			}
			sdkutils.WaitForNoANFResource(cntx, resourceID, 60, 60)
			utils.ConsoleOutput("\tVolume successfully deleted")
		}

		// Pool Cleanup
		utils.ConsoleOutput("\tCleaning up capacity pool...")
		err = sdkutils.DeleteAnfCapacityPool(
			cntx,
			resourceGroupName,
			anfAccountName,
			capacityPoolName,
		)
		if err != nil {
			utils.ConsoleOutput(fmt.Sprintf("an error ocurred while deleting capacity pool: %v", err))
			exitCode = 1
			return
		}
		sdkutils.WaitForNoANFResource(cntx, capacityPoolID, 60, 60)
		utils.ConsoleOutput("\tCapacity pool successfully deleted")

		// Account Cleanup
		utils.ConsoleOutput("\tCleaning up account...")
		err = sdkutils.DeleteAnfAccount(
			cntx,
			resourceGroupName,
			anfAccountName,
		)
		if err != nil {
			utils.ConsoleOutput(fmt.Sprintf("an error ocurred while deleting account: %v", err))
			exitCode = 1
			return
		}
		sdkutils.WaitForNoANFResource(cntx, acccountID, 60, 60)
		utils.ConsoleOutput("\tAccount successfully deleted")
		utils.ConsoleOutput("\tCleanup completed!")
	}
}

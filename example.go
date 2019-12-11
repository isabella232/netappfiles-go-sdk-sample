// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package main

import (
	//"context"
	//"fmt"
	//"os"

	"fmt"

	"github.com/Azure-Samples/netappfiles-go-sdk-sample/internal/iam"
	"github.com/Azure-Samples/netappfiles-go-sdk-sample/internal/utils"
	//"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-05-01/resources"
	//"github.com/Azure/go-autorest/autorest/azure/auth"
)

func main() {

	utils.PrintHeader("Azure NetAppFiles Go SDK Sample - sample application that performs CRUD management operations (deploys NFSv3 and NFSv4.1 Volumes)")

	authorizer, err := iam.GetAuthorizer()
	if err != nil {
		utils.ConsoleOutput(fmt.Sprintf("an error ocurred getting authorizer token: %v\n", err))
		return
	}

	test, err := utils.GetResourceValue("aa", "aa")
	if err != nil {
		utils.ConsoleOutput(fmt.Sprintf("an error ocurred getting resource value: %v\n", err))
		return
	}

	fmt.Println(test)

	fmt.Println(*authorizer)

	// cntx := context.Background()

	// fmt.Println(os.Getenv("AZURE_AUTH_LOCATION"))
	// subscriptionID := "66bc9830-19b6-4987-94d2-0e487be7aa47"
	// resourceManagerEndpointURL := "https://management.azure.com/"

	// resourceClient := resources.NewGroupsClient(subscriptionID)

	// authorizer, err := auth.NewAuthorizerFromFile(resourceManagerEndpointURL)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// resourceClient.Authorizer = authorizer
	// resourceClient.AddToUserAgent("sdk-sample")

	// resourceGroups, err := resources.GroupsClient.List(resourceClient, cntx, "", nil)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// for _, rg := range resourceGroups.Values() {

	// 	fmt.Println(*rg.Name)
	// }

	// fmt.Println(resourceGroups)
	// for _, rg := range resourceGroups.NextWithContext(cntx) {
	// 	fmt.Println(rg.name)
	// }

}

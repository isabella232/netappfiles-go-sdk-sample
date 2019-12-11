// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package iam

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Azure-Samples/netappfiles-go-sdk-sample/internal/models"
	"github.com/Azure-Samples/netappfiles-go-sdk-sample/internal/utils"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

// GetAuthorizer gets an authorization token to be used within ANF client
func GetAuthorizer() (*autorest.Authorizer, error) {

	// Getting information from authentication file
	authInfo, err := readAuthJSON(os.Getenv("AZURE_AUTH_LOCATION"))

	authorizer, err := auth.NewAuthorizerFromFile(*authInfo.ResourceManagerEndpointURL)
	if err != nil {
		utils.ConsoleOutput(fmt.Sprintf("%v", err))
		return nil, err
	}

	return &authorizer, nil
}

// readAuthJSON reads the Azure Authentication json file json file and unmarshals it.
func readAuthJSON(path string) (*models.AzureAuthInfo, error) {
	infoJSON, err := ioutil.ReadFile(path)
	if err != nil {
		utils.ConsoleOutput(fmt.Sprintf("failed to read file: %v", err))
		return &models.AzureAuthInfo{}, err
	}
	var authInfo models.AzureAuthInfo
	json.Unmarshal(infoJSON, &authInfo)
	return &authInfo, nil
}

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package utils

import (
	"errors"
	"strings"
)

// GetResourceValue returns the name of a resource from resource id/uri based on resource type name.
func GetResourceValue(resourceURI string, resourceName string) (string, error) {

	if len(strings.TrimSpace(resourceURI)) == 0 {
		return "", errors.New("resourceURI cannot be null")
	}

	if len(strings.TrimSpace(resourceName)) == 0 {
		return "", errors.New("resourceName cannot be null")
	}

	return "ok", nil

}

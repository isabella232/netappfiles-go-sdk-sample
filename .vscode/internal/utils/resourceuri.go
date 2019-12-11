// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package resourceuri

import (
	"strings"
)

// GetResourceValue returns the name of a resource from resource id/uri based on resource type name.
func GetResourceValue(resourceURI string, resourceName string) string, error {

	if Len(TrimSpace(resourceUri)) == 0 || Len(TrimSpace(resourceName)) == 0 {
		return nil, error.New("resourceURI and resourceName cannot be null")
	}

}

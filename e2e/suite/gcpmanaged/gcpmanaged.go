/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
limitations under the License.
*/
package gcpmanaged

import (
	"os"

	// nolint
	. "github.com/onsi/ginkgo"
	// nolint
	. "github.com/onsi/ginkgo/extensions/table"

	// nolint
	// . "github.com/onsi/gomega"
	"github.com/external-secrets/external-secrets/e2e/framework"
	"github.com/external-secrets/external-secrets/e2e/suite/common"
	"github.com/external-secrets/external-secrets/e2e/suite/gcp"
)

var _ = Describe("[gcpmanaged] ", func() {
	if os.Getenv("FOCUS") == "gcpmanaged" {
		f := framework.New("eso-gcp-managed")
		projectID := os.Getenv("GCP_PROJECT_ID")
		prov := &gcp.GcpProvider{}

		if projectID != "" {
			prov = gcp.NewgcpProvider(f, "", projectID)
		}

		DescribeTable("sync secrets", framework.TableFunc(f, prov),
			Entry(common.SimpleDataSync(f)),
			Entry(common.JSONDataWithProperty(f)),
			Entry(common.JSONDataFromSync(f)),
			Entry(common.NestedJSONWithGJSON(f)),
			Entry(common.JSONDataWithTemplate(f)),
			Entry(common.DockerJSONConfig(f)),
			Entry(common.DataPropertyDockerconfigJSON(f)),
			Entry(common.SSHKeySync(f)),
			Entry(common.SSHKeySyncDataProperty(f)),
			Entry(common.SyncWithoutTargetName(f)),
			Entry(common.JSONDataWithoutTargetName(f)),
		)
	}
})

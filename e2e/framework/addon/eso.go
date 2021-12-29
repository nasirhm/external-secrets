/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package addon

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	// nolint
	. "github.com/onsi/ginkgo"
)

type ESO struct {
	Addon
}

func NewESO() *ESO {
	return &ESO{
		&HelmChart{
			Namespace:   "default",
			ReleaseName: "eso",
			Chart:       "/k8s/deploy/charts/external-secrets",
			Values:      []string{"/k8s/eso.values.yaml"},
		},
	}
}

func (l *ESO) Install() error {
	By("Installing eso\n")
	err := l.Addon.Install()
	if err != nil {
		return err
	}

	By("afterInstall eso\n")
	err = l.afterInstall()
	if err != nil {
		return err
	}

	return nil
}

func (l *ESO) afterInstall() error {
	gcpProjectID := os.Getenv("GCP_PROJECT_ID")
	gcpGSAName := os.Getenv("GCP_GSA_NAME")
	gcpKSAName := os.Getenv("GCP_KSA_NAME")

	var sout, serr bytes.Buffer

	cmd := exec.Command(
		"kubectl", "annotate", "--overwrite", "serviceaccount", gcpKSAName, "--namespace",
		"default", fmt.Sprintf("iam.gke.io/gcp-service-account=%s@%s.iam.gserviceaccount.com", gcpGSAName, gcpProjectID),
	)
	cmd.Stdout = &sout
	cmd.Stderr = &serr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("unable to annotate SA: %w: %s, %s", err, sout.String(), serr.String())
	}

	cmd = exec.Command( //nolint:gosec
		"kubectl", "annotate", "--overwrite", "serviceaccount", "external-secrets-e2e", "--namespace",
		"default", fmt.Sprintf("iam.gke.io/gcp-service-account=%s@%s.iam.gserviceaccount.com", gcpGSAName, gcpProjectID),
	)
	cmd.Stdout = &sout
	cmd.Stderr = &serr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("unable to annotate SA: %w: %s, %s", err, sout.String(), serr.String())
	}
	return nil
}

func NewScopedESO() *ESO {
	return &ESO{
		&HelmChart{
			Namespace:   "default",
			ReleaseName: "eso-aws-sm",
			Chart:       "/k8s/deploy/charts/external-secrets",
			Values:      []string{"/k8s/eso.scoped.values.yaml"},
		},
	}
}

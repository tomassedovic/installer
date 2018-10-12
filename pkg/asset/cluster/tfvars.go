package cluster

import (
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/tfvars"
	"github.com/pkg/errors"
)

const (
	tfvarsFilename  = "terraform.tfvars"
	tfvarsAssetName = "Terraform Variables"
)

// TerraformVariables depends on InstallConfig and
// Ignition to generate the terrafor.tfvars.
type TerraformVariables struct {
	Platform string
	File     *asset.File
}

var _ asset.WritableAsset = (*TerraformVariables)(nil)

// Name returns the human-friendly name of the asset.
func (t *TerraformVariables) Name() string {
	return tfvarsAssetName
}

// Dependencies returns the dependency of the TerraformVariable
func (t *TerraformVariables) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
		&bootstrap.Bootstrap{},
		&machine.Master{},
		&machine.Worker{},
	}
}

// Generate generates the terraform.tfvars file.
func (t *TerraformVariables) Generate(parents asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	bootstrap := &bootstrap.Bootstrap{}
	master := &machine.Master{}
	worker := &machine.Worker{}
	parents.Get(installConfig, bootstrap, master, worker)

	t.Platform = installConfig.Config.Platform.Name()

	bootstrapIgn := string(bootstrap.Files()[0].Data)

	masterFiles := master.Files()
	masterIgns := make([]string, len(masterFiles))
	for i, f := range masterFiles {
		masterIgns[i] = string(f.Data)
	}

	workerIgn := string(worker.Files()[0].Data)

	data, err := tfvars.TFVars(installConfig.Config, bootstrapIgn, masterIgns, workerIgn)
	if err != nil {
		return errors.Wrap(err, "failed to get Tfvars")
	}
	t.File = &asset.File{
		Filename: tfvarsFilename,
		Data:     data,
	}

	return nil
}

// Files returns the files generated by the asset.
func (t *TerraformVariables) Files() []*asset.File {
	if t.File != nil {
		return []*asset.File{t.File}
	}
	return []*asset.File{}
}

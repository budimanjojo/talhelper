package config

import (
	"github.com/siderolabs/talos/pkg/machinery/config/types/v1alpha1"
)

func (mfs MachineFiles) GetMFs() []*v1alpha1.MachineFile {
	result := make([]*v1alpha1.MachineFile, 0, len(mfs))

	for _, mf := range mfs {
		result = append(result, &mf.MachineFile)
	}
	return result
}

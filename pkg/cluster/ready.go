package cluster

import (
	"context"
	"sort"
	"strings"

	"github.com/openshift/openshift-azure/pkg/api"
	"github.com/openshift/openshift-azure/pkg/cluster/names"
)

func (u *Upgrade) WaitForNodesInAgentPoolProfile(ctx context.Context, app *api.AgentPoolProfile, suffix string) error {
	vms, err := u.Vmc.List(ctx, u.Cs.Properties.AzProfile.ResourceGroup, names.GetScalesetName(app, suffix), "", "", "")
	if err != nil {
		return err
	}
	for _, vm := range vms {
		hostname := strings.ToLower(*vm.VirtualMachineScaleSetVMProperties.OsProfile.ComputerName)
		u.Log.Infof("waiting for %s to be ready", hostname)
		if app.Role == api.AgentPoolProfileRoleMaster {
			err = u.Interface.WaitForReadyMaster(ctx, hostname)
		} else {
			err = u.Interface.WaitForReadyWorker(ctx, hostname)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// SortedAgentPoolProfilesForRole returns a shallow copy of the
// AgentPoolProfiles of a given role, sorted by name
func (u *Upgrade) SortedAgentPoolProfilesForRole(role api.AgentPoolProfileRole) (apps []api.AgentPoolProfile) {
	for _, app := range u.Cs.Properties.AgentPoolProfiles {
		if app.Role == role {
			apps = append(apps, app)
		}
	}
	sort.Slice(apps, func(i, j int) bool { return apps[i].Name < apps[j].Name })
	return apps
}

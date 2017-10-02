package hostwatch

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/rancher/go-rancher-metadata/metadata"
	"github.com/rancher/kubernetes-agent/kubernetesclient"
)

type hostStatusSyncer struct {
	kClient        *kubernetesclient.Client
	metadataClient metadata.Client
}

var (
	maxRetryCount int = 3
)

func (h *hostStatusSyncer) syncHoststatus(version string) {
	err := statusSync(h.kClient, h.metadataClient)
	if err != nil {
		log.Errorf("Error syncing host status: [%v]", err)
	}
}

func statusSync(kClient *kubernetesclient.Client, metadataClient metadata.Client) error {
	hosts, err := metadataClient.GetHosts()
	if err != nil {
		return fmt.Errorf("Error reading host list from metadata service: [%v]", err)
	}
	for _, host := range hosts {
		mNodeStatus := host.State
		switch mNodeStatus {
		case ActiveState:
			err = cordonUncordon(host, kClient, metadataClient, false)
		case InActiveState:
			err = cordonUncordon(host, kClient, metadataClient, true)
		case EvictedState:
			//todo
		}
	}
	return err
}

func cordonUncordon(host metadata.Host, kClient *kubernetesclient.Client, metadataClient metadata.Client, unschedulable bool) error {
	changed := false
	for retryCount := 0; retryCount <= maxRetryCount; retryCount++ {
		node, err := kClient.Node.ByName(host.Hostname)
		if err != nil {
			log.Errorf("Error getting node: [%s] by name from kubernetes, err: [%v]", host.Hostname, err)
			continue
		}
		node.Spec.Unschedulable = unschedulable
		_, err = kClient.Node.ReplaceNode(node)
		if err != nil {
			log.Errorf("Error updating node [%s] with new schedulable state, err :[%v]", host.Hostname, err)
			continue
		}
		changed = true
		break
	}
	if !changed {
		return fmt.Errorf("Failed to cordon/uncordon node: [%s]", host.Hostname)
	}
	return nil
}

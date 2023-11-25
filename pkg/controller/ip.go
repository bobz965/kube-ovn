package controller

import (
	"fmt"
	"strings"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"

	kubeovnv1 "github.com/kubeovn/kube-ovn/pkg/apis/kubeovn/v1"
	"github.com/kubeovn/kube-ovn/pkg/util"
)

func (c *Controller) enqueueAddOrDelIP(obj interface{}) {
	if _, ok := obj.(*kubeovnv1.IP); !ok {
		klog.Errorf("object is not an IP, ignore it")
		return
	}

	ipObj := obj.(*kubeovnv1.IP)
	klog.V(3).Infof("enqueue update status subnet %s", ipObj.Spec.Subnet)
	if strings.HasPrefix(ipObj.Name, util.U2OInterconnName[0:19]) {
		return
	}
	c.updateSubnetStatusQueue.Add(ipObj.Spec.Subnet)
	for _, as := range ipObj.Spec.AttachSubnets {
		klog.V(3).Infof("enqueue update status subnet %s", as)
		c.updateSubnetStatusQueue.Add(as)
	}
}

func (c *Controller) enqueueUpdateIP(oldObj, newObj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(newObj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	oldIP := oldObj.(*kubeovnv1.IP)
	newIP := newObj.(*kubeovnv1.IP)
	if newIP.Spec.PodType == util.VM {
		if oldIP.Annotations[util.CNIDelVMPod] != newIP.Annotations[util.CNIDelVMPod] {
			klog.Infof("enqueue update ip %s for cni del vm pod", key)
			c.updateIPQueue.Add(key)
		}
	}
	klog.V(3).Infof("enqueue update status subnet %s", newIP.Spec.Subnet)
	for _, as := range newIP.Spec.AttachSubnets {
		klog.V(3).Infof("enqueue update status subnet %s", as)
		c.updateSubnetStatusQueue.Add(as)
	}
}

func (c *Controller) runUpdateIPWorker() {
	for c.processNextUpdateIPWorkItem() {
	}
}

func (c *Controller) processNextUpdateIPWorkItem() bool {
	obj, shutdown := c.updateIPQueue.Get()
	if shutdown {
		return false
	}
	err := func(obj interface{}) error {
		defer c.updateIPQueue.Done(obj)
		var key string
		var ok bool
		if key, ok = obj.(string); !ok {
			c.updateIPQueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		if err := c.handleUpdateIP(key); err != nil {
			c.updateIPQueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}
		c.updateIPQueue.Forget(obj)
		return nil
	}(obj)
	if err != nil {
		utilruntime.HandleError(err)
		return true
	}
	return true
}

func (c *Controller) handleUpdateIP(key string) error {
	cachedIP, err := c.ipsLister.Get(key)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		klog.Error(err)
		return err
	}
	needUpdate := false
	delPod, ok := cachedIP.Annotations[util.CNIDelVMPod]
	if ok {
		if delPod != "" {
			needUpdate = true
		} else {
			return nil
		}
	} else {
		needUpdate = true
	}
	if !needUpdate {
		return nil
	}

	klog.Infof("clean migrate options for ip %s", cachedIP.Name)
	if err := c.cleanVMLSPMigration(cachedIP.Name); err != nil {
		err := fmt.Errorf("failed to clean migrate options for ip %s, %v", cachedIP.Name, err)
		klog.Error(err)
		return err
	}
	return nil
}

func (c *Controller) cleanVMLSPMigration(portName string) error {
	if err := c.OVNNbClient.CleanLogicalSwitchPortMigrateOptions(portName); err != nil {
		err = fmt.Errorf("failed to clean migrate options for port %s, %v", portName, err)
		klog.Error(err)
		return err
	}
	return nil
}

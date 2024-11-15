package ovs

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func (suite *OvnClientTestSuite) testSetInterfaceBandwidth() {
	t := suite.T()
	t.Parallel()

	if !suite.enableOvsSandbox {
		return
	}

	err := SetInterfaceBandwidth("podName", "podNS", "eth0", "10", "10")
	require.NoError(t, err)
	return
}

func (suite *OvnClientTestSuite) testClearHtbQosQueue() {
	t := suite.T()
	t.Parallel()

	if !suite.enableOvsSandbox {
		return
	}

	err := ClearHtbQosQueue("podName", "podNS", "eth0")
	require.NoError(t, err)
	return

}

func (suite *OvnClientTestSuite) testIsHtbQos() {
	t := suite.T()
	t.Parallel()

	if !suite.enableOvsSandbox {
		return
	}

	isHtbQos, err := IsHtbQos("eth0")
	require.NoError(t, err)
	require.False(t, isHtbQos)
}

func (suite *OvnClientTestSuite) testSetHtbQosQueueRecord() {
	t := suite.T()
	t.Parallel()

	if !suite.enableOvsSandbox {
		return
	}

	// get a new id
	id, err := SetHtbQosQueueRecord("podName", "podNS", "eth0", 10, nil)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	// get a exist id
	queueIfaceUIDMap := make(map[string]string)
	queueIfaceUIDMap["eth0"] = "123"
	id, err = SetHtbQosQueueRecord("podName", "podNS", "eth0", 10, queueIfaceUIDMap)
	require.NoError(t, err)
	require.NotEmpty(t, id)
}

func (suite *OvnClientTestSuite) testSetQosQueueBinding() {
	t := suite.T()
	t.Parallel()

	if !suite.enableOvsSandbox {
		return
	}

	// get with invalid id
	err := SetQosQueueBinding("podName", "podNS", "podName.podNS", "eth0", "123", nil)
	require.Error(t, err)
	uid := uuid.New().String()
	err = SetQosQueueBinding("podName", "podNS", "podName.podNS", "eth0", uid, nil)
	require.Error(t, err)
	// get a exist id
	queueIfaceUIDMap := make(map[string]string)
	queueIfaceUIDMap["eth0"] = "123"
	err = SetQosQueueBinding("podName", "podNS", "podName.podNS", "eth0", "123", queueIfaceUIDMap)
	require.Error(t, err)
}

func (suite *OvnClientTestSuite) testSetNetemQos() {
	t := suite.T()
	t.Parallel()

	if !suite.enableOvsSandbox {
		return
	}

	err := SetNetemQos("podName", "podNS", "eth0", "10", "10", "10", "10")
	if suite.enableOvsSandbox {
		require.NoError(t, err)
		return
	}
}

func (suite *OvnClientTestSuite) testGetNetemQosConfig() {
	t := suite.T()
	t.Parallel()

	if !suite.enableOvsSandbox {
		return
	}

	latencyMs := 1
	latencyUs := latencyMs * 1000
	jitterMs := 1
	jitterUs := jitterMs * 1000
	limitPkts := 1
	lossPercent := 1
	iface := "test-get-netem-qos"
	err := ovsAdd("interface", iface, "type", "internal")
	require.NoError(t, err)
	err = ovsSet("interface", iface, fmt.Sprintf("external-ids:iface-id=%s", iface))
	require.NoError(t, err)
	interfaceList, err := ovsFind("interface", "name", fmt.Sprintf("external-ids:iface-id=%s", iface))
	require.NoError(t, err)
	require.NotEmpty(t, interfaceList)
	qosList, err := GetQosList("", "", iface)
	require.NoError(t, err)
	require.Empty(t, qosList)

	var qosCommandValues []string
	qosCommandValues = append(qosCommandValues, fmt.Sprintf("other_config:latency=%d", latencyUs))
	qosCommandValues = append(qosCommandValues, fmt.Sprintf("other_config:jitter=%d", jitterUs))
	qosCommandValues = append(qosCommandValues, fmt.Sprintf("other_config:limit=%d", limitPkts))
	qosCommandValues = append(qosCommandValues, fmt.Sprintf("other_config:loss=%v", lossPercent))
	qosCommandValues = append(qosCommandValues, "type=linux-netem", fmt.Sprintf(`external-ids:iface-id="%s"`, iface))
	podNamespace := "test-namespace"
	podName := "test-pod"
	if podNamespace != "" && podName != "" {
		qosCommandValues = append(qosCommandValues, fmt.Sprintf("external-ids:pod=%s/%s", podNamespace, podName))
	}
	qos, err := ovsCreate("qos", qosCommandValues...)
	require.NoError(t, err)
	err = ovsSet("port", iface, fmt.Sprintf("qos=%s", qos))
	require.NoError(t, err)
	latency, loss, limit, jitter, err := getNetemQosConfig("name")
	require.NoError(t, err)
	require.NotEmpty(t, latency)
	require.NotEmpty(t, loss)
	require.NotEmpty(t, limit)
	require.NotEmpty(t, jitter)
}

func (suite *OvnClientTestSuite) testDeleteNetemQosByID() {
	t := suite.T()
	t.Parallel()

	if !suite.enableOvsSandbox {
		return
	}

	err := deleteNetemQosByID("qosID", "eth0", "podName", "podNS")
	require.NoError(t, err)
}

func (suite *OvnClientTestSuite) testIsUserspaceDataPath() {
	t := suite.T()
	t.Parallel()
	if !suite.enableOvsSandbox {
		return
	}

	ret, err := Exec("set", "Bridge", "br-int", "datapath_type=netdev")
	require.NoError(t, err)
	require.Empty(t, ret)

	isUserspace, err := IsUserspaceDataPath()
	require.NoError(t, err)
	require.True(t, isUserspace)
}

func (suite *OvnClientTestSuite) testCheckAndUpdateHtbQos() {
	t := suite.T()
	t.Parallel()

	if !suite.enableOvsSandbox {
		return
	}

	// get a new id
	err := CheckAndUpdateHtbQos("podName", "podNS", "eth0", nil)
	require.NoError(t, err)
	// get a exist id
	queueIfaceUIDMap := make(map[string]string)
	queueIfaceUIDMap["eth0"] = "name"
	err = CheckAndUpdateHtbQos("podName", "podNS", "eth0", queueIfaceUIDMap)
	require.NoError(t, err)
}

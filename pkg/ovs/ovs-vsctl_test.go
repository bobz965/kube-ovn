package ovs

import (
	"fmt"

	"github.com/stretchr/testify/require"
)

func (suite *OvnClientTestSuite) testUpdateOVSVsctlLimiter() {
	t := suite.T()
	t.Parallel()

	UpdateOVSVsctlLimiter(int32(10))
}

func (suite *OvnClientTestSuite) testOvsExec() {
	t := suite.T()
	t.Parallel()

	if !suite.enableOvsSandbox {
		return
	}

	ret, err := Exec("add-port", "br-int", "test-ovs-exec")
	require.NoError(t, err)
	require.Empty(t, ret)
}

func (suite *OvnClientTestSuite) testOvsCreate() {
	t := suite.T()
	t.Parallel()

	if !suite.enableOvsSandbox {
		return
	}

	var qosCommandValues []string
	qosCommandValues = append(qosCommandValues, fmt.Sprintf("other_config:latency=%d", 10))
	qosCommandValues = append(qosCommandValues, fmt.Sprintf("other_config:jitter=%d", 10))
	qosCommandValues = append(qosCommandValues, fmt.Sprintf("other_config:limit=%d", 10))
	qosCommandValues = append(qosCommandValues, fmt.Sprintf("other_config:loss=%v", 10))
	ret, err := ovsCreate("qos", qosCommandValues...)
	require.NoError(t, err)
	require.NotEmpty(t, ret)
}

func (suite *OvnClientTestSuite) testOvsDestroy() {
	t := suite.T()
	t.Parallel()

	if !suite.enableOvsSandbox {
		return
	}

	err := ovsDestroy("qos", "qos-uuid")
	require.NoError(t, err)
}

func (suite *OvnClientTestSuite) testOvsSet() {
	t := suite.T()
	t.Parallel()

	if !suite.enableOvsSandbox {
		return
	}

	ret, err := Exec("add-port", "br-int", "test-ovs-set")
	require.NoError(t, err)
	require.Empty(t, ret)
	err = ovsAdd("port", "test-ovs-set", "qos=qos-uuid")
	require.Error(t, err)
	err = ovsSet("port", "ovs-set-port-test", "qos=qos-uuid1")
	require.Error(t, err)
}

func (suite *OvnClientTestSuite) testOvsAdd() {
	t := suite.T()
	t.Parallel()

	if !suite.enableOvsSandbox {
		return
	}

	err := ovsAdd("port", "ovs-add", "tag", "qos-uuid")
	require.NoError(t, err)
}

func (suite *OvnClientTestSuite) testOvsFind() {
	t := suite.T()
	t.Parallel()

	if !suite.enableOvsSandbox {
		return
	}
	// port not exist
	_, err := ovsFind("mirror", "m0", "select_all")
	require.Error(t, err)
	// port exist
	err = ovsSet("bridge", "br-int", "mirrors=@m")
	require.NoError(t, err)
	ret, err := ovsFind("mirror", "m0", "select_all")
	require.NoError(t, err)
	require.NotEmpty(t, ret)
}

func (suite *OvnClientTestSuite) testParseOvsFindOutput() {
	t := suite.T()
	t.Parallel()
	input := `br-int

br-businessnet
`
	ret := parseOvsFindOutput(input)
	require.Len(t, ret, 2)
}

func (suite *OvnClientTestSuite) testOvsClear() {
	t := suite.T()
	t.Parallel()

	if !suite.enableOvsSandbox {
		return
	}

	ret, err := Exec("add-port", "br-int", "test-ovs-clear")
	require.NoError(t, err)
	require.Empty(t, ret)
	err = ovsAdd("port", "test-ovs-clear", "qos=qos-uuid")
	require.NoError(t, err)
	err = ovsClear("port", "test-ovs-clear", "qos")
	require.NoError(t, err)
}

func (suite *OvnClientTestSuite) testOvsGet() {
	t := suite.T()
	t.Parallel()

	if !suite.enableOvsSandbox {
		return
	}

	// port not exist
	ret, err := ovsGet("port", "br-int", "name", "test-ovs-get")
	require.Error(t, err)
	require.Empty(t, ret)
	// port exist
	ret, err = Exec("add-port", "br-int", "test-ovs-get")
	require.NoError(t, err)
	require.Empty(t, ret)
	err = ovsAdd("port", "test-ovs-get", "tag", "10")
	require.NoError(t, err)
	ret, err = ovsGet("port", "br-int", "name", "test-ovs-get")
	require.Error(t, err)
	require.NotEmpty(t, ret)
}

func (suite *OvnClientTestSuite) testOvsFindBridges() {
	t := suite.T()
	t.Parallel()

	if !suite.enableOvsSandbox {
		return
	}

	bridges, err := Bridges()
	require.NoError(t, err)
	require.NotEmpty(t, bridges)
}

func (suite *OvnClientTestSuite) testOvsBridgeExists() {
	t := suite.T()
	t.Parallel()

	if !suite.enableOvsSandbox {
		return
	}

	// not exist
	ret, err := BridgeExists("not-exist-bridge")
	require.NoError(t, err)
	require.False(t, ret)
	// exist
	ret, err = BridgeExists("br-int")
	require.NoError(t, err)
	require.True(t, ret)
}

func (suite *OvnClientTestSuite) testOvsPortExists() {
	t := suite.T()
	t.Parallel()

	if !suite.enableOvsSandbox {
		return
	}

	// not exist
	ret, err := PortExists("not-exist-port")
	require.NoError(t, err)
	require.False(t, ret)

	// exist
	ret1, err := Exec("add-port", "br-int", "test-ovs-set")
	require.NoError(t, err)
	require.Empty(t, ret1)
	ret, err = PortExists("not-exist-port")
	require.NoError(t, err)
	require.False(t, ret)

}

func (suite *OvnClientTestSuite) testGetOvsQosList() {
	t := suite.T()
	t.Parallel()

	if !suite.enableOvsSandbox {
		return
	}

	// 1. qos not exist
	ret, err := GetQosList("pod-name", "pod-namespace", "iface-id")
	require.NoError(t, err)
	require.Empty(t, ret)

	ret, err = GetQosList("pod-name", "pod-namespace", "")
	require.NoError(t, err)
	require.Empty(t, ret)

	// 2. qos exist
	ret1, err := Exec("add-port", "br-int", "ovs-get-qos-list")
	require.NoError(t, err)
	require.Empty(t, ret1)
	err = ovsAdd("port", "ovs-get-qos-list", "qos=qos-uuid")
	require.NoError(t, err)
	ret, err = GetQosList("pod-name", "pod-namespace", "ovs-get-qos-list")
	require.NoError(t, err)
	require.Empty(t, ret)
}

func (suite *OvnClientTestSuite) testOvsClearPodBandwidth() {
	t := suite.T()
	t.Parallel()

	if !suite.enableOvsSandbox {
		return
	}

	err := ClearPodBandwidth("pod-name", "pod-namespace", "iface-id")
	require.NoError(t, err)
}

func (suite *OvnClientTestSuite) testOvsCleanLostInterface() {
	t := suite.T()
	t.Parallel()

	if !suite.enableOvsSandbox {
		return
	}

	CleanLostInterface()
}

func (suite *OvnClientTestSuite) testOvsCleanDuplicatePort() {
	t := suite.T()
	t.Parallel()

	if !suite.enableOvsSandbox {
		return
	}

	CleanDuplicatePort("iface-id", "port-name")
}

func (suite *OvnClientTestSuite) testOvsSetPortTag() {
	t := suite.T()
	t.Parallel()

	if !suite.enableOvsSandbox {
		return
	}

	// 1. port not exist
	err := SetPortTag("ovs-set-port-tag", "tag")
	require.Error(t, err)

	// 2. create port
	ret, err := Exec("add-port", "br-int", "ovs-set-port-tag")
	require.NoError(t, err)
	require.Empty(t, ret)
	err = ovsAdd("port", "ovs-set-port-tag", "tag", "ovs-set-port-tag")
	require.Error(t, err)
	err = SetPortTag("ovs-set-port-tag", "tag")
	require.Error(t, err)
}

func (suite *OvnClientTestSuite) testValidatePortVendor() {
	t := suite.T()
	t.Parallel()

	if !suite.enableOvsSandbox {
		return
	}

	ret, err := Exec("add-port", "br-int", "test-validate-port-vendor")
	require.NoError(t, err)
	require.Empty(t, ret)
	ok, err := ValidatePortVendor("test-validate-port-vendor")
	require.NoError(t, err)
	require.False(t, ok)
}

func (suite *OvnClientTestSuite) testGetInterfacePodNs() {
	t := suite.T()
	t.Parallel()

	if !suite.enableOvsSandbox {
		return
	}

	// 1. port not exist
	ret, err := GetInterfacePodNs("port-not-exist")
	require.NoError(t, err)
	require.Empty(t, ret)
	// 2. port exist
	ret1, err := Exec("add-port", "br-int", "ovs-get-interface-pod-ns")
	require.NoError(t, err)
	require.Empty(t, ret1)
	err = ovsAdd("port", "ovs-get-interface-pod-ns", "tag", "10")
	require.NoError(t, err)
	ret, err = GetInterfacePodNs("ovs-get-interface-pod-ns")
	require.NoError(t, err)
	require.Empty(t, ret)
}

func (suite *OvnClientTestSuite) testConfigInterfaceMirror() {
	t := suite.T()
	t.Parallel()

	if !suite.enableOvsSandbox {
		return
	}

	err := ConfigInterfaceMirror(true, "open", "m0")
	require.NoError(t, err)

	err = ConfigInterfaceMirror(false, "close", "m0")
	require.NoError(t, err)
}

func (suite *OvnClientTestSuite) testGetResidualInternalPorts() {
	t := suite.T()
	t.Parallel()

	if !suite.enableOvsSandbox {
		return
	}

	// 1. port not exist
	ret := GetResidualInternalPorts()
	require.Empty(t, ret)

	// 2. port exist
	ret1, err := Exec("add-port", "br-int", "ovs-get-residual-internal-ports", "type=internal")
	require.NoError(t, err)
	require.Empty(t, ret1)
	ret = GetResidualInternalPorts()
	require.NotEmpty(t, ret)
}

func (suite *OvnClientTestSuite) testClearPortQosBinding() {
	t := suite.T()
	t.Parallel()

	if !suite.enableOvsSandbox {
		return
	}
	// 1. port not exist
	err := ClearPortQosBinding("ovs-clear-port-qos-binding")
	require.NoError(t, err)
	// 2. port exist
	ret1, err := Exec("add-port", "br-int", "ovs-clear-port-qos-binding")
	require.NoError(t, err)
	require.Empty(t, ret1)
	err = ovsAdd("port", "ovs-clear-port-qos-binding", "qos=qos-uuid")
	require.NoError(t, err)
	err = ClearPortQosBinding("ovs-clear-port-qos-binding")
	require.NoError(t, err)
}

func (suite *OvnClientTestSuite) testOvsListExternalIDs() {
	t := suite.T()
	t.Parallel()

	if !suite.enableOvsSandbox {
		return
	}

	// 1. port not exist
	ret, err := ListExternalIDs("port")
	require.NoError(t, err)
	require.Empty(t, ret)

	// 2. port exist
	ret1, err := Exec("add-port", "br-int", "ovs-list-external-ids", "external_ids:iface-id=ovs-list-external-ids")
	require.NoError(t, err)
	require.Empty(t, ret1)
	ret, err = ListExternalIDs("port")
	require.NoError(t, err)
	require.NotEmpty(t, ret)
}

func (suite *OvnClientTestSuite) testListQosQueueIDs() {
	t := suite.T()
	t.Parallel()

	if !suite.enableOvsSandbox {
		return
	}

	// 1. qos not exist
	ret, err := ListQosQueueIDs()
	require.NoError(t, err)
	require.Empty(t, ret)

	// 2. qos exist
	ret1, err := Exec("add-port", "br-int", "ovs-list-qos-queue-ids", "qos", "queues:0!=list-qos-queue-ids")
	require.NoError(t, err)
	require.Empty(t, ret1)
	ret, err = ListQosQueueIDs()
	require.NoError(t, err)
	require.NotEmpty(t, ret)
}

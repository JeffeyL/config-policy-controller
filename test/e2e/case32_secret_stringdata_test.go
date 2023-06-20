// Copyright Contributors to the Open Cluster Management project

package e2e

import (
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"open-cluster-management.io/config-policy-controller/test/utils"
)

const (
	case32ConfigPolicyName string = "case32config"
	case32CreatePolicyYaml string = "../resources/case32_secret_stringdata/case32_create_secret.yaml"
)

// create config policy that has configmap with stringData
// verify creation of configuration, check that the configmap has been created
// check that configmap data is base64 encoded

var _ = FDescribe("Test converted stringData being decoded before comparison for Secrets", Ordered, func() {
	It("Config should be created properly on the managed cluster", func() {
		By("Creating " + case32ConfigPolicyName + " on managed")
		utils.Kubectl("apply", "-f", case32CreatePolicyYaml, "-n", testNamespace)
		cfg := utils.GetWithTimeout(clientManagedDynamic, gvrConfigPolicy,
			case32ConfigPolicyName, testNamespace, true, defaultTimeoutSeconds)
		Expect(cfg).NotTo(BeNil())
	})

	It("Verifies the config policy is initially compliant "+case32ConfigPolicyName+" in "+testNamespace, func() {
		By("Waiting for " + case32ConfigPolicyName + " to become Compliant for up to " + strconv.Itoa(defaultTimeoutSeconds) + " seconds")
		Eventually(func() interface{} {
			cfgplc := utils.GetWithTimeout(
				clientManagedDynamic, gvrConfigPolicy, case32ConfigPolicyName, testNamespace, true, defaultTimeoutSeconds,
			)

			return utils.GetComplianceState(cfgplc)
		}, defaultTimeoutSeconds, 1).Should(Equal("Compliant"))
	})

	It("Verifies the config policy stays compliant "+case32ConfigPolicyName+" in "+testNamespace, func() {
		By("Making sure " + case32ConfigPolicyName + " is consistently Compliant for " + strconv.Itoa(defaultTimeoutSeconds) + " seconds")
		Consistently(func() interface{} {
			cfgplc := utils.GetWithTimeout(
				clientManagedDynamic, gvrConfigPolicy, case32ConfigPolicyName, testNamespace, true, defaultTimeoutSeconds,
			)

			return utils.GetComplianceState(cfgplc)
		}, defaultTimeoutSeconds, 1).ShouldNot(Equal("NonCompliant"))
	})

	AfterAll(func() {
		utils.Kubectl("delete", "configurationpolicy", case32ConfigPolicyName, "-n", testNamespace)
		utils.Kubectl("delete", "secret", "htpasswd-secret", "-n", "openshift-config")
	})
})

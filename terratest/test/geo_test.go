package test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGeo(t *testing.T) {
	t.Parallel()

	var coreDNSPods []corev1.Pod

	// Path to the Kubernetes resource config we will test
	kubeResourcePath, err := filepath.Abs("../example/geo.yaml")
	require.NoError(t, err)

	// To ensure we can reuse the resource config on the same cluster to test different scenarios, we setup a unique
	// namespace for the resources for this test.
	// Note that namespaces must be lowercase.
	namespaceName := fmt.Sprintf("coredns-test-%s", strings.ToLower(random.UniqueId()))

	options := k8s.NewKubectlOptions("", "", namespaceName)
	mainNsOptions := k8s.NewKubectlOptions("", "", "coredns")
	podFilter := metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=coredns",
	}

	k8s.CreateNamespace(t, options, namespaceName)

	defer k8s.DeleteNamespace(t, options, namespaceName)

	defer k8s.KubectlDelete(t, options, kubeResourcePath)

	k8s.KubectlApply(t, options, kubeResourcePath)

	k8s.WaitUntilNumPodsCreated(t, mainNsOptions, podFilter, 1, 60, 1*time.Second)

	coreDNSPods = k8s.ListPods(t, mainNsOptions, podFilter)

	for _, pod := range coreDNSPods {
		k8s.WaitUntilPodAvailable(t, mainNsOptions, pod.Name, 60, 1*time.Second)
	}

	t.Run("site1 gets site1 endpoints", func(t *testing.T) {
		clientIP := "192.200.1.50"
		actualIP, err := DigIPs(t, "localhost", 1053, "geo.example.org", dns.TypeA, clientIP)
		require.NoError(t, err)
		assert.NotContains(t, actualIP, "192.200.2.10")
	})
	t.Run("site2 gets site2 endpoints", func(t *testing.T) {
		clientIP := "192.200.2.30"
		actualIP, err := DigIPs(t, "localhost", 1053, "geo.example.org", dns.TypeA, clientIP)
		require.NoError(t, err)
		assert.NotContains(t, actualIP, "192.200.1.5")
	})
	t.Run("outside DC client gets all endpoints", func(t *testing.T) {
		clientIP := "192.100.1.15"
		actualIP, err := DigIPs(t, "localhost", 1053, "geo.example.org", dns.TypeA, clientIP)
		require.NoError(t, err)
		assert.Equal(t, len(actualIP), 4)
	})
}

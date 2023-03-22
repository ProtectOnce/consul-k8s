package v1alpha1

import (
	"testing"
	"time"

	"github.com/hashicorp/consul-k8s/control-plane/api/common"
	capi "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestControlPlaneRequestLimit_ToConsul(t *testing.T) {
	cases := map[string]struct {
		input    *ControlPlaneRequestLimit
		expected *capi.RateLimitIPConfigEntry
	}{
		"empty fields": {
			&ControlPlaneRequestLimit{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo",
				},
				Spec: ControlPlaneRequestLimitSpec{
					Mode: "disabled",
					ReadRate: 0,
					WriteRate: 0,
				},
			},
			&capi.RateLimitIPConfigEntry{
				Name: "foo",
				Kind: capi.RateLimitIPConfig,
				Mode: "disabled",
				Meta: map[string]string{
					common.DatacenterKey: "datacenter",
					common.SourceKey:     common.SourceValue,
				},
				ReadRate: 0,
				WriteRate: 0,
			},
		},
		"every field set": {
			&ControlPlaneRequestLimit{
				ObjectMeta: metav1.ObjectMeta{
					Name: "foo",
				},
				Spec:ControlPlaneRequestLimitSpec{
					Mode: "permissive",
					ReadRate: 100.0,
					WriteRate: 100.0,
					ACL: &ReadWriteRatesConfig{
						ReadRate: 100.0,
						WriteRate: 100.0,
					},
					Catalog: &ReadWriteRatesConfig{
						ReadRate: 100.0,
						WriteRate: 100.0,
					},
					ConfigEntry: &ReadWriteRatesConfig{
						ReadRate: 100.0,
						WriteRate: 100.0,
					},
					ConnectCA: &ReadWriteRatesConfig{
						ReadRate: 100.0,
						WriteRate: 100.0,
					},
					Coordinate: &ReadWriteRatesConfig{
						ReadRate: 100.0,
						WriteRate: 100.0,
					},
					DiscoveryChain: &ReadWriteRatesConfig{
						ReadRate: 100.0,
						WriteRate: 100.0,
					},
					Health: &ReadWriteRatesConfig{
						ReadRate: 100.0,
						WriteRate: 100.0,
					},
					Intention: &ReadWriteRatesConfig{
						ReadRate: 100.0,
						WriteRate: 100.0,
					},
					KV: &ReadWriteRatesConfig{
						ReadRate: 100.0,
						WriteRate: 100.0,
					},
					Tenancy: &ReadWriteRatesConfig{
						ReadRate: 100.0,
						WriteRate: 100.0,
					},
					PreparedQuery: &ReadWriteRatesConfig{
						ReadRate: 100.0,
						WriteRate: 100.0,
					},
					Session: &ReadWriteRatesConfig{
						ReadRate: 100.0,
						WriteRate: 100.0,
					},
					Txn: &ReadWriteRatesConfig{
						ReadRate: 100.0,
						WriteRate: 100.0,
					},	
				},
			},
			&capi.RateLimitIPConfigEntry{
				Kind:     capi.RateLimitIPConfig,
				Name:     "foo",
				Mode: 	  "permissive",
				ReadRate: 100.0,
				WriteRate: 100.0,
				Meta: map[string]string{
					common.DatacenterKey: "datacenter",
					common.SourceKey:     common.SourceValue,
				},
				ACL: &capi.ReadWriteRatesConfig{
					ReadRate: 100.0,
					WriteRate: 100.0,
				},
				Catalog: &capi.ReadWriteRatesConfig{
					ReadRate: 100.0,
					WriteRate: 100.0,
				},
				ConfigEntry: &capi.ReadWriteRatesConfig{
					ReadRate: 100.0,
					WriteRate: 100.0,
				},
				ConnectCA: &capi.ReadWriteRatesConfig{
					ReadRate: 100.0,
					WriteRate: 100.0,
				},
				Coordinate: &capi.ReadWriteRatesConfig{
					ReadRate: 100.0,
					WriteRate: 100.0,
				},
				DiscoveryChain: &capi.ReadWriteRatesConfig{
					ReadRate: 100.0,
					WriteRate: 100.0,
				},
				Health: &capi.ReadWriteRatesConfig{
					ReadRate: 100.0,
					WriteRate: 100.0,
				},
				Intention: &capi.ReadWriteRatesConfig{
					ReadRate: 100.0,
					WriteRate: 100.0,
				},
				KV: &capi.ReadWriteRatesConfig{
					ReadRate: 100.0,
					WriteRate: 100.0,
				},
				Tenancy: &capi.ReadWriteRatesConfig{
					ReadRate: 100.0,
					WriteRate: 100.0,
				},
				PreparedQuery: &capi.ReadWriteRatesConfig{
					ReadRate: 100.0,
					WriteRate: 100.0,
				},
				Session: &capi.ReadWriteRatesConfig{
					ReadRate: 100.0,
					WriteRate: 100.0,
				},
				Txn: &capi.ReadWriteRatesConfig{
					ReadRate: 100.0,
					WriteRate: 100.0,
				},	
			},
		},
	}

	for name, testCase := range cases {
		t.Run(name, func(t *testing.T) {
			output := testCase.input.ToConsul("datacenter")
			require.Equal(t, testCase.expected, output)
		})
	}
}

func TestControlPlaneRequestLimit_MatchesConsul(t *testing.T) {
	cases := map[string]struct {
		internal *ControlPlaneRequestLimit
		consul   capi.ConfigEntry
		matches  bool
	}{
		"empty fields matches": {
			&ControlPlaneRequestLimit{
				ObjectMeta: metav1.ObjectMeta{
					Name: "my-test-service",
				},
				Spec: ControlPlaneRequestLimitSpec{},
			},
			&capi.RateLimitIPConfigEntry{
				Kind:        capi.RateLimitIPConfig,
				Name:        "my-test-service",
				Namespace:   "namespace",
				CreateIndex: 1,
				ModifyIndex: 2,
				Meta: map[string]string{
					common.SourceKey:     common.SourceValue,
					common.DatacenterKey: "datacenter",
				},
			},
			true,
		},
		"all fields populated matches": {
			&ControlPlaneRequestLimit{
				ObjectMeta: metav1.ObjectMeta{
					Name: "my-test-service",
				},
				Spec:ControlPlaneRequestLimitSpec{
					Mode: "permissive",
					ReadRate: 100.0,
					WriteRate: 100.0,
					ACL: &ReadWriteRatesConfig{
						ReadRate: 100.0,
						WriteRate: 100.0,
					},
					Catalog: &ReadWriteRatesConfig{
						ReadRate: 100.0,
						WriteRate: 100.0,
					},
					ConfigEntry: &ReadWriteRatesConfig{
						ReadRate: 100.0,
						WriteRate: 100.0,
					},
					ConnectCA: &ReadWriteRatesConfig{
						ReadRate: 100.0,
						WriteRate: 100.0,
					},
					Coordinate: &ReadWriteRatesConfig{
						ReadRate: 100.0,
						WriteRate: 100.0,
					},
					DiscoveryChain: &ReadWriteRatesConfig{
						ReadRate: 100.0,
						WriteRate: 100.0,
					},
					Health: &ReadWriteRatesConfig{
						ReadRate: 100.0,
						WriteRate: 100.0,
					},
					Intention: &ReadWriteRatesConfig{
						ReadRate: 100.0,
						WriteRate: 100.0,
					},
					KV: &ReadWriteRatesConfig{
						ReadRate: 100.0,
						WriteRate: 100.0,
					},
					Tenancy: &ReadWriteRatesConfig{
						ReadRate: 100.0,
						WriteRate: 100.0,
					},
					PreparedQuery: &ReadWriteRatesConfig{
						ReadRate: 100.0,
						WriteRate: 100.0,
					},
					Session: &ReadWriteRatesConfig{
						ReadRate: 100.0,
						WriteRate: 100.0,
					},
					Txn: &ReadWriteRatesConfig{
						ReadRate: 100.0,
						WriteRate: 100.0,
					},	
				},
			},
			&capi.RateLimitIPConfigEntry{
				Kind:     capi.RateLimitIPConfig,
				Name:     "my-test-service",
				Mode: 	  "permissive",
				ReadRate: 100.0,
				WriteRate: 100.0,
				Meta: map[string]string{
					common.DatacenterKey: "datacenter",
					common.SourceKey:     common.SourceValue,
				},
				ACL: &capi.ReadWriteRatesConfig{
					ReadRate: 100.0,
					WriteRate: 100.0,
				},
				Catalog: &capi.ReadWriteRatesConfig{
					ReadRate: 100.0,
					WriteRate: 100.0,
				},
				ConfigEntry: &capi.ReadWriteRatesConfig{
					ReadRate: 100.0,
					WriteRate: 100.0,
				},
				ConnectCA: &capi.ReadWriteRatesConfig{
					ReadRate: 100.0,
					WriteRate: 100.0,
				},
				Coordinate: &capi.ReadWriteRatesConfig{
					ReadRate: 100.0,
					WriteRate: 100.0,
				},
				DiscoveryChain: &capi.ReadWriteRatesConfig{
					ReadRate: 100.0,
					WriteRate: 100.0,
				},
				Health: &capi.ReadWriteRatesConfig{
					ReadRate: 100.0,
					WriteRate: 100.0,
				},
				Intention: &capi.ReadWriteRatesConfig{
					ReadRate: 100.0,
					WriteRate: 100.0,
				},
				KV: &capi.ReadWriteRatesConfig{
					ReadRate: 100.0,
					WriteRate: 100.0,
				},
				Tenancy: &capi.ReadWriteRatesConfig{
					ReadRate: 100.0,
					WriteRate: 100.0,
				},
				PreparedQuery: &capi.ReadWriteRatesConfig{
					ReadRate: 100.0,
					WriteRate: 100.0,
				},
				Session: &capi.ReadWriteRatesConfig{
					ReadRate: 100.0,
					WriteRate: 100.0,
				},
				Txn: &capi.ReadWriteRatesConfig{
					ReadRate: 100.0,
					WriteRate: 100.0,
				},	
			},
			true,
		},
		"mismatched types does not match": {
			&ControlPlaneRequestLimit{
				ObjectMeta: metav1.ObjectMeta{
					Name: "my-test-service",
				},
				Spec: ControlPlaneRequestLimitSpec{},
			},
			&capi.ProxyConfigEntry{
				Kind:        capi.RateLimitIPConfig,
				Name:        "my-test-service",
				Namespace:   "namespace",
				CreateIndex: 1,
				ModifyIndex: 2,
			},
			false,
		},
	}

	for name, testCase := range cases {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, testCase.matches, testCase.internal.MatchesConsul(testCase.consul))
		})
	}
}

func TestControlPlaneRequestLimit_AddFinalizer(t *testing.T) {
	controlPlaneRequestLimit := &ControlPlaneRequestLimit{}
	controlPlaneRequestLimit.AddFinalizer("finalizer")
	require.Equal(t, []string{"finalizer"}, controlPlaneRequestLimit.ObjectMeta.Finalizers)
}

func TestControlPlaneRequestLimit_RemoveFinalizer(t *testing.T) {
	controlPlaneRequestLimit := &ControlPlaneRequestLimit{
		ObjectMeta: metav1.ObjectMeta{
			Finalizers: []string{"f1", "f2"},
		},
	}
	controlPlaneRequestLimit.RemoveFinalizer("f1")
	require.Equal(t, []string{"f2"}, controlPlaneRequestLimit.ObjectMeta.Finalizers)
}

func TestControlPlaneRequestLimit_SetSyncedCondition(t *testing.T) {
	controlPlaneRequestLimit := &ControlPlaneRequestLimit{}
	controlPlaneRequestLimit.SetSyncedCondition(corev1.ConditionTrue, "reason", "message")

	require.Equal(t, corev1.ConditionTrue, controlPlaneRequestLimit.Status.Conditions[0].Status)
	require.Equal(t, "reason", controlPlaneRequestLimit.Status.Conditions[0].Reason)
	require.Equal(t, "message", controlPlaneRequestLimit.Status.Conditions[0].Message)
	now := metav1.Now()
	require.True(t, controlPlaneRequestLimit.Status.Conditions[0].LastTransitionTime.Before(&now))
}

func TestControlPlaneRequestLimit_SetLastSyncedTime(t *testing.T) {
	controlPlaneRequestLimit := &ControlPlaneRequestLimit{}
	syncedTime := metav1.NewTime(time.Now())
	controlPlaneRequestLimit.SetLastSyncedTime(&syncedTime)

	require.Equal(t, &syncedTime, controlPlaneRequestLimit.Status.LastSyncedTime)
}

func TestControlPlaneRequestLimit_GetSyncedConditionStatus(t *testing.T) {
	cases := []corev1.ConditionStatus{
		corev1.ConditionUnknown,
		corev1.ConditionFalse,
		corev1.ConditionTrue,
	}
	for _, status := range cases {
		t.Run(string(status), func(t *testing.T) {
			controlPlaneRequestLimit := &ControlPlaneRequestLimit{
				Status: Status{
					Conditions: []Condition{{
						Type:   ConditionSynced,
						Status: status,
					}},
				},
			}

			require.Equal(t, status, controlPlaneRequestLimit.SyncedConditionStatus())
		})
	}
}

func TestControlPlaneRequestLimit_GetConditionWhenStatusNil(t *testing.T) {
	require.Nil(t, (&ControlPlaneRequestLimit{}).GetCondition(ConditionSynced))
}

func TestControlPlaneRequestLimit_SyncedConditionStatusWhenStatusNil(t *testing.T) {
	require.Equal(t, corev1.ConditionUnknown, (&ControlPlaneRequestLimit{}).SyncedConditionStatus())
}

func  TestControlPlaneRequestLimit_SyncedConditionWhenStatusNil(t *testing.T) {
	status, reason, message := (&ControlPlaneRequestLimit{}).SyncedCondition()
	require.Equal(t, corev1.ConditionUnknown, status)
	require.Equal(t, "", reason)
	require.Equal(t, "", message)
}

func TestControlPlaneRequestLimit_ConsulKind(t *testing.T) {
	require.Equal(t, capi.RateLimitIPConfig, (&ControlPlaneRequestLimit{}).ConsulKind())
}

func TestControlPlaneRequestLimit_KubeKind(t *testing.T) {
	require.Equal(t, "controlplanerequestlimit", (&ControlPlaneRequestLimit{}).KubeKind())
}

func TestControlPlaneRequestLimit_ConsulName(t *testing.T) {
	require.Equal(t, "foo", (&ControlPlaneRequestLimit{ObjectMeta: metav1.ObjectMeta{Name: "foo"}}).ConsulName())
}

func TestControlPlaneRequestLimit_KubernetesName(t *testing.T) {
	require.Equal(t, "foo", (&ControlPlaneRequestLimit{ObjectMeta: metav1.ObjectMeta{Name: "foo"}}).KubernetesName())
}

func TestControlPlaneRequestLimit_ConsulNamespace(t *testing.T) {
	require.Equal(t, "bar", (&ControlPlaneRequestLimit{ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "bar"}}).ConsulMirroringNS())
}

func TestControlPlaneRequestLimit_ConsulGlobalResource(t *testing.T) {
	require.False(t, (&ControlPlaneRequestLimit{}).ConsulGlobalResource())
}

func TestControlPlaneRequestLimit_ObjectMeta(t *testing.T) {
	meta := metav1.ObjectMeta{
		Name:      "name",
		Namespace: "namespace",
	}
	controlPlaneRequestLimit := &ControlPlaneRequestLimit{
		ObjectMeta: meta,
	}
	require.Equal(t, meta, controlPlaneRequestLimit.GetObjectMeta())
}

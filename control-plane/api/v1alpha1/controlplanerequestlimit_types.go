package v1alpha1

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/consul-k8s/control-plane/api/common"
	capi "github.com/hashicorp/consul/api"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ControlPlaneRequestLimitKubeKind = "controlplanerequestlimit"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

func init() {
	SchemeBuilder.Register(&ControlPlaneRequestLimit{}, &ControlPlaneRequestLimitList{})
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ControlPlaneRequestLimit is the Schema for the controlplanerequestlimits API
type ControlPlaneRequestLimit struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   			  ControlPlaneRequestLimitSpec   `json:"spec,omitempty"`
	Status 			  `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ControlPlaneRequestLimitList contains a list of ControlPlaneRequestLimit
type ControlPlaneRequestLimitList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ControlPlaneRequestLimit `json:"items"`
}

// ControlPlaneRequestLimitSpec defines the desired state of ControlPlaneRequestLimit
type ControlPlaneRequestLimitSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// {permissive, enforcing, disabled}
	Mode string

	// overall limits
	ReadRate  float64
	WriteRate float64

	//limits specific to a type of call
	ACL            *ReadWriteRatesConfig `json:acl",omitempty"`
	Catalog        *ReadWriteRatesConfig `json:catalog",omitempty"`
	ConfigEntry    *ReadWriteRatesConfig `json:configEntry",omitempty"`
	ConnectCA      *ReadWriteRatesConfig `json:connectCA",omitempty"`
	Coordinate     *ReadWriteRatesConfig `json:coordinate",omitempty"`
	DiscoveryChain *ReadWriteRatesConfig `json:discoveryChain",omitempty"`
	Health         *ReadWriteRatesConfig `json:health",omitempty"`
	Intention      *ReadWriteRatesConfig `json:intention",omitempty"`
	KV             *ReadWriteRatesConfig `json:kv",omitempty"`
	Tenancy        *ReadWriteRatesConfig `json:tenancy",omitempty"`
	PreparedQuery  *ReadWriteRatesConfig `json:perparedQuery",omitempty"`
	Session        *ReadWriteRatesConfig `json:session",omitempty"`
	Txn            *ReadWriteRatesConfig `json:txn",omitempty"`
}

type ReadWriteRatesConfig struct {
	ReadRate  float64
	WriteRate float64
}


func (cprl *ControlPlaneRequestLimit) GetObjectMeta() metav1.ObjectMeta {
	return cprl.ObjectMeta
}

// AddFinalizer adds a finalizer to the list of finalizers.
func (cprl *ControlPlaneRequestLimit) AddFinalizer(name string) {
	cprl.ObjectMeta.Finalizers = append(cprl.Finalizers(), name)
}

// RemoveFinalizer removes this finalizer from the list.
func (cprl *ControlPlaneRequestLimit) RemoveFinalizer(name string) {
	var newFinalizers []string
	for _, oldF := range cprl.Finalizers() {
		if oldF != name {
			newFinalizers = append(newFinalizers, oldF)
		}
	}
	cprl.ObjectMeta.Finalizers = newFinalizers
}

// Finalizers returns the list of finalizers for this object.
func (cprl *ControlPlaneRequestLimit) Finalizers() []string {
	return cprl.ObjectMeta.Finalizers
}

// ConsulKind returns the Consul config entry kind, i.e. service-defaults, not
// servicedefaults.
func (cprl *ControlPlaneRequestLimit) ConsulKind() string {
	return capi.RateLimitIPConfig
}

// ConsulGlobalResource returns if the resource exists in the default
// Consul namespace only.
func (cprl *ControlPlaneRequestLimit) ConsulGlobalResource() bool {
	return false
}

// ConsulMirroringNS returns the Consul namespace that the config entry should
// be created in if namespaces and mirroring are enabled.
func (cprl *ControlPlaneRequestLimit) ConsulMirroringNS() string {
	return cprl.Namespace
}

// KubeKind returns the Kube config entry kind, i.e. controlplanerequetlimit, not
// control-plane-request-limit.
func (cprl *ControlPlaneRequestLimit) KubeKind() string {
	return ControlPlaneRequestLimitKubeKind
}

// ConsulName returns the name of the config entry as saved in Consul.
// This may be different than KubernetesName() in the case of a ServiceIntentions
// config entry.
func (cprl *ControlPlaneRequestLimit) ConsulName() string {
	return cprl.ObjectMeta.Name
}

// KubernetesName returns the name of the Kubernetes resource.
func (cprl *ControlPlaneRequestLimit) KubernetesName() string {
	return cprl.ObjectMeta.Name
}

// SetSyncedCondition updates the synced condition.
func (cprl *ControlPlaneRequestLimit) SetSyncedCondition(status corev1.ConditionStatus, reason string, message string) {
	cprl.Status.Conditions = Conditions{
		{
			Type:               ConditionSynced,
			Status:             status,
			LastTransitionTime: metav1.Now(),
			Reason:             reason,
			Message:            message,
		},
	}
}

// SetLastSyncedTime updates the last synced time.
func (cprl *ControlPlaneRequestLimit) SetLastSyncedTime(time *metav1.Time) {
	cprl.Status.LastSyncedTime = time
}

// SyncedCondition gets the synced condition.
func (cprl *ControlPlaneRequestLimit) SyncedCondition() (status corev1.ConditionStatus, reason string, message string) {
	cond := cprl.Status.GetCondition(ConditionSynced)
	if cond == nil {
		return corev1.ConditionUnknown, "", ""
	}
	return cond.Status, cond.Reason, cond.Message
}

// SyncedConditionStatus returns the status of the synced condition.
func (cprl *ControlPlaneRequestLimit) SyncedConditionStatus() corev1.ConditionStatus {
	condition := cprl.Status.GetCondition(ConditionSynced)
	if condition == nil {
		return corev1.ConditionUnknown
	}
	return condition.Status
}

// ToConsul converts the resource to the corresponding Consul API definition.
// Its return type is the generic ConfigEntry but a specific config entry
// type should be constructed e.g. RateLimitIPConfigEntry.
func (cprl *ControlPlaneRequestLimit) ToConsul(datacenter string) capi.ConfigEntry {
	return &capi.RateLimitIPConfigEntry{
		Kind:                      cprl.ConsulKind(),
		Name:                      cprl.ConsulName(),
		Mode:                  	   cprl.Spec.Mode,
		ReadRate:                  cprl.Spec.ReadRate,
		WriteRate:                 cprl.Spec.WriteRate,
		Meta:                      meta(datacenter),
		ACL:                       cprl.Spec.ACL.toConsul(),
		Catalog:                   cprl.Spec.Catalog.toConsul(),
		ConfigEntry:               cprl.Spec.ConfigEntry.toConsul(),
		ConnectCA:                 cprl.Spec.ConnectCA.toConsul(),
		Coordinate:                cprl.Spec.Coordinate.toConsul(),
		DiscoveryChain:      	   cprl.Spec.DiscoveryChain.toConsul(),
		Health:      			   cprl.Spec.Health.toConsul(),
		Intention:      		   cprl.Spec.Intention.toConsul(),
		KV:  					   cprl.Spec.KV.toConsul(),
		Tenancy:            	   cprl.Spec.Tenancy.toConsul(),
		PreparedQuery:      	   cprl.Spec.PreparedQuery.toConsul(),
		Session:  				   cprl.Spec.Session.toConsul(),
		Txn:            		   cprl.Spec.Txn.toConsul(),
	}
}

func (cprl *ReadWriteRatesConfig) toConsul() *capi.ReadWriteRatesConfig {
	if cprl == nil {
		return nil
	}
	return &capi.ReadWriteRatesConfig{
		ReadRate: 		cprl.ReadRate,
		WriteRate:      cprl.WriteRate,
	}
}

// MatchesConsul returns true if the resource has the same fields as the Consul
// config entry.
func (cprl *ControlPlaneRequestLimit) MatchesConsul(candidate capi.ConfigEntry) bool {
	configEntry, ok := candidate.(*capi.RateLimitIPConfigEntry)
	if !ok {
		return false
	}
	// No datacenter is passed to ToConsul as we ignore the Meta field when checking for equality.
	return cmp.Equal(cprl.ToConsul(""), configEntry, cmpopts.IgnoreFields(capi.RateLimitIPConfigEntry{}, "Partition", "Namespace", "Meta", "ModifyIndex", "CreateIndex"), cmpopts.IgnoreUnexported(), cmpopts.EquateEmpty(),
		cmp.Comparer(transparentProxyConfigComparer))
}

// Validations will be performed in Consul
func (s *ControlPlaneRequestLimit) Validate(consulMeta common.ConsulMeta) error {
	return nil
}

// DefaultNamespaceFields has no behaviour here as control-plane-request-limit have no namespace specific fields.
func (s *ControlPlaneRequestLimit) DefaultNamespaceFields(_ common.ConsulMeta) {
}

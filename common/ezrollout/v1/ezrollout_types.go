/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EzRolloutSpec defines the desired state of EzRollout
type EzRolloutSpec struct {
	// Selector is a label query over pods that should match the replica count.
	// +required
	Selector *metav1.LabelSelector `json:"selector"`

	// OnlineVersion represents the version that is currently online
	// +optional
	OnlineVersion string `json:"onlineVersion,omitempty"`

	// ScaleUpMetrics contains the specifications for which to use to calculate the
	// desired replica count (the maximum replica count across all metrics will
	// be used).
	// +optional
	ScaleUpMetrics []autoscalingv2.MetricSpec `json:"scaleUpMetrics,omitempty"`

	// ScaleDownMetrics contains the specifications for which to use to calculate the
	// desired replica count (the maximum replica count across all metrics will
	// be used).
	// +optional
	ScaleDownMetrics []autoscalingv2.MetricSpec `json:"scaleDownMetrics,omitempty"`

	// OnlineScaler defines the behavior of the online version scaler
	// +optional
	OnlineScaler *EzOnlineScaler `json:"onlineScaler,omitempty"`

	// OfflineScaler defines the behavior of the offline version scaler
	// +optional
	OfflineScaler *EzOfflineScaler `json:"offlineScaler,omitempty"`
}

// EzOnlineScaler defines the scaling behavior for online version
type EzOnlineScaler struct {
	// MinReplicas is the lower limit for the number of replicas to which the autoscaler
	// can scale down.
	// +optional
	MinReplicas *int32 `json:"minReplicas,omitempty"`

	// MaxReplicas is the upper limit for the number of replicas to which the autoscaler can scale up.
	// +optional
	MaxReplicas *int32 `json:"maxReplicas,omitempty"`

	// Behavior configures the scaling behavior of the target
	// in both Up and Down directions (scaleUp and scaleDown fields respectively).
	// +optional
	Behavior *autoscalingv2.HorizontalPodAutoscalerBehavior `json:"behavior,omitempty"`
}

// EzOfflineScaler defines the scaling behavior for offline version
type EzOfflineScaler struct {
	// Deadline represents the time limit for scaling operations
	// +optional
	Deadline      int64 `json:"deadline,omitempty"`
	EnableScaleUp bool  `json:"enableScaleUp,omitempty"`
}

// EzRolloutStatus defines the observed state of EzRollout
type EzRolloutStatus struct {
	// Ready indicates that the rollout is ready
	Ready bool `json:"ready"`
	// LatestError indicates the latest error that occurred
	LatestError EzRolloutLatestError `json:"latestError"`
	// ObservedGeneration indicates the generation observed by controller
	ObservedGeneration int64 `json:"observedGeneration"`
}

type EzRolloutLatestError struct {
	// Message indicates when the latest error's info
	Message string `json:"message"`
	// Timestamp indicates when the latest error occurred
	Timestamp int64 `json:"timestamp"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// EzRollout is the Schema for the ezrollouts API
type EzRollout struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EzRolloutSpec   `json:"spec,omitempty"`
	Status EzRolloutStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EzRolloutList contains a list of EzRollout
type EzRolloutList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EzRollout `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EzRollout{}, &EzRolloutList{})
}

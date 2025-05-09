/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
// Important: Run "make generate-deepcopy" to regenerate code after modifying this file

// BuildSpec defines the list of tasks to be execute for a Build. From Camel K version 2, it would be more appropriate
// to think it as pipeline.
type BuildSpec struct {
	// The sequence of tasks (pipeline) to be performed.
	Tasks []Task `json:"tasks,omitempty"`
	// The configuration that should be used to perform the Build.
	// Deprecated: no longer in use in Camel K 2 - maintained for backward compatibility
	Configuration BuildConfiguration `json:"configuration,omitempty"`
	// The container image to be used to run the build.
	// Deprecated: no longer in use in Camel K 2 - maintained for backward compatibility
	ToolImage string `json:"toolImage,omitempty"`
	// The namespace where to run the builder Pod (must be the same of the operator in charge of this Build reconciliation).
	// Deprecated: no longer in use in Camel K 2 - maintained for backward compatibility
	BuilderPodNamespace string `json:"operatorNamespace,omitempty"`
	// Timeout defines the Build maximum execution duration.
	// The Build deadline is set to the Build start time plus the Timeout duration.
	// If the Build deadline is exceeded, the Build context is canceled,
	// and its phase set to BuildPhaseFailed.
	// +kubebuilder:validation:Format=duration
	Timeout metav1.Duration `json:"timeout,omitempty"`
	// the maximum amount of parallel running builds started by this operator instance
	// Deprecated: no longer in use in Camel K 2 - maintained for backward compatibility
	MaxRunningBuilds int32 `json:"maxRunningBuilds,omitempty"`
}

// Task represents the abstract task. Only one of the task should be configured to represent the specific task chosen.
type Task struct {
	// Application building

	// a BuilderTask, used to generate and build the project
	Builder *BuilderTask `json:"builder,omitempty"`

	// User customizable task execution. These are executed after the build and before the package task.
	Custom *UserTask `json:"custom,omitempty"`

	// Application pre publishing
	// a PackageTask, used to package the project
	Package *BuilderTask `json:"package,omitempty"`

	// Application Publishing

	// a BuildahTask, for Buildah strategy
	// Deprecated: use jib or a custom publishing strategy instead
	Buildah *BuildahTask `json:"buildah,omitempty"`
	// a KanikoTask, for Kaniko strategy
	// Deprecated: use jib or a custom publishing strategy instead
	Kaniko *KanikoTask `json:"kaniko,omitempty"`
	// a SpectrumTask, for Spectrum strategy
	// Deprecated: use jib or a custom publishing strategy instead
	Spectrum *SpectrumTask `json:"spectrum,omitempty"`
	// a S2iTask, for S2I strategy
	// Deprecated: use jib or a custom publishing strategy instead
	S2i *S2iTask `json:"s2i,omitempty"`
	// a JibTask, for Jib strategy
	Jib *JibTask `json:"jib,omitempty"`
}

// BaseTask is a base for the struct hierarchy.
type BaseTask struct {
	// name of the task
	Name string `json:"name,omitempty"`
	// The configuration that should be used to perform the Build.
	Configuration BuildConfiguration `json:"configuration,omitempty"`
}

// BuilderTask is the generic task in charge of building the application image.
type BuilderTask struct {
	BaseTask `json:",inline"`
	// the base image layer
	BaseImage string `json:"baseImage,omitempty"`
	// the configuration required for the runtime application
	Runtime RuntimeSpec `json:"runtime,omitempty"`
	// the list of dependencies to use for this build
	Dependencies []string `json:"dependencies,omitempty"`
	// the list of steps to execute (see pkg/builder/)
	Steps []string `json:"steps,omitempty"`
	// the configuration required by Maven for the application build phase
	Maven MavenBuildSpec `json:"maven,omitempty"`
	// workspace directory to use
	BuildDir string `json:"buildDir,omitempty"`
	// the sources to add at build time
	Sources []SourceSpec `json:"sources,omitempty"`
	// the configuration of the project to build on Git
	Git *GitConfigSpec `json:"git,omitempty"`
}

// GitConfigSpec defines the Git configuration of a project.
type GitConfigSpec struct {
	// the URL of the project
	URL string `json:"url,omitempty"`
	// the Kubernetes secret where token is stored
	Secret string `json:"secret,omitempty"`
}

// MavenBuildSpec defines the Maven configuration plus additional repositories to use.
type MavenBuildSpec struct {
	// base Maven specification
	MavenSpec `json:",inline"`
	// additional repositories
	Repositories []Repository `json:"repositories,omitempty"`
	// Servers (auth)
	Servers []Server `json:"servers,omitempty"`
}

// PublishTask image publish configuration.
type PublishTask struct {
	// can be useful to share info with other tasks
	ContextDir string `json:"contextDir,omitempty"`
	// base image layer
	BaseImage string `json:"baseImage,omitempty"`
	// final image name
	Image string `json:"image,omitempty"`
	// where to publish the final image
	Registry RegistrySpec `json:"registry,omitempty"`
}

// BuildahTask is used to configure Buildah.
// Deprecated: no longer in use.
type BuildahTask struct {
	BaseTask    `json:",inline"`
	PublishTask `json:",inline"`
	// The platform of build image
	Platform string `json:"platform,omitempty"`
	// log more information
	Verbose *bool `json:"verbose,omitempty"`
	// docker image to use
	ExecutorImage string `json:"executorImage,omitempty"`
}

// KanikoTask is used to configure Kaniko.
// Deprecated: no longer in use.
type KanikoTask struct {
	BaseTask    `json:",inline"`
	PublishTask `json:",inline"`
	// log more information
	Verbose *bool `json:"verbose,omitempty"`
	// use a cache
	Cache KanikoTaskCache `json:"cache,omitempty"`
	// docker image to use
	ExecutorImage string `json:"executorImage,omitempty"`
}

// KanikoTaskCache is used to configure Kaniko cache.
// Deprecated: no longer in use.
type KanikoTaskCache struct {
	// true if a cache is enabled
	Enabled *bool `json:"enabled,omitempty"`
	// the PVC used to store the cache
	PersistentVolumeClaim string `json:"persistentVolumeClaim,omitempty"`
}

// JibTask is used to configure Jib.
type JibTask struct {
	BaseTask    `json:",inline"`
	PublishTask `json:",inline"`
}

// SpectrumTask is used to configure Spectrum.
type SpectrumTask struct {
	BaseTask    `json:",inline"`
	PublishTask `json:",inline"`
}

// S2iTask is used to configure S2I.
type S2iTask struct {
	BaseTask    `json:",inline"`
	PublishTask `json:",inline"`
	// used by the ImageStream
	Tag string `json:"tag,omitempty"`
}

// UserTask is used to execute any generic custom operation.
type UserTask struct {
	BaseTask `json:",inline"`
	// the container image to use
	ContainerImage string `json:"image,omitempty"`
	// the user id used to run the container
	ContainerUserID *int64 `json:"userId,omitempty"`
	// the command to execute
	// Deprecated: use ContainerCommands
	ContainerCommand string `json:"command,omitempty"`
	// the command to execute
	ContainerCommands []string `json:"commands,omitempty"`
	// the desired image build name
	PublishingImage string `json:"publishingImage,omitempty"`
}

// BuildStatus defines the observed state of Build.
type BuildStatus struct {
	// ObservedGeneration is the most recent generation observed for this Build.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// describes the phase
	Phase BuildPhase `json:"phase,omitempty"`
	// the image name built
	Image string `json:"image,omitempty"`
	// the digest from image
	Digest string `json:"digest,omitempty"`
	// root image (the first image from which the incremental image has started)
	RootImage string `json:"rootImage,omitempty"`
	// the base image used for this build
	BaseImage string `json:"baseImage,omitempty"`
	// a list of artifacts contained in the build
	Artifacts []Artifact `json:"artifacts,omitempty"`
	// the error description (if any)
	Error string `json:"error,omitempty"`
	// the reason of the failure (if any)
	Failure *Failure `json:"failure,omitempty"`
	// the time when it started
	StartedAt *metav1.Time `json:"startedAt,omitempty"`
	// a list of conditions occurred during the build
	Conditions []BuildCondition `json:"conditions,omitempty"`
	// how long it took for the build
	// Change to Duration / ISO 8601 when CRD uses OpenAPI spec v3
	// https://github.com/OAI/OpenAPI-Specification/issues/845
	Duration string `json:"duration,omitempty"`
}

// BuildPhase -- .
type BuildPhase string

// BuildConditionType -- .
type BuildConditionType string

const (
	// BuildKind -- .
	BuildKind string = "Build"

	// BuildPhaseNone -- .
	BuildPhaseNone BuildPhase = ""
	// BuildPhaseInitialization -- .
	BuildPhaseInitialization BuildPhase = "Initialization"
	// BuildPhaseScheduling -- .
	BuildPhaseScheduling BuildPhase = "Scheduling"
	// BuildPhasePending -- .
	BuildPhasePending BuildPhase = "Pending"
	// BuildPhaseRunning -- .
	BuildPhaseRunning BuildPhase = "Running"
	// BuildPhaseSucceeded -- .
	BuildPhaseSucceeded BuildPhase = "Succeeded"
	// BuildPhaseFailed -- .
	BuildPhaseFailed BuildPhase = "Failed"
	// BuildPhaseInterrupted -- .
	BuildPhaseInterrupted = "Interrupted"
	// BuildPhaseError -- .
	BuildPhaseError BuildPhase = "Error"

	// BuildConditionScheduled --.
	BuildConditionScheduled BuildConditionType = "Scheduled"

	// BuildConditionReadyReason --.
	BuildConditionReadyReason string = "Ready"
	// BuildConditionWaitingReason --.
	BuildConditionWaitingReason string = "Waiting"
)

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=builds,scope=Namespaced,shortName=ikb,categories=kamel;camel
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`,description="The build phase"
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`,description="The time at which the build was created"
// +kubebuilder:printcolumn:name="Started",type=date,JSONPath=`.status.startedAt`,description="The time at which the build was last (re-)started"
// Change format to 'duration' when CRD uses OpenAPI spec v3 (https://github.com/OAI/OpenAPI-Specification/issues/845)
// +kubebuilder:printcolumn:name="Duration",type=string,JSONPath=`.status.duration`,description="The build last execution duration"
// +kubebuilder:printcolumn:name="Attempts",type=integer,JSONPath=`.status.failure.recovery.attempt`,description="The number of execution attempts"

// Build is the Schema for the builds API.
type Build struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BuildSpec   `json:"spec,omitempty"`
	Status BuildStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// BuildList contains a list of Build.
type BuildList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Build `json:"items"`
}

// BuildCondition describes the state of a resource at a certain point.
type BuildCondition struct {
	// Type of integration condition.
	Type BuildConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status corev1.ConditionStatus `json:"status"`
	// The last time this condition was updated.
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
	// Last time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// The reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`
	// A human-readable message indicating details about the transition.
	Message string `json:"message,omitempty"`
}

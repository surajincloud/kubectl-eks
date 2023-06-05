package kube

const (
	// capacity Type
	CapacityTypeLabel          = "eks.amazonaws.com/capacityType"
	KarpenterCapacityTypeLabel = "karpenter.sh/capacity-type"

	NodeGroupLabel = "eks.amazonaws.com/nodegroup"
	ComputeType    = "eks.amazonaws.com/compute-type" // to detect fargate nodes
	ArchLabel      = "kubernetes.io/arch"
	OsLabel        = "kubernetes.io/os"
	HostNameLabel  = "kubernetes.io/hostname"

	InstanceTypeLabel = "node.kubernetes.io/instance-type"

	ZoneLabel = "topology.kubernetes.io/zone"

	// Ami ID
	NodeGroupImage = "eks.amazonaws.com/nodegroup-image"
	KarpenterImage = "karpenter.k8s.aws/instance-ami-id"
)

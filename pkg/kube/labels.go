package kube

const (
	CapacityTypeLabel = "eks.amazonaws.com/capacityType"
	NodeGroupLabel    = "eks.amazonaws.com/nodegroup"
	ComputeType       = "eks.amazonaws.com/compute-type" // to detect fargate nodes
	ArchLabel         = "kubernetes.io/arch"
	OsLabel           = "kubernetes.io/os"
	HostNameLabel     = "kubernetes.io/hostname"

	InstanceTypeLabel = "node.kubernetes.io/instance-type"

	ZoneLabel      = "topology.kubernetes.io/zone"
	NodeGroupImage = "eks.amazonaws.com/nodegroup-image"
)

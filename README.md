# Kubectl EKS Plugin

`kubectl-eks` is kubectl plugin for amazon EKS

## Usage

* **List nodes but get more information**

```
kubectl eks nodes
```

* **List nodes but get more information**

```
kubectl eks ssm <name-of-the-node>
```

**Note**: required SSM IAM role needs to be present on the node. no need of aws ssm plugin.

## Installation

### Build from Source

```
git clone https://github.com/surajincloud/kubectl-eks
cd kubectl-eks
make
```

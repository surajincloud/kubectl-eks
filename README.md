# Kubectl EKS Plugin

`kubectl-eks` is kubectl plugin for amazon EKS

## Usage

### List nodes but get more information

```
kubectl eks nodes
```

### Access to EKS node via SSM

```
kubectl eks ssm <name-of-the-node>
```

**Note**: above command will only work if node IAM role has predefined IAM policy `AmazonSSMManagedInstanceCore` policy attached. Click [here](https://docs.aws.amazon.com/systems-manager/latest/userguide/setup-instance-profile.html) for more reference.

### List nodes but get more information


* List serviceaccount with IRSA information from all namespaces

```
kubectl eks irsa
```

* List serviceaccount with IRSA information from given namespace

```
kubectl eks irsa -n app-staging
```

## Installation

### Using Krew

* You can add custom index as shown below and install the plugin from there. We are planning to submit this plugin to official Krew index as well, you can track the progress [here](https://github.com/surajincloud/kubectl-eks/issues/3).

```
kubectl krew index add surajincloud git@github.com:surajincloud/krew-index.git
kubectl krew search eks
kubectl krew install surajincloud/kubectl-eks
```

### Download the binary

* Download the binary from the [release pages](https://github.com/surajincloud/kubectl-eks/releases).
* Place it into any of the location from the PATH.

### Build from Source

```
git clone https://github.com/surajincloud/kubectl-eks
cd kubectl-eks
make
```

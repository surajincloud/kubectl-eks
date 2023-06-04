# Kubectl EKS Plugin

`kubectl-eks` is kubectl plugin for amazon EKS

Check out docs website for more information: https://surajincloud.github.io/kubectl-eks/

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

# Installation

### Using Krew

* You can add custom index as shown below and install the plugin from there. We are planning to submit this plugin to official Krew index as well, you can track the progress [here](https://github.com/surajincloud/kubectl-eks/issues/3).

```
kubectl krew index add surajincloud git@github.com:surajincloud/kubectl-eks.git
kubectl krew search eks
kubectl krew install surajincloud/kubectl-eks
```


## Install from Source

```
git clone https://github.com/surajincloud/kubectl-eks
cd kubectl-eks
make
```

## Install from Releases

* Download latest binary from [release page](https://github.com/surajincloud/kubectl-eks/releases).

* make it executable

```
chmod +x kubectl-eks
```

* move it to one of the location from the PATH

```
mv kubectl-eks ~/.local/bin
```

## Install using Brew

```
brew tap surajincloud/tools
brew install kubectl-eks
```

## Verify Installation

* Verify the installation by running the following command.

```
$ kubectl eks --help
A kubectl plugin for Amazon EKS

Usage:
  kubectl-eks [command]

Available Commands:
  addons      A brief description of your command
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  irsa        A brief description of your command
  nodes       List all EKS Nodes
  ssm         Access given EKS node via SSM
  version     Print the version of kubectl-eks
...
...
```

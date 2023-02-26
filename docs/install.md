# Installation


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

TODO

you can contribute to this section, check out [issue #17](https://github.com/surajincloud/kubectl-eks/issues/17).

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
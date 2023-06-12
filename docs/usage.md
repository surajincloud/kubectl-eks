# Usage

> Note: `kubectl-eks` is able to read `region` from aws credentials file, aws profile and environment variable, optionally you can pass it via `--region` flag. it is also able to read cluster name from the kubeconfig context, optionally you can pass it via `--cluster` flag or `AWS_EKS_CLUSTER` environment variable.

## Creates Kubeconfig

```
kubectl eks kubeconfig --cluster your-cluster --region your-region --out
```

## Updates existing Kubeconfig

* fetch and update existing kubeconfig (~/.kube/config)

```
kubectl eks kubeconfig --cluster your-cluster --region your-region
```

## List addons

```
kubectl eks addons --cluster your-cluster --region your-region
```

## List serviceaccount with IRSA information from all namespaces

```
kubectl eks irsa
```

## List serviceaccount with IRSA information from given namespace

```
kubectl eks irsa -n app-staging
```

## List nodes but get more information

```
kubectl eks nodes
```
## Access to EKS node via SSM

```
kubectl eks ssm <name-of-the-node>
```

**Note**: above command will only work if node IAM role has predefined IAM policy AmazonSSMManagedInstanceCore policy attached. Click here for more reference.

## List fargate profiles

```
kubectl eks fargate --cluster your-cluster --region your-region
```

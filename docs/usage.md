# Usage

## List addons

```
kubectl eks addons --cluster-name your-cluster --region your-region
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
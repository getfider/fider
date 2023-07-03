# fider

![Version: 0.0.0](https://img.shields.io/badge/Version-0.0.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.21.1](https://img.shields.io/badge/AppVersion-0.21.1-informational?style=flat-square)

An open platform to collect and prioritize feedback

## Installation

### Add Helm repository

```bash
helm repo add fider https://footur.github.io/fider
helm repo update
```

## Configuration

The following table lists the configurable parameters of the chart and the default values.

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` | Defines rules for pod scheduling based on node or pod properties, allowing for affinity (preferred) or anti-affinity (avoidance) relationships between pods and nodes |
| autoscaling | object | `{"enabled":false,"maxReplicas":100,"minReplicas":1,"targetCPUUtilizationPercentage":80}` | Automatically adjusts the number of replicas of a deployment or a replica set based on predefined metrics or custom rules to handle varying workload demands |
| autoscaling.enabled | bool | `false` | Create a HPA for the server deployment |
| fider.env[0] | object | `{"name":"BASE_URL","value":"https://feedback.yourdomain.com"}` | Public Host Name |
| fider.env[1].name | string | `"EMAIL_NOREPLY"` |  |
| fider.env[1].value | string | `"noreply@yourdomain.com"` |  |
| fider.legalPages | object | `{"enabled":false,"privacy.md":"bar","terms.md":"foo"}` | https://fider.io/docs/how-to-show-legal-pages |
| fider.legalPages."privacy.md" | string | `"bar"` | Write your own Privacy Policy in Markdown here |
| fider.legalPages."terms.md" | string | `"foo"` | Write your own Terms of Service in Markdown here |
| fider.legalPages.enabled | bool | `false` | Enables legal pages |
| fider.secretEnv | object | `{"DATABASE_URL":"postgres://fider:s0m3g00dp4ssw0rd@db:5432/fider?sslmode=disable","JWT_SECRET":"VERY_STRONG_SECRET_SHOULD_BE_USED_HERE"}` | These environment variables are stored in a Kubernetes secret |
| fider.secretEnv.DATABASE_URL | string | `"postgres://fider:s0m3g00dp4ssw0rd@db:5432/fider?sslmode=disable"` | Connection string to the PostgreSQL database |
| fider.secretEnv.JWT_SECRET | string | `"VERY_STRONG_SECRET_SHOULD_BE_USED_HERE"` | Use a 512-bit secret here |
| fullnameOverride | string | `nil` | Override the fully qualified app name |
| image.pullPolicy | string | `"IfNotPresent"` | "IfNotPresent" to pull the image if no image with the specified tag exists on the node, "Always" to always pull the image or "Never" to try and use pre-pulled images |
| image.repository | string | `"getfider/fider"` | Repository to pull fider image from |
| imagePullSecrets | list | `[]` | Names of the Kubernetes secrets for imagePullSecrets |
| ingress.annotations | object | `{}` | Ingress annotations |
| ingress.className | string | `""` | The name of the Ingress Class associated with the ingress |
| ingress.enabled | bool | `false` | If `true``, an Ingress is created |
| ingress.hosts | list | `[{"host":"fider.local","paths":[{"path":"/","pathType":"prefix"}]}]` | Domain name Kubernetes Ingress rule looks for. Set it to the domain Fider will be hosted on |
| ingress.hosts[0] | object | `{"host":"fider.local","paths":[{"path":"/","pathType":"prefix"}]}` | List of domain names Kubernetes Ingress rule looks for. Set it to the domains in which Fider will be hosted on |
| ingress.hosts[0].paths | list | `[{"path":"/","pathType":"prefix"}]` | List of paths to use in Kubernetes Ingress rules |
| ingress.hosts[0].paths[0] | object | `{"path":"/","pathType":"prefix"}` | Path to use in the Ingress |
| ingress.hosts[0].paths[0].pathType | string | `"prefix"` | Ingress path type |
| ingress.tls | list | `[]` | Ingress TLS settings |
| kubeVersionOverride | string | `nil` | Override the version of Kubernetes being used in a cluster |
| nameOverride | string | `nil` | Override the name of the chart |
| nodeSelector | object | `{}` | Selects the nodes where a pod can be scheduled based on node labels |
| podAnnotations | object | `{}` | Key-value metadata for individual pods |
| podSecurityContext | object | `{}` | Defines the security settings and privileges for a pod |
| postgresql.enabled | bool | `false` | Enable the PostgreSQL subchart? |
| postgresql.ingress.enabled | bool | `false` | Set up ingress for external access (optional) |
| postgresql.persistence | object | `{"accessModes":["ReadWriteOnce"],"enabled":true,"size":"8Gi","storageClass":"standard"}` | Persistent volume configuration |
| postgresql.persistence.accessModes | list | `["ReadWriteOnce"]` | Set the access modes |
| postgresql.persistence.enabled | bool | `true` | Enables persistent storage |
| postgresql.persistence.size | string | `"8Gi"` | Set the size of the persistent volume |
| postgresql.persistence.storageClass | string | `"standard"` | Set the storage class |
| postgresql.postgresqlDatabase | string | `"mydatabase"` | PostgreSQL database |
| postgresql.postgresqlExtendedConf.max_connections | int | `200` | Maximum of DB connections |
| postgresql.postgresqlExtendedConf.shared_buffers | string | `"256MB"` | Shared buffers are a dedicated portion of memory used to cache frequently accessed data and improve database performance |
| postgresql.postgresqlPassword | string | `"mypassword"` | PostgreSQL password |
| postgresql.postgresqlUsername | string | `"myuser"` | PostgreSQL username |
| postgresql.resources | object | `{}` | Kubernetes ressources configuration |
| postgresql.service | object | `{"port":5432,"type":"ClusterIP"}` | Kubernetes service configuration |
| postgresql.service.port | int | `5432` | PostgreSQL port |
| postgresql.service.type | string | `"ClusterIP"` | The service type |
| replicaCount | int | `1` | The number of replicas to create |
| resources | object | `{}` | Kubernetes resources |
| securityContext | object | `{}` | Fider Container-level security-context |
| service | object | `{"port":80,"type":"ClusterIP"}` | Kubernetes service configuration |
| service.port | int | `80` | The HTTP service port |
| service.type | string | `"ClusterIP"` | The service type |
| serviceAccount.annotations | object | `{}` | Annotations to add to the service account |
| serviceAccount.create | bool | `true` | Specifies whether a service account should be created |
| serviceAccount.name | string | `""` | The name of the service account to use. If not set and create is true, a name is generated using the fullname template |
| tolerations | list | `[]` | Allows a pod to tolerate or ignore specific node taints, enabling it to be scheduled on tainted nodes |

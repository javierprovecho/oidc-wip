# GKE OIDC workload identity for external resources

## Dependency credentials for Kubernetes applications 

When developing software, security should always be a top priority. Software development often involves integrating with external dependencies and services, both first and third party. These external resources usually require authentication or credentials to access, which can become a complex and time-consuming process to manage.

## Injecting credentials

Kubernetes provides several options to inject a secret value into a workload (pod). One option is to use environment variables:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  containers:
  - name: my-container
    image: my-image
    env:
    - name: SECRET_API_KEY
      value: topSecretValue!123
```

Using a Kubernetes Secret object as source for the environment variable helps segregating sensitive information from the pod specification:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  containers:
  - name: my-container
    image: my-image
    env:
    - name: SECRET_API_KEY
      valueFrom:
        secretKeyRef:
          name: my-secret
          key: my-key
---
apiVersion: v1
kind: Secret
metadata:
  name: my-secret
type: Opaque
data:
  # base64 encoded value for 'topSecretValue!123'
  my-key: dG9wU2VjcmV0VmFsdWUhMTIz
```

Another option is to mount the secret as a volume in the pod. This allows the secret data to be accessed as a file within the container:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  containers:
  - name: my-container
    image: my-image
    volumeMounts:
    - name: my-secret-volume
      mountPath: /etc/my-secret
      # /etc/my-secret/my-key will contain 'topSecretValue!123'
  volumes:
  - name: my-secret-volume
    secret:
      secretName: my-secret
---
apiVersion: v1
kind: Secret
metadata:
  name: my-secret
type: Opaque
data:
  # base64 encoded value for 'topSecretValue!123'
  my-key: dG9wU2VjcmV0VmFsdWUhMTIz
```

Using a CSI driver like GoogleCloudPlatform/secrets-store-csi-driver-provider-gcp to mount secrets in a pod allows for more flexibility and extensibility in how secrets are stored and managed.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  containers:
  - name: my-container
    image: my-image
    volumeMounts:
    - name: my-secret-volume
      mountPath: /etc/my-secret
  volumes:
  - name: my-secret-volume
    csi:
      driver: secrets-store.csi.k8s.io
      readOnly: true
      volumeAttributes:
        secretProviderClass: my-secret
---
apiVersion: secrets-store.csi.x-k8s.io/v1
kind: SecretProviderClass
metadata:
  name: my-secret
spec:
  provider: gcp
  parameters:
    secrets: |
      - resourceName: "projects/$PROJECT_ID/secrets/my-secret/versions/latest"
        path: "my-key"
```

Although injecting secrets into Kubernetes pods can help secure applications, it still requires someone or something to generate and rotate the secrets. If the secrets are used for authenticating against external resources and are also required for development environments, this can add management overhead and increase the risk of secrets being exposed if not properly managed.

## Kubernetes Projected Service Account Tokens

Kubernetes automatically mounts a service account token in every pod by default, which can be used to authenticate with the Kubernetes API server and other resources. A service account in Kubernetes is used to authenticate and authorize operations within the cluster, such as API requests and communication between pods.

```shell
/ $ ls /var/run/secrets/kubernetes.io/serviceaccount/
ca.crt
namespace
token
```

Since 1.24, Kubernetes provides a way to project a service account token into a pod as a volume. The Kubernetes projected service account tokens feature allows users to customize the audience and duration of the tokens.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  containers:
  - name: my-container
    image: my-image
    volumeMounts:
    - name: my-token-volume
      mountPath: /var/run/secrets/tokens
      readOnly: true
  volumes:
  - name: my-token-volume
    projected:
      sources:
      - serviceAccountToken:
          path: external-resource-a-token
          expirationSeconds: 7200
          audience: external-resource-a
      - serviceAccountToken:
          path: external-resource-b-token
          expirationSeconds: 1800
          audience: external-resource-b
```

This token is in JWT format, compatible with OIDC and can be used to authenticate to external services.

## What is OIDC?

OIDC, or OpenID Connect, is a simple identity layer on top of the OAuth 2.0 protocol. It enables clients to both obtain basic profile information about the user and verify the identity of the user based on the authentication performed by an authorization server.

When it comes to Kubernetes projected service accounts, the following applies:

- Verification of the token's signature is performed using the issuer's public keys exposed through an OIDC discovery endpoint.
- The profile information contains pod information where the token is mounted, including the namespace, pod name, and service account name.

These are the JWT claims present in a token generated by a GKE cluster:

```json
{
  "aud": [
    "external-resource-a"
  ],
  "exp": 0,
  "iat": 0,
  "iss": "https://container.googleapis.com/v1/projects/PROJECT/locations/LOCATION/clusters/CLUSTER",
  "nbf": 0,
  "sub": "system:serviceaccount:NAMESPACE:SERVICEACCOUNT",
  "kubernetes.io": {
    "namespace": "NAMESPACE",
    "pod": {
      "name": "POD",
      "uid": ""
    },
    "serviceaccount": {
      "name": "SERVICEACCOUNT",
      "uid": ""
    },
    "warnafter": 0
  }
}
```

## GKE Clusters and their OIDC Endpoint


GKE clusters have a public OIDC discovery endpoint, even if the cluster is private. This endpoint can be used by external services to verify tokens issued by the cluster.

```
https://container.googleapis.com/v1/projects/PROJECT/locations/LOCATION/clusters/CLUSTER
```

An OIDC discovery endpoint consists of joining the "`iss`" (issuer) claim value with "`.well-known/openid-configuration`".

```json
{
  "issuer": "https://container.googleapis.com/v1/projects/PROJECT/locations/LOCATION/clusters/CLUSTER",
  "jwks_uri": "https://container.googleapis.com/v1/projects/PROJECT/locations/LOCATION/clusters/CLUSTER/jwks",
  "response_types_supported": [
    "id_token"
  ],
  "subject_types_supported": [
    "public"
  ],
  "id_token_signing_alg_values_supported": [
    "RS256"
  ],
  "claims_supported": [
    "iss",
    "sub",
    "kubernetes.io"
  ],
  "grant_types": [
    "urn:kubernetes:grant_type:programmatic_authorization"
  ]
}
```

## Demo: an authentication library with simple schema support

TBC

## Conclusion 

In conclusion, GKE OIDC workload identity for external resources is a powerful feature that enables you to securely authenticate with external resources or services. The OIDC token validation library simplifies the process of validating OIDC tokens issued by GKE, making it easier to ensure that your tokens are valid and secure. This feature is available now in Kubernetes 1.24 and can help simplify authentication in complex Kubernetes environments.

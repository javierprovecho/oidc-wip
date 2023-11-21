package parse

import (
	"gopkg.in/square/go-jose.v2/jwt"
)

type k8SClaims struct {
	jwt.Claims
	Kubernetes kubernetesClaims `json:"kubernetes.io"`
}

type kubernetesClaims struct {
	Namespace      string         `json:"namespace"`
	Pod            resourceClaims `json:"pod"`
	ServiceAccount resourceClaims `json:"serviceaccount"`
}

type resourceClaims struct {
	Name string `json:"name"`
	UID  string `json:"uid"`
}

// GetK8SSub returns namespace, serviceaccount, pod
func GetK8SSub(token string) (string, string, string, error) {
	jwtToken, err := jwt.ParseSigned(token)
	if err != nil {
		return "", "", "", err
	}

	unsafeClaims := k8SClaims{}
	if err := jwtToken.UnsafeClaimsWithoutVerification(&unsafeClaims); err != nil {
		return "", "", "", err
	}

	return unsafeClaims.Kubernetes.Namespace, unsafeClaims.Kubernetes.ServiceAccount.Name, unsafeClaims.Kubernetes.Pod.Name, nil
}

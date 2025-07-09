package dep

import (
	ezrolloutv1 "codo-cnmp/common/ezrollout/v1"
	"github.com/openkruise/kruise-api/apps/v1alpha1"
	gamekruisev1alpha1 "github.com/openkruise/kruise-game/apis/v1alpha1"
	"k8s.io/client-go/kubernetes/scheme"
)

type CRD struct {
}

func (c *CRD) Register() {
	// 注册 CRD
	_ = v1alpha1.AddToScheme(scheme.Scheme)
	_ = gamekruisev1alpha1.AddToScheme(scheme.Scheme)
	_ = ezrolloutv1.AddToScheme(scheme.Scheme)
}

func NewCRD() *CRD {
	return &CRD{}
}

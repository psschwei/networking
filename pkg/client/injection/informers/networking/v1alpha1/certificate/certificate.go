/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by injection-gen. DO NOT EDIT.

package certificate

import (
	context "context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	cache "k8s.io/client-go/tools/cache"
	apisnetworkingv1alpha1 "knative.dev/networking/pkg/apis/networking/v1alpha1"
	versioned "knative.dev/networking/pkg/client/clientset/versioned"
	v1alpha1 "knative.dev/networking/pkg/client/informers/externalversions/networking/v1alpha1"
	client "knative.dev/networking/pkg/client/injection/client"
	factory "knative.dev/networking/pkg/client/injection/informers/factory"
	networkingv1alpha1 "knative.dev/networking/pkg/client/listers/networking/v1alpha1"
	controller "knative.dev/pkg/controller"
	injection "knative.dev/pkg/injection"
	logging "knative.dev/pkg/logging"
)

func init() {
	injection.Default.RegisterInformer(withInformer)
	injection.Dynamic.RegisterDynamicInformer(withDynamicInformer)
}

// Key is used for associating the Informer inside the context.Context.
type Key struct{}

func withInformer(ctx context.Context) (context.Context, controller.Informer) {
	f := factory.Get(ctx)
	inf := f.Networking().V1alpha1().Certificates()
	return context.WithValue(ctx, Key{}, inf), inf.Informer()
}

func withDynamicInformer(ctx context.Context) context.Context {
	inf := &wrapper{client: client.Get(ctx)}
	return context.WithValue(ctx, Key{}, inf)
}

// Get extracts the typed informer from the context.
func Get(ctx context.Context) v1alpha1.CertificateInformer {
	untyped := ctx.Value(Key{})
	if untyped == nil {
		logging.FromContext(ctx).Panic(
			"Unable to fetch knative.dev/networking/pkg/client/informers/externalversions/networking/v1alpha1.CertificateInformer from context.")
	}
	return untyped.(v1alpha1.CertificateInformer)
}

type wrapper struct {
	client versioned.Interface

	namespace string
}

var _ v1alpha1.CertificateInformer = (*wrapper)(nil)
var _ networkingv1alpha1.CertificateLister = (*wrapper)(nil)

func (w *wrapper) Informer() cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(nil, &apisnetworkingv1alpha1.Certificate{}, 0, nil)
}

func (w *wrapper) Lister() networkingv1alpha1.CertificateLister {
	return w
}

func (w *wrapper) Certificates(namespace string) networkingv1alpha1.CertificateNamespaceLister {
	return &wrapper{client: w.client, namespace: namespace}
}

func (w *wrapper) List(selector labels.Selector) (ret []*apisnetworkingv1alpha1.Certificate, err error) {
	lo, err := w.client.NetworkingV1alpha1().Certificates(w.namespace).List(context.TODO(), v1.ListOptions{
		LabelSelector: selector.String(),
		// TODO(mattmoor): Incorporate resourceVersion bounds based on staleness criteria.
	})
	if err != nil {
		return nil, err
	}
	for idx := range lo.Items {
		ret = append(ret, &lo.Items[idx])
	}
	return ret, nil
}

func (w *wrapper) Get(name string) (*apisnetworkingv1alpha1.Certificate, error) {
	return w.client.NetworkingV1alpha1().Certificates(w.namespace).Get(context.TODO(), name, v1.GetOptions{
		// TODO(mattmoor): Incorporate resourceVersion bounds based on staleness criteria.
	})
}

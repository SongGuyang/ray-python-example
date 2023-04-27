/*
Copyright 2023.

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
// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1 "github.com/ray-automl/apis/automl/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeProxies implements ProxyInterface
type FakeProxies struct {
	Fake *FakeAutomlV1
	ns   string
}

var proxiesResource = corev1.SchemeGroupVersion.WithResource("proxies")

var proxiesKind = corev1.SchemeGroupVersion.WithKind("Proxy")

// Get takes name of the proxy, and returns the corresponding proxy object, and an error if there is any.
func (c *FakeProxies) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.Proxy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(proxiesResource, c.ns, name), &v1.Proxy{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.Proxy), err
}

// List takes label and field selectors, and returns the list of Proxies that match those selectors.
func (c *FakeProxies) List(ctx context.Context, opts metav1.ListOptions) (result *v1.ProxyList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(proxiesResource, proxiesKind, c.ns, opts), &v1.ProxyList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1.ProxyList{ListMeta: obj.(*v1.ProxyList).ListMeta}
	for _, item := range obj.(*v1.ProxyList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested proxies.
func (c *FakeProxies) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(proxiesResource, c.ns, opts))

}

// Create takes the representation of a proxy and creates it.  Returns the server's representation of the proxy, and an error, if there is any.
func (c *FakeProxies) Create(ctx context.Context, proxy *v1.Proxy, opts metav1.CreateOptions) (result *v1.Proxy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(proxiesResource, c.ns, proxy), &v1.Proxy{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.Proxy), err
}

// Update takes the representation of a proxy and updates it. Returns the server's representation of the proxy, and an error, if there is any.
func (c *FakeProxies) Update(ctx context.Context, proxy *v1.Proxy, opts metav1.UpdateOptions) (result *v1.Proxy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(proxiesResource, c.ns, proxy), &v1.Proxy{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.Proxy), err
}

// Delete takes name of the proxy and deletes it. Returns an error if one occurs.
func (c *FakeProxies) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(proxiesResource, c.ns, name, opts), &v1.Proxy{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeProxies) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(proxiesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1.ProxyList{})
	return err
}

// Patch applies the patch and returns the patched proxy.
func (c *FakeProxies) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Proxy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(proxiesResource, c.ns, name, pt, data, subresources...), &v1.Proxy{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.Proxy), err
}

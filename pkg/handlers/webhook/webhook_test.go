/*
Copyright 2018 Bitnami

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

package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/bitnami-labs/kubewatch/config"
	"github.com/bitnami-labs/kubewatch/pkg/event"
)

func TestWebhookInit(t *testing.T) {
	s := &Webhook{}
	expectedError := fmt.Errorf(webhookErrMsg, "Missing Webhook url")

	var Tests = []struct {
		webhook config.Webhook
		err     error
	}{
		{config.Webhook{Url: "foo"}, nil},
		{config.Webhook{}, expectedError},
	}

	for _, tt := range Tests {
		c := &config.Config{}
		c.Handler.Webhook = tt.webhook
		if err := s.Init(c); !reflect.DeepEqual(err, tt.err) {
			t.Fatalf("Init(): %v", err)
		}
	}
}

func checkHeaderAndHMAC(hmacKey []byte, headerName string, r *http.Request) (headerOk, hmacSigOk bool) {
	headerOk = r.Header.Get(headerName) != ""
	hmacValue := r.Header.Get(headerName)

	if hmacValue != "" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		_ = r.Body.Close()

		mac := hmac.New(sha256.New, hmacKey)
		mac.Write(body)

		hmacSigOk = hex.EncodeToString(mac.Sum(nil)) == hmacValue
	}

	return
}

func TestWebhook_Handle(t *testing.T) {
	ev := event.New(event.Event{
		Namespace: "default",
		Kind:      "pod",
		Component: "",
		Host:      "some-node",
		Reason:    "created",
		Status:    "Normal",
		Name:      "cool-pod",
	}, "created")

	t.Run("default header", func(t *testing.T) {
		headerOk := false
		hmacSigOk := false

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headerOk, hmacSigOk = checkHeaderAndHMAC([]byte(`123`), "X-KubeWatch-Signature", r)
			_, _ = w.Write([]byte(`ok`))
		}))
		defer srv.Close()

		s := &Webhook{}
		c := &config.Config{}
		c.Handler.Webhook = config.Webhook{
			Url:     srv.URL,
			HMACKey: "MTIz", // 123
		}

		if err := s.Init(c); err != nil {
			t.Fatal(err)
		}

		s.Handle(ev)

		if !headerOk {
			t.Fatal("header does not match")
		}
		if !hmacSigOk {
			t.Fatal("hmac signature does not match")
		}
	})

	t.Run("custom header", func(t *testing.T) {
		headerOk := false
		hmacSigOk := false

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headerOk, hmacSigOk = checkHeaderAndHMAC([]byte(`123`), "X-Custom-Header", r)
			_, _ = w.Write([]byte(`ok`))
		}))
		defer srv.Close()

		s := &Webhook{}
		c := &config.Config{}
		c.Handler.Webhook = config.Webhook{
			Url:                 srv.URL,
			HMACKey:             "MTIz", // 123
			HMACSignatureHeader: "X-Custom-Header",
		}

		if err := s.Init(c); err != nil {
			t.Fatal(err)
		}

		s.Handle(ev)

		if !headerOk {
			t.Fatal("header does not match")
		}
		if !hmacSigOk {
			t.Fatal("hmac signature does not match")
		}
	})
}

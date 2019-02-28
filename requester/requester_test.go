// Copyright 2014 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package requester

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestN(t *testing.T) {
	var count int64
	handler := func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&count, int64(1))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()
    //urls := []string{"https://api.dongdakid.com/v2/reservableitems"}
    //oncework := "classes"
    //oncework := "teachers"
    //oncework := "sessioninfos"  //跑了7，8分钟一遍没跑完  跑分页
	oncework := "sessioninfos?page[number]=1&page[size]=20"   //sessioninfos的分页查询
	//oncework := "units"     //有问题不稳定
	//oncework := "students"   //跑了7，8分钟一遍没跑完   跑分页
	//oncework := "students?page[number]=1&page[size]=20"   //students的分页查询
	//oncework := "brands"
	//oncework := "categories"
	//oncework := "yards"
	//oncework := "rooms"
	//oncework := "kids"
	//oncework := "images"
	//oncework := "applies"     //跑了7，8分钟一遍没跑完   跑分页
	//oncework := "applies?page[number]=1&page[size]=20"   //applies的分页查询
	//oncework := "reservableitems"
	//oncework := "Catenodes"
	for i:=0;i<10;i++ {
		req, _ := http.NewRequest("GET","https://api.dongdakid.com/v2/" + oncework , nil)
		header := make(http.Header)
		header.Add("Authorization", "bearer 46fa418a3ec8ec32cd8662a589f3b403")
		req.Header = header
		w := &Work{
			Request: req,
			N:       200,
			C:       100,
		}
		w.Run()
		fmt.Println(i,"--------------------------------------------------------------------")
		if i==6{
			fmt.Println(i)
		}
	}
	if count != 20 {
		t.Errorf("Expected to send 20 requests, found %v", count)
	}
}

func TestQps(t *testing.T) {
	var wg sync.WaitGroup
	var count int64
	handler := func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&count, int64(1))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL, nil)
	header := make(http.Header)
	header.Add("Authorization", "bearer 46fa418a3ec8ec32cd8662a589f3b403")
	req.Header = header
	w := &Work{
		Request: req,
		N:       20,
		C:       2,
		QPS:     1,
	}
	wg.Add(1)
	time.AfterFunc(time.Second, func() {
		if count > 2 {
			t.Errorf("Expected to work at most 2 times, found %v", count)
		}
		wg.Done()
	})
	go w.Run()
	wg.Wait()
}

/*func TestRequest(t *testing.T) {
	var uri, contentType, some, method, auth string
	handler := func(w http.ResponseWriter, r *http.Request) {
		uri = r.RequestURI
		method = r.Method
		contentType = r.Header.Get("Content-type")
		some = r.Header.Get("X-some")
		auth = r.Header.Get("Authorization")
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	header := make(http.Header)
	header.Add("Content-type", "text/html")
	header.Add("X-some", "value")
	req, _ := http.NewRequest("GET", server.URL, nil)
	req.Header = header
	req.SetBasicAuth("username", "password")
	w := &Work{
		Request: req,
		N:       1,
		C:       1,
	}
	w.Run()
	if uri != "/" {
		t.Errorf("Uri is expected to be /, %v is found", uri)
	}
	if contentType != "text/html" {
		t.Errorf("Content type is expected to be text/html, %v is found", contentType)
	}
	if some != "value" {
		t.Errorf("X-some header is expected to be value, %v is found", some)
	}
	if auth != "Basic dXNlcm5hbWU6cGFzc3dvcmQ=" {
		t.Errorf("Basic authorization is not properly set")
	}
}
*/
func TestBody(t *testing.T) {
	var count int64
	handler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		if string(body) == "Body" {
			atomic.AddInt64(&count, 1)
		}
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	req, _ := http.NewRequest("POST", server.URL, bytes.NewBuffer([]byte("Body")))
	w := &Work{
		Request:     req,
		RequestBody: []byte("Body"),
		N:           10,
		C:           1,
	}
	w.Run()
	if count != 10 {
		t.Errorf("Expected to work 10 times, found %v", count)
	}
}

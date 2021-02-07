package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type teststruct struct {
	ID string `json:"id" validate:"required"`
}

func TestRespondJSON(t *testing.T) {

	/*
		ts := httptest.NewServer(func(w http.ResponseWriter, r *http.Request) {
			respondJSON(w, http.StatusOK, teststruct{ID: "asid"})
		})
		defer ts.Close()
	*/

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	status := http.StatusOK
	if resp.Result().StatusCode != status {
		t.Fatalf("status")
	}

	handle := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, status, teststruct{ID: "test"})
	})

	handle.ServeHTTP(resp, req)
	hh := resp.Header()
	if hh.Get("Content-Type") != "application/json" {
		t.Fatalf("bad content type")
	}

	if p, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Fatalf("ioutil.ReadAll: %s", err.Error())
	} else {
		obj := &teststruct{}

		err := json.Unmarshal(p, obj)
		if err != nil {
			t.Fatalf("json.Unmarshal: %s", err.Error())
		} else {
			t.Logf("%v", obj)
			if obj.ID != "test" {
				t.Fatalf("unwanted value, bad test...")
			}
		}
	}
}

func TestRespondError(t *testing.T) {
	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	status := http.StatusForbidden

	x := errors.New("test error")
	handle := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respondError(w, status, x)
	})

	handle.ServeHTTP(resp, req)
	hh := resp.Header()
	if hh.Get("Content-Type") != "application/json" {
		t.Fatalf("bad content type")
	}

	if resp.Result().StatusCode != status {
		t.Fatalf("status")
	}

	if p, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Fatalf("ioutil.ReadAll: %s", err.Error())
	} else {
		obj := map[string]string{}

		err := json.Unmarshal(p, &obj)
		if err != nil {
			t.Fatalf("json.Unmarshal: %s", err.Error())
		} else {
			t.Logf("%v", obj)
			if obj["error"] != "test error" {
				t.Fatalf("unwanted value, bad test...")
			}
		}
	}
}

func TestBindJSON(t *testing.T) {
	temp := teststruct{
		ID: "test 2",
	}
	body, err := json.Marshal(temp)
	if err != nil {
		t.Fatalf("json.Marshal: %s", err.Error())
	}

	resp := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	status := http.StatusOK

	handle := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		x := teststruct{}
		err := bindJSON(r, &x)
		if err != nil {
			respondError(w, http.StatusBadRequest, errors.New("input is not json"))
			return
		}

		respondJSON(w, status, x)
	})

	handle.ServeHTTP(resp, req)
	hh := resp.Header()
	if hh.Get("Content-Type") != "application/json" {
		t.Fatalf("bad content type")
	}

	if resp.Result().StatusCode != status {
		t.Fatalf("status")
	}

	if p, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Fatalf("ioutil.ReadAll: %s", err.Error())
	} else {
		t.Logf("%s", string(p))
		obj := teststruct{}

		err := json.Unmarshal(p, &obj)
		if err != nil {
			t.Fatalf("json.Unmarshal: %s", err.Error())
		} else {
			t.Logf("%v", obj)
			if obj.ID != "test 2" {
				t.Fatalf("unwanted value, bad test...")
			}
		}
	}

}

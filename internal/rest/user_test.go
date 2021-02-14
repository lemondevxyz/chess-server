package rest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	us = User{}
)

func TestNewUser(t *testing.T) {
	us = AddClient(cl1)
}

func TestGetUser(t *testing.T) {
	resp := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", us.Token))

	handle := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := GetUser(r)
		if err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		respondJSON(w, http.StatusOK, "success")
	})

	handle.ServeHTTP(resp, req)
	if resp.Result().StatusCode != http.StatusOK {
		t.Fatalf("header authentication doesnt work")
	}

}

func TestUserClient(t *testing.T) {
	cl := us.Client()
	if cl1 != cl {
		t.Fatalf("not the same pointers")
	}
}

func TestUserDelete(t *testing.T) {
	us.Delete()
	if us.Client() != nil {
		t.Fatalf("delete does not delete")
	}
}

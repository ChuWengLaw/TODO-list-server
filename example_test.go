package main

import (
	"net/http"
	"net/http/httptest"
	todo "server/utils/todo"
	"testing"
)

func TestList(t *testing.T) {
	var ww http.ResponseWriter
	var rr *http.Request
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww = w
		rr = r
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	client := mockServer.Client()
	resp, err := client.Get(mockServer.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}

	defer func() {
		if err := recover(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}()

	todo.List(ww, rr)
}

func TestAdd(t *testing.T) {
	var ww http.ResponseWriter
	var rr *http.Request
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww = w
		rr = r
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	client := mockServer.Client()
	resp, err := client.Get(mockServer.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}

	defer func() {
		if err := recover(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}()

	todo.Add(ww, rr)
}

func TestMark(t *testing.T) {
	var ww http.ResponseWriter
	var rr *http.Request
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww = w
		rr = r
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	client := mockServer.Client()
	resp, err := client.Get(mockServer.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}

	defer func() {
		if err := recover(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}()

	todo.Mark(ww, rr)
}

func TestDelete(t *testing.T) {
	var ww http.ResponseWriter
	var rr *http.Request
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww = w
		rr = r
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	client := mockServer.Client()
	resp, err := client.Get(mockServer.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}

	defer func() {
		if err := recover(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}()

	todo.Delete(ww, rr)
}

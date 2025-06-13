package main

import "testing"

func TestCreateQueryParamsWithoutExclusions(t *testing.T) {
	query := CreateQueryParams("example.com", []string{})
	if query != "site:example.com" {
		t.Errorf("Expected 'site:example.com', got '%s'", query)
	}
}

func TestCreateQueryParamsWithExclusions(t *testing.T) {
	query := CreateQueryParams("example.com", []string{"www.example.com", "wiki.example.com"})
	if query != "site:example.com -site:www.example.com -site:wiki.example.com" {
		t.Errorf("Expected 'site:example.com -site:www.example.com -site:wiki.example.com', got '%s'", query)
	}
}

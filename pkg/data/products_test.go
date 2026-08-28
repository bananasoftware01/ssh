package data

import (
	"strings"
	"testing"
)

func TestProductsData(t *testing.T) {
	if len(Products) == 0 {
		t.Fatal("expected at least one product")
	}

	for i, p := range Products {
		if p.Name == "" {
			t.Errorf("product %d has empty Name", i)
		}
		if p.URL == "" || !strings.HasPrefix(p.URL, "https://") {
			t.Errorf("product %s has invalid URL: %s", p.Name, p.URL)
		}
		if p.DescVI == "" || p.DescEN == "" {
			t.Errorf("product %s missing descriptions", p.Name)
		}
	}
}

func TestCompanyData(t *testing.T) {
	if Company.TaxID != "2803238388" {
		t.Errorf("expected TaxID 2803238388, got %s", Company.TaxID)
	}
	if Company.Phone != "0335581402" {
		t.Errorf("expected Phone 0335581402, got %s", Company.Phone)
	}
}

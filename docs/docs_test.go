package docs

import (
	"strings"
	"testing"

	"github.com/swaggo/swag"
)

func TestSwaggerInfo_BasicMetadata(t *testing.T) {
	if SwaggerInfo == nil {
		t.Fatalf("SwaggerInfo is nil")
	}
	if SwaggerInfo.InstanceName() != "swagger" {
		t.Fatalf("SwaggerInfo.InstanceName() = %q, want %q", SwaggerInfo.InstanceName(), "swagger")
	}
	if SwaggerInfo.Title == "" {
		t.Fatalf("SwaggerInfo.Title is empty")
	}
	if SwaggerInfo.Version == "" {
		t.Fatalf("SwaggerInfo.Version is empty")
	}
	if SwaggerInfo.BasePath == "" {
		t.Fatalf("SwaggerInfo.BasePath is empty")
	}
}

func TestSwaggerDocument_IsRegisteredAndContainsExpectedSections(t *testing.T) {
	doc, err := swag.ReadDoc(SwaggerInfo.InstanceName())
	if err != nil {
		t.Fatalf("swag.ReadDoc() error = %v", err)
	}
	if doc == "" {
		t.Fatalf("swag.ReadDoc() returned empty document")
	}

	expectedFragments := []string{
		`"swagger": "2.0"`,
		`"/products"`,
		`"/users/login"`,
		`"securityDefinitions"`,
	}

	for _, fragment := range expectedFragments {
		if !strings.Contains(doc, fragment) {
			t.Fatalf("swagger doc does not contain required fragment: %s", fragment)
		}
	}
}

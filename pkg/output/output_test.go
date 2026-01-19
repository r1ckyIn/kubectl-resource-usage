package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/r1ckyIn/kubectl-resource-usage/pkg/calculator"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestTableFormatter(t *testing.T) {
	podUsages := []calculator.PodUsage{
		{
			Namespace: "default",
			Name:      "test-pod",
			Node:      "node-1",
			CPU: calculator.ResourceUsage{
				Usage:          resource.MustParse("100m"),
				Requests:       resourcePtr(resource.MustParse("200m")),
				Limits:         resourcePtr(resource.MustParse("500m")),
				RequestPercent: intPtr(50),
				LimitPercent:   intPtr(20),
			},
			Memory: calculator.ResourceUsage{
				Usage:          resource.MustParse("128Mi"),
				Requests:       resourcePtr(resource.MustParse("256Mi")),
				Limits:         resourcePtr(resource.MustParse("512Mi")),
				RequestPercent: intPtr(50),
				LimitPercent:   intPtr(25),
			},
		},
	}

	var buf bytes.Buffer
	formatter := &TableFormatter{}
	err := formatter.Format(&buf, podUsages)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	// Check header
	if !strings.Contains(output, "NAMESPACE") {
		t.Error("expected header to contain NAMESPACE")
	}
	if !strings.Contains(output, "CPU_REQ%") {
		t.Error("expected header to contain CPU_REQ%")
	}

	// Check data
	if !strings.Contains(output, "default") {
		t.Error("expected output to contain namespace 'default'")
	}
	if !strings.Contains(output, "test-pod") {
		t.Error("expected output to contain pod name 'test-pod'")
	}
	if !strings.Contains(output, "50%") {
		t.Error("expected output to contain '50%'")
	}
}

func TestTableFormatterWithNA(t *testing.T) {
	podUsages := []calculator.PodUsage{
		{
			Namespace: "default",
			Name:      "test-pod",
			Node:      "node-1",
			CPU: calculator.ResourceUsage{
				Usage:          resource.MustParse("100m"),
				RequestPercent: nil,
				LimitPercent:   nil,
			},
			Memory: calculator.ResourceUsage{
				Usage:          resource.MustParse("128Mi"),
				RequestPercent: nil,
				LimitPercent:   nil,
			},
		},
	}

	var buf bytes.Buffer
	formatter := &TableFormatter{}
	err := formatter.Format(&buf, podUsages)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "N/A") {
		t.Error("expected output to contain 'N/A' for nil percentages")
	}
}

func TestJSONFormatter(t *testing.T) {
	podUsages := []calculator.PodUsage{
		{
			Namespace: "default",
			Name:      "test-pod",
			Node:      "node-1",
			CPU: calculator.ResourceUsage{
				Usage:          resource.MustParse("100m"),
				Requests:       resourcePtr(resource.MustParse("200m")),
				Limits:         resourcePtr(resource.MustParse("500m")),
				RequestPercent: intPtr(50),
				LimitPercent:   intPtr(20),
			},
			Memory: calculator.ResourceUsage{
				Usage:          resource.MustParse("128Mi"),
				Requests:       resourcePtr(resource.MustParse("256Mi")),
				Limits:         resourcePtr(resource.MustParse("512Mi")),
				RequestPercent: intPtr(50),
				LimitPercent:   intPtr(25),
			},
		},
	}

	var buf bytes.Buffer
	formatter := &JSONFormatter{}
	err := formatter.Format(&buf, podUsages)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify it's valid JSON
	var result StructuredOutput
	err = json.Unmarshal(buf.Bytes(), &result)
	if err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	// Verify content
	if len(result.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(result.Items))
	}

	item := result.Items[0]
	if item.Namespace != "default" {
		t.Errorf("expected namespace 'default', got '%s'", item.Namespace)
	}
	if item.Pod != "test-pod" {
		t.Errorf("expected pod 'test-pod', got '%s'", item.Pod)
	}
	if item.CPU.RequestPercent == nil || *item.CPU.RequestPercent != 50 {
		t.Errorf("expected CPU request percent 50, got %v", item.CPU.RequestPercent)
	}
}

func TestJSONFormatterWithNil(t *testing.T) {
	podUsages := []calculator.PodUsage{
		{
			Namespace: "default",
			Name:      "test-pod",
			Node:      "node-1",
			CPU: calculator.ResourceUsage{
				Usage:          resource.MustParse("100m"),
				RequestPercent: nil,
				LimitPercent:   nil,
			},
			Memory: calculator.ResourceUsage{
				Usage:          resource.MustParse("128Mi"),
				RequestPercent: nil,
				LimitPercent:   nil,
			},
		},
	}

	var buf bytes.Buffer
	formatter := &JSONFormatter{}
	err := formatter.Format(&buf, podUsages)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify it's valid JSON with null values
	var result StructuredOutput
	err = json.Unmarshal(buf.Bytes(), &result)
	if err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	item := result.Items[0]
	if item.CPU.RequestPercent != nil {
		t.Errorf("expected nil CPU request percent, got %v", item.CPU.RequestPercent)
	}
	if item.CPU.Requests != nil {
		t.Errorf("expected nil CPU requests, got %v", item.CPU.Requests)
	}
}

func TestYAMLFormatter(t *testing.T) {
	podUsages := []calculator.PodUsage{
		{
			Namespace: "default",
			Name:      "test-pod",
			Node:      "node-1",
			CPU: calculator.ResourceUsage{
				Usage:          resource.MustParse("100m"),
				Requests:       resourcePtr(resource.MustParse("200m")),
				Limits:         resourcePtr(resource.MustParse("500m")),
				RequestPercent: intPtr(50),
				LimitPercent:   intPtr(20),
			},
			Memory: calculator.ResourceUsage{
				Usage:          resource.MustParse("128Mi"),
				Requests:       resourcePtr(resource.MustParse("256Mi")),
				Limits:         resourcePtr(resource.MustParse("512Mi")),
				RequestPercent: intPtr(50),
				LimitPercent:   intPtr(25),
			},
		},
	}

	var buf bytes.Buffer
	formatter := &YAMLFormatter{}
	err := formatter.Format(&buf, podUsages)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	// Check YAML content
	if !strings.Contains(output, "namespace: default") {
		t.Error("expected output to contain 'namespace: default'")
	}
	if !strings.Contains(output, "pod: test-pod") {
		t.Error("expected output to contain 'pod: test-pod'")
	}
	if !strings.Contains(output, "requestPercent: 50") {
		t.Error("expected output to contain 'requestPercent: 50'")
	}
}

func TestWideFormatter(t *testing.T) {
	podUsages := []calculator.PodUsage{
		{
			Namespace: "default",
			Name:      "test-pod",
			Node:      "node-1",
			CPU: calculator.ResourceUsage{
				Usage:          resource.MustParse("100m"),
				Requests:       resourcePtr(resource.MustParse("200m")),
				Limits:         resourcePtr(resource.MustParse("500m")),
				RequestPercent: intPtr(50),
				LimitPercent:   intPtr(20),
			},
			Memory: calculator.ResourceUsage{
				Usage:          resource.MustParse("128Mi"),
				Requests:       resourcePtr(resource.MustParse("256Mi")),
				Limits:         resourcePtr(resource.MustParse("512Mi")),
				RequestPercent: intPtr(50),
				LimitPercent:   intPtr(25),
			},
		},
	}

	var buf bytes.Buffer
	formatter := &WideFormatter{}
	err := formatter.Format(&buf, podUsages)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()

	// Check header includes raw value columns
	if !strings.Contains(output, "CPU_REQ") {
		t.Error("expected header to contain CPU_REQ")
	}
	if !strings.Contains(output, "CPU_LIM") {
		t.Error("expected header to contain CPU_LIM")
	}
	if !strings.Contains(output, "MEM_REQ") {
		t.Error("expected header to contain MEM_REQ")
	}
	if !strings.Contains(output, "MEM_LIM") {
		t.Error("expected header to contain MEM_LIM")
	}

	// Check data includes raw values
	if !strings.Contains(output, "200m") {
		t.Error("expected output to contain CPU request '200m'")
	}
	if !strings.Contains(output, "500m") {
		t.Error("expected output to contain CPU limit '500m'")
	}
	if !strings.Contains(output, "256Mi") {
		t.Error("expected output to contain memory request '256Mi'")
	}
	if !strings.Contains(output, "512Mi") {
		t.Error("expected output to contain memory limit '512Mi'")
	}
}

func TestWideFormatterWithNA(t *testing.T) {
	podUsages := []calculator.PodUsage{
		{
			Namespace: "default",
			Name:      "test-pod",
			Node:      "node-1",
			CPU: calculator.ResourceUsage{
				Usage:          resource.MustParse("100m"),
				Requests:       nil,
				Limits:         nil,
				RequestPercent: nil,
				LimitPercent:   nil,
			},
			Memory: calculator.ResourceUsage{
				Usage:          resource.MustParse("128Mi"),
				Requests:       nil,
				Limits:         nil,
				RequestPercent: nil,
				LimitPercent:   nil,
			},
		},
	}

	var buf bytes.Buffer
	formatter := &WideFormatter{}
	err := formatter.Format(&buf, podUsages)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "N/A") {
		t.Error("expected output to contain 'N/A' for nil values")
	}
}

func TestNewFormatter(t *testing.T) {
	// Test table formatter (default)
	tableFormatter := NewFormatter("table")
	if _, ok := tableFormatter.(*TableFormatter); !ok {
		t.Error("expected TableFormatter for 'table' format")
	}

	// Test JSON formatter
	jsonFormatter := NewFormatter("json")
	if _, ok := jsonFormatter.(*JSONFormatter); !ok {
		t.Error("expected JSONFormatter for 'json' format")
	}

	// Test YAML formatter
	yamlFormatter := NewFormatter("yaml")
	if _, ok := yamlFormatter.(*YAMLFormatter); !ok {
		t.Error("expected YAMLFormatter for 'yaml' format")
	}

	// Test Wide formatter
	wideFormatter := NewFormatter("wide")
	if _, ok := wideFormatter.(*WideFormatter); !ok {
		t.Error("expected WideFormatter for 'wide' format")
	}

	// Test unknown format defaults to table
	unknownFormatter := NewFormatter("unknown")
	if _, ok := unknownFormatter.(*TableFormatter); !ok {
		t.Error("expected TableFormatter for unknown format")
	}
}

func intPtr(i int) *int {
	return &i
}

func resourcePtr(r resource.Quantity) *resource.Quantity {
	return &r
}

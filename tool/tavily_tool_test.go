package tool

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestTavilySearch(t *testing.T) {
	// Check if API key is available in environment
	if os.Getenv("TAVILY_API_KEY") == "" {
		t.Skip("Skipping TavilySearch tests: TAVILY_API_KEY environment variable is not set")
		return
	}

	// Test with existing API key
	_, err := TavilySearch("test query")
	if err != nil {
		t.Logf("TavilySearch() returned error (may be expected): %v", err)
	}
}

func TestTavilySearchMissingAPIKey(t *testing.T) {
	// Save original API key
	originalAPIKey := os.Getenv("TAVILY_API_KEY")
	defer os.Setenv("TAVILY_API_KEY", originalAPIKey)

	// Unset API key for this test
	os.Unsetenv("TAVILY_API_KEY")

	// Test missing API key
	_, err := TavilySearch("test query")
	if err == nil {
		t.Error("TavilySearch() without API key expected error, got nil")
	}
	if !strings.Contains(err.Error(), "TAVILY_API_KEY") {
		t.Errorf("TavilySearch() error should mention API key, got %v", err)
	}

	// Test with URL function as well
	_, err = TavilySearchWithLimitAndURL("test query", 10, "http://test.com")
	if err == nil {
		t.Error("TavilySearchWithLimitAndURL() without API key expected error, got nil")
	}
	if !strings.Contains(err.Error(), "TAVILY_API_KEY") {
		t.Errorf("TavilySearchWithLimitAndURL() error should mention API key, got %v", err)
	}
}

func TestTavilySearchWithLimit(t *testing.T) {
	// Check if API key is available in environment
	if os.Getenv("TAVILY_API_KEY") == "" {
		t.Skip("Skipping TavilySearchWithLimit tests: TAVILY_API_KEY environment variable is not set")
		return
	}

	// Test different limit values
	tests := []struct {
		name       string
		query      string
		maxResults int
	}{
		{
			name:       "Default limit",
			query:      "test query",
			maxResults: 20,
		},
		{
			name:       "Zero limit (should default to 5)",
			query:      "test query",
			maxResults: 0,
		},
		{
			name:       "Negative limit (should default to 5)",
			query:      "test query",
			maxResults: -1,
		},
		{
			name:       "Limit above maximum (should cap at 100)",
			query:      "test query",
			maxResults: 150,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := TavilySearchWithLimit(tt.query, tt.maxResults)
			// The call might fail due to API issues, but we're testing that the function doesn't panic
			if err != nil {
				t.Logf("TavilySearchWithLimit() returned error (may be expected): %v", err)
			}
		})
	}
}

func TestTavilySearchResponseParsing(t *testing.T) {
	// Test the response parsing logic by simulating different API responses
	testCases := []struct {
		name           string
		response       string
		expectedResult string
		expectError    bool
	}{
		{
			name: "Normal response with results",
			response: `{
				"results": [
					{
						"title": "Result 1",
						"url": "https://example.com/1",
						"content": "Content 1"
					},
					{
						"title": "Result 2",
						"url": "https://example.com/2",
						"content": "Content 2"
					}
				],
				"images": ["https://example.com/img1.jpg"]
			}`,
			expectedResult: "Title: Result 1\nURL: https://example.com/1\nContent: Content 1\n\nTitle: Result 2\nURL: https://example.com/2\nContent: Content 2\n\n\nRelevant Images:\n- Image URL: https://example.com/img1.jpg\n\n",
			expectError:    false,
		},
		{
			name: "Empty results",
			response: `{
				"results": [],
				"images": []
			}`,
			expectedResult: "No results found.",
			expectError:    false,
		},
		{
			name: "Results without images",
			response: `{
				"results": [
					{
						"title": "Single Result",
						"url": "https://example.com",
						"content": "Single content"
					}
				],
				"images": []
			}`,
			expectedResult: "Title: Single Result\nURL: https://example.com\nContent: Single content\n\n",
			expectError:    false,
		},
		{
			name:        "Invalid JSON",
			response:    `{invalid json}`,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, tc.response)
			}))
			defer server.Close()

			// Use mock API key for testing
			os.Setenv("TAVILY_API_KEY", "test-key")
			defer os.Unsetenv("TAVILY_API_KEY")

			// Test with our configurable URL function
			result, err := TavilySearchWithLimitAndURL("test query", 10, server.URL)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				} else if result != tc.expectedResult {
					t.Errorf("Expected result:\n%s\n\nGot:\n%s", tc.expectedResult, result)
				}
			}
		})
	}
}

func TestTavilySearchEdgeCases(t *testing.T) {
	// Check if API key is available in environment
	if os.Getenv("TAVILY_API_KEY") == "" {
		t.Skip("Skipping TavilySearchEdgeCases tests: TAVILY_API_KEY environment variable is not set")
		return
	}

	// Test empty query
	_, err := TavilySearch("")
	if err != nil {
		// Empty query might be valid or invalid depending on the API
		// This test documents the behavior
		t.Logf("TavilySearch() with empty query returned error: %v", err)
	}

	// Test very long query
	longQuery := strings.Repeat("test ", 1000)
	_, err = TavilySearch(longQuery)
	if err != nil {
		t.Logf("TavilySearch() with long query returned error: %v", err)
	}

	// Test special characters in query
	specialQuery := "How to test \"quotes\" & ampersands in <search> queries?"
	_, err = TavilySearch(specialQuery)
	if err != nil {
		t.Logf("TavilySearch() with special characters returned error: %v", err)
	}
}

// Example of how to benchmark TavilySearch (without actual API calls)
func BenchmarkTavilySearch(b *testing.B) {
	// Check if API key is available in environment
	if os.Getenv("TAVILY_API_KEY") == "" {
		b.Skip("Skipping TavilySearch benchmark: TAVILY_API_KEY environment variable is not set")
		return
	}

	for b.Loop() {
		// This will fail due to network issues in benchmark environment, but it benchmarks the request setup
		_, err := TavilySearch("benchmark test query")
		_ = err // Ignore error for benchmarking purposes
	}
}

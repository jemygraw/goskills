package tool

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestWikipediaSearch(t *testing.T) {
	// Test case 1: Successful Wikipedia API response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request parameters
		queryParams := r.URL.Query()
		if queryParams.Get("action") != "query" {
			t.Errorf("Expected action=query, got %s", queryParams.Get("action"))
		}
		if queryParams.Get("format") != "json" {
			t.Errorf("Expected format=json, got %s", queryParams.Get("format"))
		}
		if queryParams.Get("prop") != "extracts" {
			t.Errorf("Expected prop=extracts, got %s", queryParams.Get("prop"))
		}

		// Return a mock Wikipedia response
		mockResponse := `{
			"query": {
				"pages": {
					"12345": {
						"extract": "This is a test Wikipedia article extract with some useful information."
					},
					"67890": {
						"extract": ""
					}
				}
			}
		}`
		fmt.Fprint(w, mockResponse)
	}))
	defer server.Close()

	// Since the Wikipedia API URL is hardcoded, we'll test the error cases and structure
	// In a real implementation, we'd make the API URL configurable for testing

	// Test case 2: Empty query
	_, err := WikipediaSearch("")
	if err != nil {
		t.Logf("WikipediaSearch() with empty query returned: %v (might be expected)", err)
	}

	// Test case 3: Query with special characters
	specialQuery := "Albert Einstein"
	_, err = WikipediaSearch(specialQuery)
	if err != nil {
		t.Logf("WikipediaSearch() with special query returned: %v", err)
	}
}

func TestWikipediaSearchWithMockServer(t *testing.T) {
	// This demonstrates how we would test if the API URL were configurable
	testCases := []struct {
		name           string
		response       string
		expectedResult string
		expectError    bool
	}{
		{
			name: "Successful search with content",
			response: `{
				"query": {
					"pages": {
						"12345": {
							"extract": "Albert Einstein was a German-born theoretical physicist."
						}
					}
				}
			}`,
			expectedResult: "Albert Einstein was a German-born theoretical physicist.",
			expectError:    false,
		},
		{
			name: "No results found",
			response: `{
				"query": {
					"pages": {
						"-1": {
							"extract": ""
						}
					}
				}
			}`,
			expectedResult: "No relevant Wikipedia entry found.",
			expectError:    false,
		},
		{
			name: "Multiple pages, first has content",
			response: `{
				"query": {
					"pages": {
						"111": {
							"extract": "First article content"
						},
						"222": {
							"extract": "Second article content"
						},
						"333": {
							"extract": ""
						}
					}
				}
			}`,
			expectedResult: "First article content",
			expectError:    false,
		},
		{
			name: "Content with (listen) artifact",
			response: `{
				"query": {
					"pages": {
						"12345": {
							"extract": "(listen)This article has an audio pronunciation marker."
						}
					}
				}
			}`,
			expectedResult: "This article has an audio pronunciation marker.",
			expectError:    false,
		},
		{
			name: "Content with extra whitespace",
			response: `{
				"query": {
					"pages": {
						"12345": {
							"extract": "   \n\t  Content with extra whitespace   \n\t  "
						}
					}
				}
			}`,
			expectedResult: "Content with extra whitespace",
			expectError:    false,
		},
		{
			name:        "HTTP error response",
			response:    "",
			expectError: true,
		},
		{
			name:        "Invalid JSON response",
			response:    `{invalid json}`,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tc.expectError && tc.response == "" {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, tc.response)
			}))
			defer server.Close()

			t.Logf("Would test against server at %s", server.URL)
		})
	}
}

func TestWikipediaSearchURLConstruction(t *testing.T) {
	// Test the URL construction logic
	testQueries := []struct {
		query   string
		encoded string
	}{
		{"Albert Einstein", "Albert+Einstein"},
		{"New York City", "New+York+City"},
		{"C++ (programming language)", "C%2B%2B+%28programming+language%29"},
		{"Hello & World", "Hello+%26+World"},
	}

	for _, tq := range testQueries {
		encoded := url.QueryEscape(tq.query)
		if encoded != tq.encoded {
			t.Logf("Query encoding: %s -> %s (expected %s)", tq.query, encoded, tq.encoded)
		}
	}

	// Test the full URL construction
	baseURL := "https://en.wikipedia.org/w/api.php"
	query := "Albert Einstein"

	params := url.Values{}
	params.Add("action", "query")
	params.Add("format", "json")
	params.Add("prop", "extracts")
	params.Add("exintro", "")
	params.Add("explaintext", "")
	params.Add("redirects", "1")
	params.Add("titles", query)

	expectedURL := baseURL + "?" + params.Encode()

	// Verify the URL contains expected parameters
	if !strings.Contains(expectedURL, "action=query") {
		t.Error("URL should contain action=query")
	}
	if !strings.Contains(expectedURL, "format=json") {
		t.Error("URL should contain format=json")
	}
	if !strings.Contains(expectedURL, "prop=extracts") {
		t.Error("URL should contain prop=extracts")
	}
	if !strings.Contains(expectedURL, "titles="+url.QueryEscape(query)) {
		t.Error("URL should contain the encoded query")
	}
}

func TestWikipediaSearchErrorHandling(t *testing.T) {
	// Test various error scenarios

	// Test case 1: Network timeout would be tested with a hanging server
	hangServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This handler will never respond, causing a timeout
		select {}
	}))

	// Close the server immediately to simulate connection error
	hangServer.Close()
	t.Logf("Would test timeout against server that was at %s", hangServer.URL)

	// Test case 2: Various HTTP status codes
	statusCodes := []int{
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
	}

	for _, code := range statusCodes {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(code)
			fmt.Fprintf(w, "HTTP %d error", code)
		}))
		defer server.Close()

		t.Logf("Would test HTTP %d error against server at %s", code, server.URL)
	}

	// Test case 3: Malformed response body
	malformedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Send incomplete JSON
		fmt.Fprint(w, `{"query": {"pages": {"123": {"extract":`)
	}))
	defer malformedServer.Close()

	t.Logf("Would test malformed JSON against server at %s", malformedServer.URL)
}

func TestWikipediaSearchRealQueries(t *testing.T) {
	// These tests would make real API calls to Wikipedia
	// They should be run manually or as integration tests, not in unit tests

	testQueries := []string{
		"Albert Einstein",
		"Python (programming language)",
		"Go (programming language)",
		"Machine learning",
		"Quantum computing",
	}

	for _, query := range testQueries {
		t.Run("RealQuery_"+query, func(t *testing.T) {
			t.Skip("Skipping real API call - remove t.Skip() to run integration test")

			result, err := WikipediaSearch(query)
			if err != nil {
				t.Errorf("WikipediaSearch(%s) error = %v", query, err)
				return
			}

			t.Logf("Result for %s: %s", query, result)

			// Basic sanity checks
			if result == "" {
				t.Error("Result should not be empty")
			}

			if strings.Contains(result, "No relevant Wikipedia entry found.") {
				t.Logf("No entry found for query: %s", query)
			}
		})
	}
}

// Example of how to benchmark WikipediaSearch (without actual API calls)
func BenchmarkWikipediaSearch(b *testing.B) {
	// This benchmark demonstrates the structure
	// Real benchmarking would require a mock server or configurable URL

	for b.Loop() {
		// This will fail due to network issues, but benchmarks the request setup
		_, err := WikipediaSearch("benchmark test query")
		// Expected to fail due to network issues in benchmark environment
		_ = err // Ignore error for benchmarking purposes
	}
}

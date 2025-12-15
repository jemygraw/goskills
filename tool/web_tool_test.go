package tool

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWebFetch(t *testing.T) {
	// Test case 1: Valid HTML response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		html := `<!DOCTYPE html>
<html>
<head>
	<title>Test Page</title>
	<script>console.log('test');</script>
	<style>body { color: red; }</style>
</head>
<body>
	<h1>Hello World</h1>
	<p>This is a test paragraph.</p>
</body>
</html>`
		fmt.Fprint(w, html)
	}))
	defer server.Close()

	result, err := WebFetch(server.URL)
	if err != nil {
		t.Errorf("WebFetch() error = %v", err)
		return
	}

	// Should contain text content but not script/style tags
	if !containsString(result, "Hello World") || !containsString(result, "This is a test paragraph") {
		t.Errorf("WebFetch() result should contain page text, got %q", result)
	}

	if containsString(result, "console.log") || containsString(result, "color: red") {
		t.Errorf("WebFetch() result should not contain script or style content, got %q", result)
	}

	// Test case 2: Invalid URL
	_, err = WebFetch("invalid-url")
	if err == nil {
		t.Error("WebFetch() with invalid URL expected error, got nil")
	}

	// Test case 3: Non-200 status code
	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Not Found")
	}))
	defer errorServer.Close()

	_, err = WebFetch(errorServer.URL)
	if err == nil {
		t.Error("WebFetch() with non-200 status expected error, got nil")
	}
}

func TestWebFetchWithDifferentContent(t *testing.T) {
	tests := []struct {
		name             string
		contentType      string
		content          string
		expectedError    bool
		expectedInResult string
	}{
		{
			name:             "Plain text",
			contentType:      "text/plain",
			content:          "This is plain text content",
			expectedError:    false,
			expectedInResult: "This is plain text content",
		},
		{
			name:             "HTML with script and style",
			contentType:      "text/html",
			content:          `<html><head><script>alert('test');</script><style>body{margin:0;}</style></head><body>Content here</body></html>`,
			expectedError:    false,
			expectedInResult: "Content here",
		},
		{
			name:             "Empty body",
			contentType:      "text/html",
			content:          `<html><head></head><body></body></html>`,
			expectedError:    true,
			expectedInResult: "",
		},
		{
			name:             "HTML with special characters",
			contentType:      "text/html",
			content:          `<html><body>Hello &amp; Welcome &copy; 2024</body></html>`,
			expectedError:    false,
			expectedInResult: "Hello & Welcome © 2024",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", tt.contentType)
				fmt.Fprint(w, tt.content)
			}))
			defer server.Close()

			result, err := WebFetch(server.URL)

			if (err != nil) != tt.expectedError {
				t.Errorf("WebFetch() error = %v, wantErr %v", err, tt.expectedError)
				return
			}

			if !tt.expectedError && tt.expectedInResult != "" {
				// For HTML, goquery will decode HTML entities
				// We check if the essential content is there
				if !strings.Contains(result, tt.expectedInResult) {
					// Try checking for unencoded version
					unencoded := strings.ReplaceAll(tt.expectedInResult, "&amp;", "&")
					unencoded = strings.ReplaceAll(unencoded, "&copy;", "©")
					if !strings.Contains(result, unencoded) {
						t.Errorf("WebFetch() result = %q, expected to contain %q", result, tt.expectedInResult)
					}
				}
			}
		})
	}
}

func TestWebFetchWithNetworkErrors(t *testing.T) {
	// Test case 1: Connection timeout
	// Create a server that hangs to test timeout
	hangServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This handler will never respond, causing a timeout
		select {} // Block forever
	}))

	// Close the server immediately to simulate connection error
	hangServer.Close()

	_, err := WebFetch(hangServer.URL)
	if err == nil {
		t.Error("WebFetch() with connection error expected error, got nil")
	}

	// Test case 2: Malformed URL
	malformedURLs := []string{
		"not-a-url",
		"http://",
		"https://[invalid-ipv6",
		"ftp://example.com", // Unsupported protocol
	}

	for _, url := range malformedURLs {
		_, err := WebFetch(url)
		if err == nil {
			t.Errorf("WebFetch() with malformed URL %q expected error, got nil", url)
		}
	}
}

func TestWebFetchUserAgent(t *testing.T) {
	// Test that WebFetch sets a proper User-Agent
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userAgent := r.Header.Get("User-Agent")
		if userAgent == "" {
			t.Error("WebFetch() should set User-Agent header")
		}

		// Check if it contains "Chrome" as specified in the implementation
		if !strings.Contains(userAgent, "Chrome") {
			t.Errorf("WebFetch() User-Agent should contain 'Chrome', got %q", userAgent)
		}

		fmt.Fprint(w, "<html><body>User-Agent test</body></html>")
	}))
	defer server.Close()

	_, err := WebFetch(server.URL)
	if err != nil {
		t.Errorf("WebFetch() error = %v", err)
	}
}

// Test with large HTML content
func TestWebFetchWithLargeContent(t *testing.T) {
	// Generate a large HTML content
	var htmlBuilder strings.Builder
	htmlBuilder.WriteString("<html><head><title>Large Page</title></head><body>")

	for i := 0; i < 1000; i++ {
		htmlBuilder.WriteString(fmt.Sprintf("<p>This is paragraph %d with some content.</p>", i))
	}

	htmlBuilder.WriteString("</body></html>")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, htmlBuilder.String())
	}))
	defer server.Close()

	result, err := WebFetch(server.URL)
	if err != nil {
		t.Errorf("WebFetch() with large content error = %v", err)
		return
	}

	// Should contain some expected content
	if !containsString(result, "This is paragraph 0") || !containsString(result, "This is paragraph 999") {
		t.Error("WebFetch() should handle large content correctly")
	}

	// Result should be reasonably sized (not just empty)
	if len(result) < 1000 {
		t.Errorf("WebFetch() result seems too short: %d characters", len(result))
	}
}

// Example of how to benchmark WebFetch
func BenchmarkWebFetch(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, "<html><body>Benchmark test content</body></html>")
	}))
	defer server.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := WebFetch(server.URL)
		if err != nil {
			b.Fatalf("WebFetch() error = %v", err)
		}
	}
}

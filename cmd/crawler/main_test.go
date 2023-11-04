package main

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/gocolly/colly"
	"github.com/jonesrussell/crawler/internal/crawlResult"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	logs []string
}

func (l *MockLogger) Info(args ...interface{}) {
	l.logs = append(l.logs, fmt.Sprint(args...))
}

func (l *MockLogger) Infof(template string, args ...interface{}) {
	l.logs = append(l.logs, fmt.Sprintf(template, args...))
}

func (l *MockLogger) Errorw(msg string, keysAndValues ...interface{}) {
	l.logs = append(l.logs, fmt.Sprint(keysAndValues...))
}

func (l *MockLogger) Sync() error {
	// Implement the method
	return nil
}

// MockDrug is a mocked object for termmatcher.Related
type MockDrug struct {
	mock.Mock
}

// Related is a mocked implementation for termmatcher.Related
func (m *MockDrug) Related(href string) bool {
	args := m.Called(href)
	return args.Bool(0)
}

// MockRedisWrapper is a mocked object for rediswrapper functions
type MockRedisWrapper struct {
	mock.Mock
}

// SAdd is a mocked implementation for rediswrapper.SAdd
func (m *MockRedisWrapper) SAdd(href string) error {
	args := m.Called(href)
	return args.Error(0)
}

// PublishHref is a mocked implementation for rediswrapper.PublishHref
func (m *MockRedisWrapper) PublishHref(stream, href, group string) error {
	args := m.Called(stream, href, group)
	return args.Error(0)
}

// Del is a mocked implementation for rediswrapper.Del
func (m *MockRedisWrapper) Del() error {
	args := m.Called()
	return args.Error(0)
}

func TestProcessFlags(t *testing.T) {
	// Define a helper function to set os.Args for the test
	setArgs := func(url, searchTerms, crawlSiteID string, maxDepth int, debug bool) {
		os.Args = []string{
			"cmd",
			"-url=" + url,
			"-searchterms=" + searchTerms,
			"-crawlsiteid=" + crawlSiteID,
			"-maxdepth=" + strconv.Itoa(maxDepth),
			"-debug=" + strconv.FormatBool(debug),
		}
	}

	// Test case 1: All required arguments provided
	setArgs("https://www.example.com", "search-term-1,search-term-2", "99", 1, false)
	args := processFlags()
	assert.Equal(t, "https://www.example.com", args.URL, "Expected URL to match")
	assert.Equal(t, "search-term-1,search-term-2", args.SearchTerms, "Expected search terms to match")
	assert.Equal(t, "99", args.CrawlSiteID, "Expected crawlsite id to be 99")
	assert.Equal(t, 1, args.MaxDepth, "Expected max depth to be 1")
	assert.Equal(t, false, args.Debug, "Expected debug to be false")

	// Add more test cases as needed
}

func TestSetupRedis(t *testing.T) {
	// Use the mock logger
	logger = &MockLogger{}

	// Create a Config with your test settings
	config := Config{
		RedisHost: "localhost",
		RedisPort: "6379",
		RedisAuth: "",
	}

	// Call your function with the test config
	setupRedis(config, true)

	// Add assertions here to verify the behavior of your function.
}

func TestConfigureCollector(t *testing.T) {
	domain := "example.com"
	maxDepth := 3

	collector := configureCollector([]string{domain}, maxDepth)

	// Assertions
	assert.NotNil(t, collector, "Expected collector not to be nil")
	assert.True(t, collector.Async, "Expected collector to be asynchronous")
	assert.Equal(t, maxDepth, collector.MaxDepth, "Expected MaxDepth to match the input")
	// You can add more assertions based on your requirements
}

func TestSetupCrawlingLogic(t *testing.T) {
	// Create a new collector
	collector := colly.NewCollector()

	// Define search terms for testing
	searchTerms := []string{
		"DRUG",
		"SMOKE JOINT",
		"GROW OP",
		"CANNABI",
		"IMPAIR",
		"SHOOT",
		"FIREARM",
		"MURDER",
		"COCAIN",
		"POSSESS",
		"BREAK ENTER",
	}

	// Create an empty slice of crawlResult.PageData
	var results []crawlResult.PageData

	// Pass the address of results to setupCrawlingLogic
	setupCrawlingLogic(collector, searchTerms, &results)

	// Your test assertions here
	// ...

	// Assert that the expectations of your mocks are met
	// ...

	// Clean up any resources
	// ...
}

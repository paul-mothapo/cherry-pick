package analyzer

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/cherry-pick/pkg/loadbalancer/core"
)

type URLAnalysisResult struct {
	BaseURL         string           `json:"baseUrl"`
	DiscoveredPages []DiscoveredPage `json:"discoveredPages"`
	AnalysisTime    time.Time        `json:
	
	"analysisTime"`
	TotalPages      int              `json:"totalPages"`
}

type DiscoveredPage struct {
	Path         string `json:"path"`
	Title        string `json:"title"`
	StatusCode   int    `json:"statusCode"`
	ResponseTime int64  `json:"responseTime"`
	IsInternal   bool   `json:"isInternal"`
	Depth        int    `json:"depth"`
}

type URLAnalyzer struct {
	client     *http.Client
	maxDepth   int
	maxPages   int
	visited    map[string]bool
	discovered []DiscoveredPage
}

func NewURLAnalyzer() *URLAnalyzer {
	return &URLAnalyzer{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		maxDepth:   5,
		maxPages:   200,
		visited:    make(map[string]bool),
		discovered: make([]DiscoveredPage, 0),
	}
}

func (ua *URLAnalyzer) AnalyzeURL(baseURL string) (*URLAnalysisResult, error) {
	log.Printf("Starting URL analysis for: %s", baseURL)

	ua.visited = make(map[string]bool)
	ua.discovered = make([]DiscoveredPage, 0)

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		log.Printf("Invalid URL: %v", err)
		return nil, fmt.Errorf("invalid URL: %v", err)
	}

	log.Printf("Parsed URL - Scheme: %s, Host: %s, Path: %s", parsedURL.Scheme, parsedURL.Host, parsedURL.Path)

	log.Printf("Analyzing root page: %s", baseURL)
	ua.analyzePage(baseURL, "", 0)

	log.Printf("Analysis complete. Found %d pages", len(ua.discovered))
	for i, page := range ua.discovered {
		log.Printf("  %d. %s (%s) - Status: %d - Time: %dms",
			i+1, page.Path, page.Title, page.StatusCode, page.ResponseTime)
	}

	return &URLAnalysisResult{
		BaseURL:         baseURL,
		DiscoveredPages: ua.discovered,
		AnalysisTime:    time.Now(),
		TotalPages:      len(ua.discovered),
	}, nil
}

func (ua *URLAnalyzer) analyzePage(pageURL, referrer string, depth int) {
	if depth > ua.maxDepth || len(ua.discovered) >= ua.maxPages {
		log.Printf("Stopping analysis - depth: %d, maxDepth: %d, pages: %d, maxPages: %d",
			depth, ua.maxDepth, len(ua.discovered), ua.maxPages)
		return
	}

	if ua.visited[pageURL] {
		log.Printf("Already visited: %s", pageURL)
		return
	}

	log.Printf("Analyzing page (depth %d): %s", depth, pageURL)
	ua.visited[pageURL] = true

	req, err := http.NewRequest("GET", pageURL, nil)
	if err != nil {
		log.Printf("Failed to create request for %s: %v", pageURL, err)
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	log.Printf("Making request to: %s", pageURL)
	startTime := time.Now()
	resp, err := ua.client.Do(req)
	responseTime := time.Since(startTime).Milliseconds()

	page := DiscoveredPage{
		Path:         ua.getPathFromURL(pageURL),
		Title:        ua.extractTitle(pageURL),
		StatusCode:   http.StatusOK,
		ResponseTime: responseTime,
		IsInternal:   true,
		Depth:        depth,
	}

	if err != nil {
		log.Printf("Request failed for %s: %v", pageURL, err)
		page.StatusCode = http.StatusInternalServerError
		ua.discovered = append(ua.discovered, page)
		return
	}

	page.StatusCode = resp.StatusCode
	log.Printf("Response from %s: Status %d, Time %dms", pageURL, resp.StatusCode, responseTime)
	ua.discovered = append(ua.discovered, page)

	if resp.StatusCode == http.StatusOK {
		var reader io.Reader = resp.Body
		if resp.Header.Get("Content-Encoding") == "gzip" {
			gzReader, err := gzip.NewReader(resp.Body)
			if err != nil {
				log.Printf("Failed to create gzip reader for %s: %v", pageURL, err)
				resp.Body.Close()
				return
			}
			defer gzReader.Close()
			reader = gzReader
		}

		body, err := io.ReadAll(reader)
		resp.Body.Close()

		if err != nil {
			log.Printf("Failed to read response body from %s: %v", pageURL, err)
		} else {
			log.Printf("Reading HTML content from %s (%d bytes)", pageURL, len(body))

			links := ua.extractLinks(string(body), pageURL)
			log.Printf("Found %d links in %s", len(links), pageURL)

			if len(links) == 0 {
				log.Printf("No links found, analyzing content for potential routes...")
				contentLinks := ua.extractRoutesFromContent(string(body), pageURL)
				links = append(links, contentLinks...)
				log.Printf("Found %d additional routes from content analysis", len(contentLinks))
			}

			internalLinks := 0
			for _, link := range links {
				if ua.isInternalLink(link, pageURL) {
					internalLinks++
					if !ua.visited[link] {
						log.Printf("Found internal link: %s", link)
						ua.analyzePage(link, pageURL, depth+1)
					} else {
						log.Printf("Skipping already visited internal link: %s", link)
					}
				} else {
					log.Printf("External link (skipping): %s", link)
				}
			}
			log.Printf("Internal links found: %d/%d", internalLinks, len(links))
		}
	} else {
		log.Printf("Non-200 status code for %s: %d", pageURL, resp.StatusCode)
	}
}

func (ua *URLAnalyzer) extractLinks(html, baseURL string) []string {
	log.Printf("Extracting links from HTML content (%d bytes)", len(html))

	patterns := []*regexp.Regexp{
		regexp.MustCompile(`href=["']([^"']+)["']`),
		regexp.MustCompile(`src=["']([^"']+)["']`),
		regexp.MustCompile(`action=["']([^"']+)["']`),
		regexp.MustCompile(`to=["']([^"']+)["']`),
		regexp.MustCompile(`pathname=["']([^"']+)["']`),
		regexp.MustCompile(`data-href=["']([^"']+)["']`),
		regexp.MustCompile(`data-to=["']([^"']+)["']`),
		regexp.MustCompile(`data-path=["']([^"']+)["']`),
		regexp.MustCompile(`router\.push\(["']([^"']+)["']`),
		regexp.MustCompile(`navigate\(["']([^"']+)["']`),
		regexp.MustCompile(`history\.push\(["']([^"']+)["']`),
		regexp.MustCompile(`window\.location\.href\s*=\s*["']([^"']+)["']`),
		regexp.MustCompile(`url\(["']?([^"')]+)["']?\)`),
		regexp.MustCompile(`background-image:\s*url\(["']?([^"')]+)["']?`),
		regexp.MustCompile(`["']/([a-zA-Z0-9\-_/]+)["']`),
		regexp.MustCompile(`\s/([a-zA-Z0-9\-_/]+)\s`),
		regexp.MustCompile(`>([a-zA-Z0-9\-_/]+)<`),
	}

	links := make([]string, 0)
	base, _ := url.Parse(baseURL)
	seen := make(map[string]bool)

	for i, pattern := range patterns {
		matches := pattern.FindAllStringSubmatch(html, -1)
		log.Printf("Pattern %d found %d matches", i+1, len(matches))
		for _, match := range matches {
			if len(match) > 1 {
				link := strings.TrimSpace(match[1])
				originalLink := link

				if link == "" || strings.HasPrefix(link, "javascript:") ||
					strings.HasPrefix(link, "mailto:") || strings.HasPrefix(link, "tel:") ||
					strings.HasPrefix(link, "#") {
					log.Printf("Skipping link: %s", link)
					continue
				}

				path := strings.TrimPrefix(link, base.Scheme+"://"+base.Host)
				if len(path) <= 1 ||
					path == "/" && len(originalLink) <= 1 ||
					strings.Count(path, "/") > 5 ||
					strings.Contains(path, "..") ||
					strings.Contains(path, " ") ||
					strings.Contains(path, "\\") {
					log.Printf("Skipping low-quality link: %s", link)
					continue
				}

				if strings.HasPrefix(link, "/") {
					link = base.Scheme + "://" + base.Host + link
				} else if strings.HasPrefix(link, "http") {
				} else if !strings.Contains(link, "://") {
					if !strings.HasPrefix(link, "/") {
						link = "/" + link
					}
					link = base.Scheme + "://" + base.Host + link
				}

				log.Printf("Processing link: %s -> %s", originalLink, link)

				if !seen[link] && ua.isInternalLink(link, baseURL) {
					seen[link] = true
					links = append(links, link)
					log.Printf("Added internal link: %s", link)
				} else if seen[link] {
					log.Printf("Already seen link: %s", link)
				} else {
					log.Printf("External link (skipping): %s", link)
				}
			}
		}
	}

	log.Printf("Link extraction complete. Found %d unique links", len(links))
	return links
}

func (ua *URLAnalyzer) extractRoutesFromContent(html, baseURL string) []string {
	log.Printf("Analyzing content for potential routes...")

	routePatterns := []*regexp.Regexp{
		regexp.MustCompile(`<nav[^>]*>.*?</nav>`),
		regexp.MustCompile(`<ul[^>]*class="[^"]*nav[^"]*"[^>]*>.*?</ul>`),
		regexp.MustCompile(`<div[^>]*class="[^"]*menu[^"]*"[^>]*>.*?</div>`),
		regexp.MustCompile(`<button[^>]*>([^<]+)</button>`),
		regexp.MustCompile(`<a[^>]*>([^<]+)</a>`),
		regexp.MustCompile(`["']/([a-zA-Z0-9\-_/]+)["']`),
		regexp.MustCompile(`\s/([a-zA-Z0-9\-_/]+)\s`),
		regexp.MustCompile(`<h[1-6][^>]*>([^<]+)</h[1-6]>`),
		regexp.MustCompile(`<section[^>]*id="([^"]+)"`),
		regexp.MustCompile(`<div[^>]*id="([^"]+)"`),
	}

	potentialRoutes := make([]string, 0)
	seen := make(map[string]bool)
	base, _ := url.Parse(baseURL)

	for i, pattern := range routePatterns {
		matches := pattern.FindAllStringSubmatch(html, -1)
		log.Printf("Content pattern %d found %d matches", i+1, len(matches))

		for _, match := range matches {
			if len(match) > 1 {
				text := strings.TrimSpace(match[1])

				route := ua.textToRoute(text)
				if route != "" && !seen[route] {
					seen[route] = true

					fullURL := base.Scheme + "://" + base.Host + route
					potentialRoutes = append(potentialRoutes, fullURL)
					log.Printf("Potential route found: %s -> %s", text, route)
				}
			}
		}
	}

	log.Printf("Content analysis found %d potential routes", len(potentialRoutes))
	return potentialRoutes
}

func (ua *URLAnalyzer) textToRoute(text string) string {
	text = strings.ToLower(strings.TrimSpace(text))

	if len(text) < 2 || strings.ContainsAny(text, " \t\n\r") {
		return ""
	}

	textToRouteMap := map[string]string{
		"about":            "/about",
		"about me":         "/about",
		"contact":          "/contact",
		"contact us":       "/contact",
		"portfolio":        "/portfolio",
		"projects":         "/projects",
		"work":             "/work",
		"experience":       "/experience",
		"skills":           "/skills",
		"resume":           "/resume",
		"cv":               "/cv",
		"services":         "/services",
		"blog":             "/blog",
		"posts":            "/posts",
		"articles":         "/articles",
		"news":             "/news",
		"docs":             "/docs",
		"documentation":    "/docs",
		"help":             "/help",
		"support":          "/support",
		"privacy":          "/privacy",
		"privacy policy":   "/privacy",
		"terms":            "/terms",
		"terms of service": "/terms",
		"home":             "/",
		"main":             "/",
		"dashboard":        "/dashboard",
		"profile":          "/profile",
		"settings":         "/settings",
		"admin":            "/admin",
		"login":            "/login",
		"signup":           "/signup",
		"register":         "/signup",
	}

	if route, exists := textToRouteMap[text]; exists {
		return route
	}

	if strings.HasPrefix(text, "/") {
		return text
	}

	route := "/" + strings.ReplaceAll(text, " ", "-")
	return route
}

func (ua *URLAnalyzer) isInternalLink(link, baseURL string) bool {
	linkURL, err := url.Parse(link)
	if err != nil {
		return false
	}

	baseURLParsed, err := url.Parse(baseURL)
	if err != nil {
		return false
	}

	return linkURL.Host == baseURLParsed.Host
}

func (ua *URLAnalyzer) getPathFromURL(pageURL string) string {
	parsed, err := url.Parse(pageURL)
	if err != nil {
		return pageURL
	}
	return parsed.Path
}

func (ua *URLAnalyzer) extractTitle(pageURL string) string {
	path := ua.getPathFromURL(pageURL)

	if path == "/" {
		return "Home Page"
	}

	title := strings.Trim(path, "/")
	title = strings.ReplaceAll(title, "-", " ")
	title = strings.ReplaceAll(title, "_", " ")
	words := strings.Fields(title)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}
	title = strings.Join(words, " ")

	if title == "" {
		title = "Page"
	}

	return title
}

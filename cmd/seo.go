package main

import (
	"fmt"
	"time"
)

// SEOData contains metadata for search engine optimization
type SEOData struct {
	Title         string
	Description   string
	Keywords      []string
	CanonicalURL  string
	OGType        string
	OGImage       string
	OGImageAlt    string
	TwitterCard   string
	Author        string
	PublishedTime string
	ModifiedTime  string
	SchemaJSON    string
}

// DefaultSEO returns SEO data for the homepage
func DefaultSEO() *SEOData {
	return &SEOData{
		Title:        "Corbin Lazarone - Software Engineer Portfolio",
		Description:  "Software engineer specializing in backend development, distributed systems, and cloud infrastructure. Explore my projects, experience, and technical blog.",
		Keywords:     []string{"software engineer", "backend developer", "golang", "distributed systems", "portfolio", "Corbin Lazarone"},
		CanonicalURL: "https://corbinlazarone.com/",
		OGType:       "website",
		OGImage:      "/static/images/og-image.png",
		OGImageAlt:   "Corbin Lazarone - Software Engineer",
		TwitterCard:  "summary_large_image",
		Author:       "Corbin Lazarone",
		SchemaJSON:   generatePersonSchema(),
	}
}

// BlogListingSEO returns SEO data for the blog listing page
func BlogListingSEO() *SEOData {
	return &SEOData{
		Title:        "Blog - Corbin Lazarone",
		Description:  "Essays and thoughts on software engineering, technology, personal growth, and life experiences. From technical deep-dives to reflections on learning and building.",
		Keywords:     []string{"blog", "software engineering", "technology", "personal development", "essays", "technical writing", "life experiences"},
		CanonicalURL: "https://corbinlazarone.com/blog",
		OGType:       "website",
		OGImage:      "/static/images/og-blog.png",
		OGImageAlt:   "Corbin Lazarone - Blog",
		TwitterCard:  "summary_large_image",
		Author:       "Corbin Lazarone",
		SchemaJSON:   generateBlogSchema(),
	}
}

// BlogPostSEO returns SEO data for a specific blog post
func BlogPostSEO(post *BlogPost, fullURL string) *SEOData {
	publishedTime := ""
	if !post.Date.IsZero() {
		publishedTime = post.Date.Format(time.RFC3339)
	}

	return &SEOData{
		Title:         fmt.Sprintf("%s - Corbin Lazarone", post.Title),
		Description:   post.Excerpt,
		Keywords:      post.Tags,
		CanonicalURL:  fullURL,
		OGType:        "article",
		OGImage:       "/static/images/og-blog.png",
		OGImageAlt:    fmt.Sprintf("%s by Corbin Lazarone", post.Title),
		TwitterCard:   "summary_large_image",
		Author:        "Corbin Lazarone",
		PublishedTime: publishedTime,
		ModifiedTime:  publishedTime,
		SchemaJSON:    generateBlogPostSchema(post, fullURL),
	}
}

// generatePersonSchema creates JSON-LD structured data for the Person schema
func generatePersonSchema() string {
	return `{
  "@context": "https://schema.org",
  "@type": "Person",
  "name": "Corbin Lazarone",
  "url": "https://corbinlazarone.com",
  "jobTitle": "Software Engineer",
  "sameAs": [
    "https://github.com/corbinlazarone",
    "https://linkedin.com/in/corbinlazarone",
    "https://x.com/ccorrzzyy"
  ],
  "knowsAbout": ["Software Engineering", "Backend Development", "Distributed Systems", "Go Programming"]
}`
}

// generateBlogSchema creates JSON-LD structured data for the Blog schema
func generateBlogSchema() string {
	return `{
  "@context": "https://schema.org",
  "@type": "Blog",
  "name": "Corbin Lazarone's Blog",
  "description": "Essays on software engineering, technology, personal growth, and life experiences",
  "url": "https://corbinlazarone.com/blog",
  "author": {
    "@type": "Person",
    "name": "Corbin Lazarone"
  }
}`
}

// generateBlogPostSchema creates JSON-LD structured data for a BlogPosting
func generateBlogPostSchema(post *BlogPost, fullURL string) string {
	datePublished := ""
	if !post.Date.IsZero() {
		datePublished = post.Date.Format(time.RFC3339)
	}

	return fmt.Sprintf(`{
  "@context": "https://schema.org",
  "@type": "BlogPosting",
  "headline": "%s",
  "description": "%s",
  "url": "%s",
  "datePublished": "%s",
  "author": {
    "@type": "Person",
    "name": "Corbin Lazarone"
  },
  "publisher": {
    "@type": "Person",
    "name": "Corbin Lazarone"
  }
}`, post.Title, post.Excerpt, fullURL, datePublished)
}

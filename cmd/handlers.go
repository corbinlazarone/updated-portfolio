package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

// TemplateData holds data passed to templates
type TemplateData struct {
	SEO   *SEOData
	Posts []*BlogPost
	Post  *BlogPost
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./ui/html/home.tmpl.html",
		"./ui/html/ascii.tmpl.html",
	}

	// read each template file and parse it
	tempSet, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Prepare template data with SEO
	data := TemplateData{
		SEO: DefaultSEO(),
	}

	err = tempSet.ExecuteTemplate(w, "home", data)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func blog(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/blog" {
		http.NotFound(w, r)
		return
	}

	// read each template file and parse it
	tempSet, err := template.ParseFiles("./ui/html/blog.tmpl.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Prepare template data with SEO and blog posts
	data := TemplateData{
		SEO:   BlogListingSEO(),
		Posts: GetAllBlogPosts(),
	}

	err = tempSet.ExecuteTemplate(w, "blog", data)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func blogPost(w http.ResponseWriter, r *http.Request) {
	// Extract slug from URL path (/blog/slug)
	slug := strings.TrimPrefix(r.URL.Path, "/blog/")
	if slug == "" || slug == "blog" {
		http.NotFound(w, r)
		return
	}

	// Get blog post by slug
	post, exists := GetBlogPostBySlug(slug)
	if !exists {
		http.NotFound(w, r)
		return
	}

	// Parse template
	tempSet, err := template.ParseFiles("./ui/html/blog-post.tmpl.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Prepare template data with SEO and blog post
	fullURL := "https://corbinlazarone.com/blog/" + slug
	data := TemplateData{
		SEO:  BlogPostSEO(post, fullURL),
		Post: post,
	}

	err = tempSet.ExecuteTemplate(w, "blog-post", data)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func robotsTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	robots := `User-agent: *
Allow: /

Sitemap: https://corbinlazarone.com/sitemap.xml`
	w.Write([]byte(robots))
}

func sitemapXML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")

	sitemap := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>https://corbinlazarone.com/</loc>
    <changefreq>weekly</changefreq>
    <priority>1.0</priority>
  </url>
  <url>
    <loc>https://corbinlazarone.com/blog</loc>
    <changefreq>weekly</changefreq>
    <priority>0.8</priority>
  </url>`

	// Add all blog posts to sitemap
	posts := GetAllBlogPosts()
	for _, post := range posts {
		sitemap += "\n  <url>\n"
		sitemap += "    <loc>https://corbinlazarone.com/blog/" + post.Slug + "</loc>\n"
		if !post.Date.IsZero() {
			sitemap += "    <lastmod>" + post.Date.Format("2006-01-02") + "</lastmod>\n"
		}
		sitemap += "    <changefreq>monthly</changefreq>\n"
		sitemap += "    <priority>0.6</priority>\n"
		sitemap += "  </url>"
	}

	sitemap += "\n</urlset>"
	w.Write([]byte(sitemap))
}

package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Load blog posts on startup
	if err := LoadBlogPosts(); err != nil {
		log.Printf("Warning: failed to load blog posts: %v", err)
	} else {
		log.Printf("Loaded %d blog posts", len(GetAllBlogPosts()))
	}

	// Get port from environment variable or flag
	portEnv := os.Getenv("PORT")
	defaultPort := ":4000"
	if portEnv != "" {
		defaultPort = ":" + portEnv
	}

	port := flag.String("port", defaultPort, "port the server runs on")
	flag.Parse()

	mux := http.NewServeMux()

	// Static file server
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// SEO routes
	mux.HandleFunc("/robots.txt", robotsTxt)
	mux.HandleFunc("/sitemap.xml", sitemapXML)

	// Page routes
	mux.HandleFunc("/", home)
	mux.HandleFunc("/blog", blog)

	// Blog post routes - handle /blog/* paths
	mux.HandleFunc("/blog/", func(w http.ResponseWriter, r *http.Request) {
		// If path is exactly "/blog", redirect to blog listing
		if r.URL.Path == "/blog" || r.URL.Path == "/blog/" {
			// Ensure we don't have a trailing slash for blog listing
			if strings.HasSuffix(r.URL.Path, "/") {
				http.Redirect(w, r, "/blog", http.StatusMovedPermanently)
				return
			}
			blog(w, r)
			return
		}
		// Otherwise, it's a blog post
		blogPost(w, r)
	})

	srv := &http.Server{
		Addr:    *port,
		Handler: mux,
	}

	log.Printf("Listening on port %s", *port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

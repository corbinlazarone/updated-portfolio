package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/russross/blackfriday/v2"
	"gopkg.in/yaml.v3"
)

// BlogPost represents a parsed blog post
type BlogPost struct {
	Title       string    `yaml:"title"`
	Date        time.Time `yaml:"date"`
	Tags        []string  `yaml:"tags"`
	Excerpt     string    `yaml:"excerpt"`
	Slug        string    `yaml:"slug"`
	Content     string    // Markdown content
	HTMLContent template.HTML // Rendered HTML
	ReadingTime int       // Minutes
}

// BlogPostMetadata is used for parsing just the frontmatter
type BlogPostMetadata struct {
	Title   string    `yaml:"title"`
	Date    time.Time `yaml:"date"`
	Tags    []string  `yaml:"tags"`
	Excerpt string    `yaml:"excerpt"`
	Slug    string    `yaml:"slug"`
}

var (
	blogPosts     []*BlogPost
	blogPostsMap  = make(map[string]*BlogPost)
)

// LoadBlogPosts loads all blog posts from the content/blog directory
func LoadBlogPosts() error {
	files, err := ioutil.ReadDir("./content/blog")
	if err != nil {
		// If directory doesn't exist yet, just return empty
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read blog directory: %w", err)
	}

	posts := []*BlogPost{}
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".md") {
			continue
		}

		filePath := filepath.Join("./content/blog", file.Name())
		post, err := ParseBlogPost(filePath)
		if err != nil {
			fmt.Printf("Warning: failed to parse %s: %v\n", file.Name(), err)
			continue
		}

		posts = append(posts, post)
		blogPostsMap[post.Slug] = post
	}

	// Sort posts by date (newest first)
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})

	blogPosts = posts
	return nil
}

// ParseBlogPost parses a markdown file with YAML frontmatter
func ParseBlogPost(filePath string) (*BlogPost, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Split frontmatter and content
	parts := bytes.SplitN(content, []byte("---"), 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid frontmatter format")
	}

	// Parse frontmatter
	var metadata BlogPostMetadata
	if err := yaml.Unmarshal(parts[1], &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	// Get markdown content
	markdownContent := string(parts[2])

	// Create blog post
	post := &BlogPost{
		Title:       metadata.Title,
		Date:        metadata.Date,
		Tags:        metadata.Tags,
		Excerpt:     metadata.Excerpt,
		Slug:        metadata.Slug,
		Content:     markdownContent,
		ReadingTime: calculateReadingTime(markdownContent),
	}

	// Render markdown to HTML with syntax highlighting
	htmlContent, err := renderMarkdown(markdownContent)
	if err != nil {
		return nil, fmt.Errorf("failed to render markdown: %w", err)
	}
	post.HTMLContent = template.HTML(htmlContent)

	return post, nil
}

// renderMarkdown converts markdown to HTML with syntax highlighting
func renderMarkdown(markdown string) (string, error) {
	// Create a custom renderer for syntax highlighting
	renderer := &chromaRenderer{
		html: blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
			Flags: blackfriday.CommonHTMLFlags,
		}),
	}

	// Parse and render markdown
	extensions := blackfriday.CommonExtensions | blackfriday.AutoHeadingIDs | blackfriday.Footnotes
	output := blackfriday.Run([]byte(markdown), blackfriday.WithRenderer(renderer), blackfriday.WithExtensions(extensions))

	return string(output), nil
}

// chromaRenderer is a custom renderer that adds syntax highlighting to code blocks
type chromaRenderer struct {
	html *blackfriday.HTMLRenderer
}

func (r *chromaRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	// Handle code blocks with syntax highlighting
	if node.Type == blackfriday.CodeBlock {
		lang := string(node.CodeBlockData.Info)
		code := string(node.Literal)

		highlighted, err := highlightCode(code, lang)
		if err != nil {
			// Fall back to default rendering if highlighting fails
			return r.html.RenderNode(w, node, entering)
		}

		w.Write([]byte(highlighted))
		return blackfriday.GoToNext
	}

	// Delegate to default renderer for other nodes
	return r.html.RenderNode(w, node, entering)
}

func (r *chromaRenderer) RenderHeader(w io.Writer, ast *blackfriday.Node) {
	r.html.RenderHeader(w, ast)
}

func (r *chromaRenderer) RenderFooter(w io.Writer, ast *blackfriday.Node) {
	r.html.RenderFooter(w, ast)
}

// highlightCode applies syntax highlighting to code
func highlightCode(code, lang string) (string, error) {
	// Get lexer for the language
	lexer := lexers.Get(lang)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	// Use a terminal-friendly style (monokai is dark)
	style := styles.Get("monokai")
	if style == nil {
		style = styles.Fallback
	}

	// Create formatter
	formatter := html.New(
		html.WithClasses(true),
		html.PreventSurroundingPre(false),
	)

	// Tokenize
	tokens, err := lexer.Tokenise(nil, code)
	if err != nil {
		return "", err
	}

	// Format to HTML
	var buf bytes.Buffer
	err = formatter.Format(&buf, style, tokens)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// calculateReadingTime estimates reading time in minutes (assuming 200 words per minute)
func calculateReadingTime(content string) int {
	words := strings.Fields(content)
	wordCount := len(words)
	minutes := wordCount / 200
	if minutes == 0 {
		minutes = 1
	}
	return minutes
}

// GetAllBlogPosts returns all blog posts sorted by date
func GetAllBlogPosts() []*BlogPost {
	return blogPosts
}

// GetBlogPostBySlug returns a blog post by its slug
func GetBlogPostBySlug(slug string) (*BlogPost, bool) {
	post, exists := blogPostsMap[slug]
	return post, exists
}

package handlers

import (
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/i18n"
	"github.com/getfider/fider/app/pkg/markdown"

	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/web"
)

type AtomFeed struct {
	XMLName  xml.Name `xml:"http://www.w3.org/2005/Atom feed"`
	Title    string   `xml:"title"`
	Subtitle Content  `xml:"subtitle"`
	Id       string   `xml:"id"`
	Updated  string   `xml:"updated"`
	Link     []Link   `xml:"link"`
	Author   *Author  `xml:"author"`
	Entries  []*Entry `xml:"entry"`
}

type Entry struct {
	Title      string      `xml:"title,omitempty"`
	Id         string      `xml:"id"`
	Published  string      `xml:"published"`
	Updated    string      `xml:"updated,omitempty"`
	Link       []Link      `xml:"link"`
	Author     *Author     `xml:"author"`
	Summary    *Content    `xml:"summary"`
	Content    *Content    `xml:"content"`
	Categories []*Category `xml:"category,omitempty"`
}

type Category struct {
	Term string `xml:"term,attr"`
}

type Link struct {
	Rel      string `xml:"rel,attr,omitempty"`
	Href     string `xml:"href,attr"`
	Type     string `xml:"type,attr,omitempty"`
	HrefLang string `xml:"hreflang,attr,omitempty"`
	Title    string `xml:"title,attr,omitempty"`
	Length   uint   `xml:"length,attr,omitempty"`
}

type Author struct {
	Name     string `xml:"name"`
	Uri      string `xml:"uri,omitempty"`
	Email    string `xml:"email,omitempty"`
	InnerXML string `xml:",innerxml"`
}

type Content struct {
	Type string `xml:"type,attr"`
	Body string `xml:",chardata"`
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02T15:04:05-07:00")
}

func generateXML(feed *AtomFeed) (string, error) {
	indent := ""
	if env.IsDevelopment() {
		indent = "	"
	}

	feedXML, err := xml.MarshalIndent(feed, "", indent)
	if err != nil {
		return "", err
	}

	return "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n" + string(feedXML), nil
}

type generatorOptions struct {
	generateTitle  bool
	generateFooter bool
}

func generatePostContent(c *web.Context, post *entity.Post, options *generatorOptions) string {
	title := ""
	if options.generateTitle {
		title = i18n.T(c, "feed.post.title", i18n.Params{
			"votes":    post.VotesCount,
			"comments": post.CommentsCount,
			"title":    post.Title,
		})
	}

	footer := ""
	if options.generateFooter {
		responseFooter := ""
		if (post.Response != nil) && (post.Response.Text != "") {
			responseFooter = i18n.T(c, "feed.post.footer.response", i18n.Params{
				"responder": post.Response.User.Name,
				"date":      post.Response.RespondedAt.Format("Jan 2, 2006"),
				"response":  strings.ReplaceAll(post.Response.Text, "\n", "\n>"),
			})
		}

		footer = i18n.T(c, "feed.post.footer", i18n.Params{
			"response_footer": responseFooter,
			"votes":           post.VotesCount,
			"comments":        post.CommentsCount,
			"web_link":        fmt.Sprintf("%s/posts/%d", web.BaseURL(c), post.Number),
			"feed_link":       fmt.Sprintf("%s/feed/posts/%d.atom", web.BaseURL(c), post.Number),
		})
	}

	return string(markdown.Full(title+post.Description+footer, true))
}

func appendTags(c *web.Context, categories []*Category, post *entity.Post) ([]*Category, error) {
	getAssignedTags := &query.GetAssignedTags{Post: post}
	if err := bus.Dispatch(c, getAssignedTags); err != nil {
		return nil, err
	}
	tags := getAssignedTags.Result

	for _, tag := range tags {
		if tag.IsPublic {
			categories = append(categories, &Category{Term: tag.Name})
		}
	}

	return categories, nil
}

// GlobalFeed Returns the global ATOM feed with the 30 most recent posts as entries
func GlobalFeed() web.HandlerFunc {
	return func(c *web.Context) error {
		if c.Tenant().IsPrivate || !c.Tenant().IsFeedEnabled {
			return c.NotFound()
		}

		searchPosts := &query.SearchPosts{
			Query: c.QueryParam("query"),
			View:  "all",
			Limit: "30",
			Tags:  c.QueryParamAsArray("tags"),
		}
		if err := bus.Dispatch(c, searchPosts); err != nil {
			return c.Failure(err)
		}
		posts := searchPosts.Result

		feed := &AtomFeed{
			Title:    c.Tenant().Name,
			Subtitle: Content{Body: string(markdown.Full(c.Tenant().WelcomeMessage, true)), Type: "html"},
			Id:       web.BaseURL(c),
			Link: []Link{
				{Href: fmt.Sprintf("%s/feed/global.atom", web.BaseURL(c)), Type: "application/atom+xml", Rel: "self"},
				{Href: web.BaseURL(c), Type: "text/html", Rel: "alternate"},
			},
			Entries: []*Entry{},
		}

		lastUpdate := time.UnixMilli(0)
		for _, post := range posts {
			if post.CreatedAt.After(lastUpdate) {
				lastUpdate = post.CreatedAt
			}
			if (post.Response != nil) && post.Response.RespondedAt.After(lastUpdate) {
				lastUpdate = post.Response.RespondedAt
			}

			categories := []*Category{{Term: i18n.T(c, "enum.poststatus."+post.Status.Name())}}
			categories, err := appendTags(c, categories, post)
			if err != nil {
				return c.Failure(err)
			}

			feed.Entries = append(feed.Entries, &Entry{
				Title: i18n.T(c, "feed.global.title", i18n.Params{
					"count": post.VotesCount,
					"title": post.Title,
				}),
				Author:    &Author{Name: post.User.Name},
				Published: formatTime(post.CreatedAt),
				Updated: func() string {
					if post.Response == nil {
						return formatTime(post.CreatedAt)
					}
					return formatTime(post.Response.RespondedAt)
				}(),
				Content: &Content{Type: "html", Body: generatePostContent(c, post, &generatorOptions{generateTitle: false, generateFooter: true})},
				Id:      fmt.Sprintf("%s/posts/%d", web.BaseURL(c), post.Number),
				Link: []Link{
					{Href: fmt.Sprintf("%s/feed/posts/%d.atom", web.BaseURL(c), post.Number), Type: "application/atom+xml", Rel: "self"},
					{Href: fmt.Sprintf("%s/posts/%d", web.BaseURL(c), post.Number), Type: "text/html", Rel: "alternate"},
				},
				Categories: categories,
			})
		}
		feed.Updated = formatTime(lastUpdate)

		feedStr, err := generateXML(feed)
		if err != nil {
			return c.Failure(err)
		}

		return c.Blob(http.StatusOK, "application/atom+xml", []byte(feedStr))
	}
}

// CommentFeed Returns the ATOM feed for a single post, with its comments as entries
func CommentFeed() web.HandlerFunc {
	return func(c *web.Context) error {
		if c.Tenant().IsPrivate || !c.Tenant().IsFeedEnabled {
			return c.NotFound()
		}

		path := c.Param("path")
		path, found := strings.CutSuffix(path, ".atom")
		if !found {
			return c.NotFound()
		}
		number, err := strconv.Atoi(path)
		if err != nil {
			return c.NotFound()
		}

		getPost := &query.GetPostByNumber{Number: number}
		if err := bus.Dispatch(c, getPost); err != nil {
			return c.Failure(err)
		}
		getComments := &query.GetCommentsByPost{Post: getPost.Result}
		if err := bus.Dispatch(c, getComments); err != nil {
			return c.Failure(err)
		}
		post := getPost.Result
		comments := getComments.Result
		comments = comments[max(0, len(comments)-30):] // get the last 30 comments

		feed := &AtomFeed{
			Title:    post.Title,
			Subtitle: Content{Body: string(markdown.Full(post.Description, true)), Type: "html"},
			Author:   &Author{Name: post.User.Name},
			Id:       fmt.Sprintf("%s/posts/%d/#comments", web.BaseURL(c), post.Number),
			Link: []Link{
				{Href: fmt.Sprintf("%s/feed/posts/%d.atom", web.BaseURL(c), post.Number), Type: "application/atom+xml", Rel: "self"},
				{Href: fmt.Sprintf("%s/posts/%d", web.BaseURL(c), post.Number), Type: "text/html", Rel: "alternate"},
			},
			Entries: []*Entry{},
		}

		categories := []*Category{{Term: i18n.T(c, "enum.poststatus."+post.Status.Name())}}
		categories, err = appendTags(c, categories, post)
		if err != nil {
			return c.Failure(err)
		}

		feed.Entries = append(feed.Entries, &Entry{
			Title: i18n.T(c, "feed.comment.op", i18n.Params{
				"author": post.User.Name,
			}),
			Author:    &Author{Name: post.User.Name},
			Published: formatTime(post.CreatedAt),
			Updated:   formatTime(post.CreatedAt),
			Id:        fmt.Sprintf("%s/posts/%d", web.BaseURL(c), post.Number),
			Link: []Link{
				{Href: fmt.Sprintf("%s/posts/%d", web.BaseURL(c), post.Number), Type: "text/html", Rel: "alternate"},
			},
			Content:    &Content{Type: "html", Body: generatePostContent(c, post, &generatorOptions{generateTitle: true, generateFooter: false})},
			Categories: categories,
		})
		if (post.Response != nil) && (post.Response.Text != "") {
			feed.Entries = append(feed.Entries, &Entry{
				Title: i18n.T(c, "feed.comment.response", i18n.Params{
					"author": post.Response.User.Name,
				}),
				Author:    &Author{Name: post.Response.User.Name},
				Published: formatTime(post.Response.RespondedAt),
				Updated:   formatTime(post.Response.RespondedAt), // so that it shows as "updated" on edit / new response
				Id:        fmt.Sprintf("%s/posts/%d/#response", web.BaseURL(c), post.Number),
				Link: []Link{
					{Href: fmt.Sprintf("%s/posts/%d", web.BaseURL(c), post.Number), Type: "text/html", Rel: "alternate"},
				},
				Content:    &Content{Type: "html", Body: string(markdown.Full(post.Response.Text, true))},
				Categories: []*Category{{Term: i18n.T(c, "enum.poststatus."+post.Status.Name())}},
			})
		}

		lastUpdate := time.UnixMilli(0)
		for _, comment := range comments {
			if comment.CreatedAt.After(lastUpdate) {
				lastUpdate = comment.CreatedAt
			}
			if (comment.EditedAt != nil) && comment.EditedAt.After(lastUpdate) {
				lastUpdate = *comment.EditedAt
			}

			feed.Entries = append(feed.Entries, &Entry{
				Title: i18n.T(c, "feed.comment.title", i18n.Params{
					"author": comment.User.Name,
				}),
				Author:    &Author{Name: comment.User.Name},
				Published: formatTime(comment.CreatedAt),
				Updated: func() string {
					if comment.EditedAt == nil {
						return formatTime(comment.CreatedAt)
					}
					return formatTime(*comment.EditedAt)
				}(),
				Content: &Content{Type: "html", Body: string(markdown.Full(html.UnescapeString(comment.Content), true))},
				Id:      fmt.Sprintf("%s/posts/%d/#comment-%d", web.BaseURL(c), post.Number, comment.ID),
				Link:    []Link{{Href: fmt.Sprintf("%s/posts/%d/#comment-%d", web.BaseURL(c), post.Number, comment.ID), Type: "text/html", Rel: "alternate"}},
			})
		}
		feed.Updated = formatTime(lastUpdate)

		feedStr, err := generateXML(feed)
		if err != nil {
			return c.Failure(err)
		}

		return c.Blob(http.StatusOK, "application/atom+xml", []byte(feedStr))
	}
}

package connector

import (
	"cognix.ch/api/v2/core/model"
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/gocolly/colly/v2"
	"go.uber.org/zap"
	"jaytaylor.com/html2text"
	"net/url"
	"strings"
	"time"
)

type (
	Web struct {
		Base
		param   *WebParameters
		scraper *colly.Collector
		history map[string]struct{}
	}
	WebParameters struct {
		URL string `url:"url"`
	}
)

func (c *Web) Config(connector *model.Connector) (Connector, error) {
	c.Base.Config(connector)
	c.param = &WebParameters{}
	c.history = make(map[string]struct{})
	if err := connector.ConnectorSpecificConfig.ToStruct(c.param); err != nil {
		return nil, err
	}
	c.scraper = colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)
	return c, nil
}

func (c *Web) Execute(ctx context.Context, param model.JSONMap) (*model.Connector, error) {

	c.scraper.OnHTML("body", c.onNextURL)
	c.history[c.param.URL] = struct{}{}
	err := c.scraper.Visit(c.param.URL)
	if err != nil {
		zap.L().Error("Failed to scrape URL", zap.String("url", c.param.URL), zap.Error(err))
	}
	return c.model, err
}

func NewWeb(connector *model.Connector) (Connector, error) {
	web := Web{}
	return web.Config(connector)
}

func (c *Web) onNextURL(e *colly.HTMLElement) {
	child := e.ChildAttrs("a", "href")
	c.processChildLinks(e.Request.URL, child)

	text, _ := html2text.FromString(e.ChildText("main"), html2text.Options{
		PrettyTables: true,
		PrettyTablesOptions: &html2text.PrettyTablesOptions{
			AutoFormatHeader: true,
			AutoWrapText:     true,
		},
		OmitLinks: true,
	})

	signature := fmt.Sprintf("%x", sha256.Sum256([]byte(text)))
	docID := e.Request.URL.String()
	doc, ok := c.model.DocsMap[docID]
	if !ok {
		doc = &model.Document{
			DocumentID:  docID,
			ConnectorID: c.model.ID,
			Link:        docID,
			CreatedDate: time.Now().UTC(),
			IsExists:    true,
			IsUpdated:   true,
		}
		c.model.DocsMap[docID] = doc
		c.model.Docs = append(c.model.Docs, doc)
	}
	doc.IsExists = true
	if doc.Signature == signature {
		return
	}
	doc.Signature = signature
	if doc.ID != 0 {
		doc.IsUpdated = true
		doc.UpdatedDate = pg.NullTime{time.Now().UTC()}
	}

	// todo send text for indexing
	fmt.Println(text)

}

func (c *Web) processChildLinks(baseURL *url.URL, urls []string) {
	for _, u := range urls {
		if u[0] == '#' || !strings.Contains(u, baseURL.Path) {
			continue
		}
		if strings.HasPrefix(u, baseURL.Path) {
			u = fmt.Sprintf("%s://%s%s", baseURL.Scheme, baseURL.Host, u)
		}
		if _, ok := c.history[u]; ok {
			continue
		}
		c.history[u] = struct{}{}

		if err := c.scraper.Visit(u); err != nil {
			zap.S().Errorf("Failed to scrape URL: %s", u)
		}
	}
	return
}

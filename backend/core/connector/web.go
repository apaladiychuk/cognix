package connector

import (
	"cognix.ch/api/v2/core/model"
	"cognix.ch/api/v2/core/proto"
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/gocolly/colly/v2"
	"go.uber.org/zap"
	"jaytaylor.com/html2text"
	"net/url"
	"strings"
)

type (
	Web struct {
		Base
		param   *WebParameters
		scraper *colly.Collector
		history map[string]string
		ctx     context.Context
	}
	WebParameters struct {
		URL string `url:"url"`
	}
)

var excludeTag = map[string]bool{
	"script": true,
	"button": true,
	"a":      true,
	"header": true,
	"nav":    true,
}

func withContext(ctx context.Context, fn func(context.Context, *colly.HTMLElement)) colly.HTMLCallback {
	return func(element *colly.HTMLElement) {
		fn(ctx, element)
	}
}

func (c *Web) Execute(ctx context.Context, param map[string]string) chan *proto.TriggerResponse {
	zap.S().Debugf("Run web connector with param %s ...", c.param.URL)
	c.ctx = ctx
	go func() {
		c.scraper.OnHTML("body", withContext(ctx, c.onBody))
		err := c.scraper.Visit(c.param.URL)
		if err != nil {
			zap.L().Error("Failed to scrape URL", zap.String("url", c.param.URL), zap.Error(err))
		}
		zap.S().Debugf("Complete web connector with param %s", c.param.URL)
		close(c.resultCh)
	}()
	return c.resultCh
}

func NewWeb(connector *model.Connector) (Connector, error) {
	web := Web{}
	web.Base.Config(connector)
	web.param = &WebParameters{}
	web.history = make(map[string]string)
	if err := connector.ConnectorSpecificConfig.ToStruct(web.param); err != nil {
		return nil, err
	}
	web.scraper = colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)
	return &web, nil
}

func (c *Web) onBody(ctx context.Context, e *colly.HTMLElement) {
	child := e.ChildAttrs("a", "href")
	text, _ := html2text.FromString(e.ChildText("main"), html2text.Options{
		PrettyTables: true,
		PrettyTablesOptions: &html2text.PrettyTablesOptions{
			AutoFormatHeader: true,
			AutoWrapText:     true,
		},
		OmitLinks: true,
	})
	c.history[e.Request.URL.String()] = text
	c.processChildLinks(e.Request.URL, child)
	signature := fmt.Sprintf("%x", sha256.Sum256([]byte(text)))
	docID := e.Request.URL.String()

	c.resultCh <- &proto.TriggerResponse{
		DocumentId: docID,
		Url:        docID,
		Content:    text,
		Signature:  signature,
		Status:     proto.Status_SUCCESS,
	}

}

func (c *Web) onBody2(ctx context.Context, e *colly.HTMLElement) {
	//child := e.ChildAttrs("a", "href")
	//
	//var rows []string
	//e.ForEach("*", func(i int, element *colly.HTMLElement) {
	//
	//	if _, ok := excludeTag[element.Name]; ok {
	//		return
	//	}
	//	fmt.Println(fmt.Sprintf("%d >> %s ", i, element.Name))
	//	rows = append(rows, element.Text)
	//})
	//
	//text, _ := html2text.FromString(strings.Join(rows, "\n"), html2text.Options{
	//	PrettyTables: true,
	//	PrettyTablesOptions: &html2text.PrettyTablesOptions{
	//		AutoFormatHeader: true,
	//		AutoWrapText:     true,
	//	},
	//	OmitLinks: true,
	//})
	//c.history[e.Request.URL.String()] = text
	//c.processChildLinks(e.Request.URL, child)
	//signature := fmt.Sprintf("%x", sha256.Sum256([]byte(text)))
	//docID := e.Request.URL.String()
	//doc, ok := c.model.DocsMap[docID]
	//if !ok {
	//	doc = &model.Document{
	//		DocumentID:  docID,
	//		ConnectorID: c.model.ID,
	//		Link:        docID,
	//		CreatedDate: time.Now().UTC(),
	//		IsExists:    true,
	//		IsUpdated:   true,
	//	}
	//	c.model.DocsMap[docID] = doc
	//	c.model.Docs = append(c.model.Docs, doc)
	//}
	//doc.IsExists = true
	//if doc.Signature == signature {
	//	return
	//}
	//doc.Signature = signature
	//if doc.ID.IntPart() != 0 {
	//	doc.IsUpdated = true
	//	doc.UpdatedDate = pg.NullTime{time.Now().UTC()}
	//}
	//c.resultCh <- &proto.TriggerResponse{
	//	DocumentId: doc.ID.IntPart(),
	//	Url:        doc.DocumentID,
	//	Content:    text,
	//	Signature:  signature,
	//}

}

func (c *Web) processChildLinks(baseURL *url.URL, urls []string) {
	for _, u := range urls {
		if len(u) == 0 || u[0] == '#' || !strings.Contains(u, baseURL.Path) ||
			(strings.HasPrefix(u, "http") && !strings.Contains(u, baseURL.Host)) {
			continue
		}
		if strings.HasPrefix(u, baseURL.Path) {
			u = fmt.Sprintf("%s://%s%s", baseURL.Scheme, baseURL.Host, u)
		}
		if _, ok := c.history[u]; ok {
			continue
		}
		if err := c.scraper.Visit(u); err != nil {
			zap.S().Errorf("Failed to scrape URL: %s", u)
		}
	}
	return
}

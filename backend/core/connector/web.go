package connector

import (
	"cognix.ch/api/v2/core/model"
	"context"
	"github.com/gocolly/colly/v2"
	"github.com/k3a/html2text"
	"go.uber.org/zap"
)

type (
	Web struct {
		Base
		param   *WebParameters
		scraper *colly.Collector
	}
	WebParameters struct {
		URL string `url:"url"`
	}
)

func (c *Web) Config(connector *model.Connector) (Connector, error) {
	c.Base.Config(connector)
	c.param = &WebParameters{}
	if err := connector.ConnectorSpecificConfig.ToStruct(c.param); err != nil {
		return nil, err
	}
	c.scraper = colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)
	return c, nil
}

func (c *Web) Execute(ctx context.Context, param model.JSONMap) error {
	c.scraper.OnHTML("a[href]", c.onNextURL)
	c.scraper.OnResponse(c.onResponse)
	if err := c.scraper.Visit(c.param.URL); err != nil {
		zap.L().Error("Failed to scrape URL", zap.String("url", c.param.URL), zap.Error(err))
	}
	return nil
}

func NewWeb(connector *model.Connector) (Connector, error) {
	var web Web
	return web.Config(connector)
}

func (c *Web) onNextURL(e *colly.HTMLElement) {
	if err := e.Request.Visit(e.Request.AbsoluteURL(e.Attr("href"))); err != nil {
		zap.S().Fatalw("Failed to visit URL", "url", e.Request.URL.String())
	}
}

func (c *Web) onResponse(e *colly.Response) {
	text := html2text.HTML2Text(string(e.Body))
	zap.S().Infof(text)
}

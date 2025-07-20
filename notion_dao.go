package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jomei/notionapi"
	"github.com/mmcdole/gofeed"
)

type NotionClient struct {
	client *notionapi.Client
}

func NewNotionClient(token string) *NotionClient {
	return &NotionClient{
		client: notionapi.NewClient(notionapi.Token(token)),
	}
}

// Obtener páginas que NO están marcadas con "Guardar"
func (n *NotionClient) GetPagesToDelete(databaseID string) ([]notionapi.PageID, error) {
	filter := notionapi.PropertyFilter{
		Property: "Guardar",
		Checkbox: &notionapi.CheckboxFilterCondition{
			Equals: false,
		},
	}
	query := &notionapi.DatabaseQueryRequest{
		Filter: &filter,
	}

	resp, err := n.client.Database.Query(context.Background(), notionapi.DatabaseID(databaseID), query)
	if err != nil {
		return nil, err
	}

	var ids []notionapi.PageID
	for _, page := range resp.Results {
		ids = append(ids, page.ID)
	}
	return ids, nil
}

// Archivar página
func (n *NotionClient) DeletePage(pageID notionapi.PageID) error {
	_, err := n.client.Page.Update(context.Background(), pageID, &notionapi.PageUpdateRequest{
		Archived: true,
	})
	return err
}

// Obtener URLs de feeds RSS desde la base de datos de feeds
func (n *NotionClient) GetFeedURLs(databaseID string) ([]string, error) {
	resp, err := n.client.Database.Query(context.Background(), notionapi.DatabaseID(databaseID), nil)
	if err != nil {
		return nil, err
	}

	var urls []string
	for _, page := range resp.Results {
		if prop, ok := page.Properties["URL"].(*notionapi.URLProperty); ok && prop.URL != "" {
			urls = append(urls, prop.URL)
		}
	}
	return urls, nil
}

// Crear página en Notion desde un ítem RSS
func (n *NotionClient) CreatePageFromRSSItem(databaseID string, item *gofeed.Item) error {
	props := notionapi.Properties{
		"Name": notionapi.TitleProperty{
			Title: []notionapi.RichText{
				{Text: &notionapi.Text{Content: item.Title}},
			},
		},
		"Link": notionapi.URLProperty{URL: item.Link},
	}

	if item.PublishedParsed != nil {
		props["Fecha"] = notionapi.DateProperty{
			Date: &notionapi.DateObject{
				Start: notionapi.DateTime{Time: *item.PublishedParsed},
			},
		}
	}

	_, err := n.client.Page.Create(context.Background(), &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			DatabaseID: notionapi.DatabaseID(databaseID),
		},
		Properties: props,
	})
	return err
}

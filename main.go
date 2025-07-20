package main

import (
	"log"
	"os"
)

func main() {
	// Cargar variables de entorno
	authToken := os.Getenv("NOTION_API_TOKEN")
	contentDB := os.Getenv("NOTION_RSS_CONTENT_DATABASE_ID")
	feedsDB := os.Getenv("NOTION_RSS_FEEDS_DATABASE_ID")

	if authToken == "" || contentDB == "" || feedsDB == "" {
		log.Fatal("Faltan variables de entorno necesarias.")
	}

	// Inicializar cliente de Notion
	notionClient := NewNotionClient(authToken)

	// Paso 1: Eliminar entradas no marcadas con "Guardar"
	pageIDs, err := notionClient.GetPagesToDelete(contentDB)
	if err != nil {
		log.Fatalf("Error obteniendo páginas a eliminar: %v", err)
	}
	for _, pageID := range pageIDs {
		if err := notionClient.DeletePage(pageID); err != nil {
			log.Printf("Error eliminando página %s: %v", pageID, err)
		}
	}

	// Paso 2: Obtener URLs de feeds RSS
	feedURLs, err := notionClient.GetFeedURLs(feedsDB)
	if err != nil {
		log.Fatalf("Error obteniendo feeds: %v", err)
	}

	// Paso 3: Leer e importar noticias
	for _, url := range feedURLs {
		items := ParseRSS(url)
		for _, item := range items {
			if err := notionClient.CreatePageFromRSSItem(contentDB, item); err != nil {
				log.Printf("Error creando página para %s: %v", item.Title, err)
			}
		}
	}
}

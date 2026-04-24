package main

import (
	"fmt"
	"log"

	"github.com/mmcdole/gofeed"
)

// Article представляет структуру статьи
type Article struct {
	Title   string
	Link    string
	Tags    []string
	Content string
}

func fetchRSSFeed(url string) ([]Article, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга RSS: %w", err)
	}

	var articles []Article
	for _, item := range feed.Items {
		// Извлекаем категории/теги из RSS
		var tags []string
		for _, category := range item.Categories {
			tags = append(tags, category)
		}

		articles = append(articles, Article{
			Title:   item.Title,
			Link:    item.Link,
			Tags:    tags,
			Content: item.Description,
		})
	}

	return articles, nil
}

func main() {
	rssURL := "https://habr.com/ru/rss/all/all/"

	articles, err := fetchRSSFeed(rssURL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Найдено статей: %d\n", len(articles))
	for i, article := range articles[:5] { // Покажем первые 5
		fmt.Printf("%d. %s\n%s\n", i+1, article.Title, article.Tags)
	}
}

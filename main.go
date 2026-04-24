package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mmcdole/gofeed"
)

type Article struct {
	Title   string
	Link    string
	Tags    []string
	Content string
}

func loadTagsFromFile(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config struct {
		TargetTags []string `json:"target_tags"`
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return config.TargetTags, nil
}

func fetchRSSFeedByTags(targetTags []string) ([]Article, error) {
	result := []Article{}
	if len(targetTags) == 0 {
		return result, fmt.Errorf("Тэги не заполнены")
	}

	for _, tag := range targetTags {
		time.Sleep(1 * time.Second)
		rssURL := "https://habr.com/ru/rss/hub/" + tag + "/articles/all/"
		log.Printf("Парсинг хаба: %s", tag)
		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(rssURL)
		if err != nil {
			log.Printf("Ошибка парсинга для тега %s: %v", tag, err)
			continue
		}
		for _, item := range feed.Items {
			// Извлекаем категории/теги из RSS
			var tags []string
			for _, category := range item.Categories {
				tags = append(tags, category)
			}

			result = append(result, Article{
				Title:   item.Title,
				Link:    item.Link,
				Tags:    tags,
				Content: item.Description,
			})
		}
	}

	return result, nil
}

func main() {
	targetTags, err := loadTagsFromFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	articles, err := fetchRSSFeedByTags(targetTags)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Найдено статей: %d\n", len(articles))

	allTags := []string{}
	for _, article := range articles {
		allTags = append(allTags, article.Tags...)
	}

	for i, article := range articles {
		log.Printf("%d. %s\n%s\n%s\n", i+1, article.Title, article.Link, article.Tags)
		log.Println("_______________________________________________________________")
	}
}

package main

import (
	"context"
	"fmt"
	"log"

	"io/ioutil"

	"time"

	"github.com/playwright-community/playwright-go"
	"github.com/vincent-petithory/dataurl"

	"github.com/mattn/go-mastodon"
)

func main() {
	c := mastodon.NewClient(&mastodon.Config{
		Server:       "https://PUT YOUR INSTANCE HERE",
		ClientID:     "PUT ID HERE",
		ClientSecret: "PUT SECRET HERE",
	})
	err := c.Authenticate(context.Background(), "PUT EMAIL HERE", "PUT PASSWORD HERE")
	if err != nil {
		log.Fatal(err)
	}
  //UNCOMMENT BELOW IF YOU WANT TO INSTALL PLAYWRIGHT DEPS (you need it)
	//err := playwright.Install(&playwright.RunOptions{Browsers: []string{"firefox"}})

	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	browser, err := pw.Firefox.Launch()
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	if _, err = page.Goto("https://place.zevent.fr/"); err != nil {
		log.Fatalf("could not goto: %v", err)
	}

	time.Sleep(5 * time.Second)
	for {
		if _, err = page.Goto("https://place.zevent.fr/"); err != nil {
			log.Fatalf("could not goto: %v", err)
		}

		time.Sleep(5 * time.Second)

		entry, err := page.QuerySelector(".game-container__inner")
		if err != nil {
			log.Fatalf("could not get entry: %v", err)
		}
		place, err := entry.QuerySelector("img")
		if err != nil {
			log.Fatalf("could not get place: %v", err)
		}

		fmt.Println(place)

		imgsrc, err := place.GetAttribute("src")
		if err != nil {
			log.Fatalf("could not get place src: %v", err)
		}

		//fmt.Println(imgsrc)

		imgsrcdecode, err := dataurl.DecodeString(imgsrc)
		if err != nil {
			log.Fatalf("could not decode place data url: %v", err)
		}

		ioutil.WriteFile("image.png", imgsrcdecode.Data, 0777)

		media, err := c.UploadMedia(context.Background(), "./image.png")
		c.PostStatus(context.Background(), &mastodon.Toot{MediaIDs: []mastodon.ID{media.ID}, Sensitive: true, Visibility: "unlisted"})

		time.Sleep(20 * time.Minute)
	}
}

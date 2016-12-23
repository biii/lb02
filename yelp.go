package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/JustinBeckwith/go-yelp/yelp"
	"github.com/guregu/null"
	"github.com/line/line-bot-sdk-go/linebot"
)

var o *yelp.AuthOptions
var food = make(map[string]string)

type UrlShortener struct {
	ShortUrl    string
	OriginalUrl string
}

func yelp_init() {
	rand.Seed(time.Now().UnixNano())

	// check environment variables
	o = &yelp.AuthOptions{
		ConsumerKey:       os.Getenv("CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("CONSUMER_SECRET"),
		AccessToken:       os.Getenv("ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
	}

	if o.ConsumerKey == "" || o.ConsumerSecret == "" || o.AccessToken == "" || o.AccessTokenSecret == "" {
		log.Fatal("Wrong environment setting about yelp-api-keys")
	}
}

func yelp_parse(bot *linebot.Client, token string, text string) {
	if _, err = bot.ReplyMessage(token, linebot.NewTextMessage(outmsg.String())).Do(); err != nil {
					log.Print(err)
				}
}
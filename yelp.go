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
	"bytes"

	"github.com/JustinBeckwith/go-yelp/yelp"
	"github.com/guregu/null"
	"github.com/line/line-bot-sdk-go/linebot"
)

var o *yelp.AuthOptions

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

func yelp_parse(bot *linebot.Client, token string, loc *linebot.LocationMessage, food string) {
	var err error
	
	// create a new yelp client with the auth keys
	client := yelp.New(o, nil)
	
	if loc == nil {
		// make a simple query for food and location
		results, err := client.DoSimpleSearch(food, "台北市大安區")
		if err == nil {
			log.Println("loc nil")
			yelp_parse_result(bot, token, results)
		}
	} else {
		// Build an advanced set of search criteria that include
		// general options, and coordinate options.
		s := yelp.SearchOptions{
			GeneralOptions: &yelp.GeneralOptions{
				Term: food,
			},
			//LocaleOptions: &yelp.LocaleOptions{
			//	lang: "zh",
			//},
			CoordinateOptions: &yelp.CoordinateOptions{
				Latitude:  null.FloatFrom(loc.Latitude),
				Longitude: null.FloatFrom(loc.Longitude),
			},
		}

		// Perform the search using the search options
		results, err := client.DoSearch(s)
		if err == nil {
			log.Println(fmt.Sprintf("loc:(%f, %f)", loc.Latitude, loc.Longitude))
			yelp_parse_result(bot, token, results)
		}
	}
	if err != nil {
		log.Println(err)
		_, err = bot.ReplyMessage(token, linebot.NewTextMessage("查無資料！")).Do()
	}	
}

func yelp_parse_result(bot *linebot.Client, token string, results yelp.SearchResult) {
	var err error
//	var msgs []linebot.Message
	var outmsg bytes.Buffer

	for i := 0; i < 3; i++ {
		//i := 0
		//if results.Total >= 20 {
		//	i = rand.Intn(20)
		//} else if results.Total >= 10 {
		//	i = rand.Intn(10)
		//} else if results.Total > j {
		//	i = j
		//} else if results.Total <= j && results.Total != 0 {
		if ( results.Total <= i) {
			//_, err = bot.ReplyMessage(token, linebot.NewTextMessage("已無更多資料！")).Do()
			break
		}
		urlOrig := UrlShortener{}
		urlOrig.short(results.Businesses[i].MobileURL)
		address := strings.Join(results.Businesses[i].Location.DisplayAddress, ",")
		//var largeImageURL = strings.Replace(results.Businesses[i].ImageURL, "ms.jpg", "l.jpg", 1)
		
		outmsg.WriteString("店名："+results.Businesses[i].Name+"\n電話："+results.Businesses[i].Phone+"\n評比："+strconv.FormatFloat(float64(results.Businesses[i].Rating), 'f', 1, 64)+"\n地址："+address+"\n更多資訊："+urlOrig.ShortUrl+"\n")

//		_, err = bot.ReplyMessage(token, linebot.NewImageMessage(largeImageURL, largeImageURL), linebot.NewTextMessage("店名："+results.Businesses[i].Name+"\n電話："+results.Businesses[i].Phone+"\n評比："+strconv.FormatFloat(float64(results.Businesses[i].Rating), 'f', 1, 64)+"\n更多資訊："+urlOrig.ShortUrl), linebot.NewLocationMessage(results.Businesses[i].Name+"\n", address, float64(results.Businesses[i].Location.Coordinate.Latitude), float64(results.Businesses[i].Location.Coordinate.Longitude)) ).Do()
//		msgs = append(msgs, linebot.NewImageMessage(largeImageURL, largeImageURL))
//		msgs = append(msgs, linebot.NewTextMessage("店名："+results.Businesses[i].Name+"\n電話："+results.Businesses[i].Phone+"\n評比："+strconv.FormatFloat(float64(results.Businesses[i].Rating), 'f', 1, 64)+"\n更多資訊："+urlOrig.ShortUrl))
//		msgs = append(msgs, linebot.NewLocationMessage(results.Businesses[i].Name+"\n", address, float64(results.Businesses[i].Location.Coordinate.Latitude), float64(results.Businesses[i].Location.Coordinate.Longitude)))
//		_, err = bot.ReplyMessage(token, linebot.NewImageMessage(largeImageURL, largeImageURL)).Do()
//		_, err = bot.ReplyMessage(token, linebot.NewTextMessage("店名："+results.Businesses[i].Name+"\n電話："+results.Businesses[i].Phone+"\n評比："+strconv.FormatFloat(float64(results.Businesses[i].Rating), 'f', 1, 64)+"\n更多資訊："+urlOrig.ShortUrl)).Do()
//		_, err = bot.ReplyMessage(token, linebot.NewLocationMessage(results.Businesses[i].Name+"\n", address, float64(results.Businesses[i].Location.Coordinate.Latitude), float64(results.Businesses[i].Location.Coordinate.Longitude))).Do()
	}
	
	if outmsg.String() != "" {
		_, err = bot.ReplyMessage(token, linebot.NewTextMessage(outmsg.String())).Do()
	}
	if err != nil {
		log.Println(err)
	}	
}

func getResponseData(urlOrig string) string {
	response, err := http.Get(urlOrig)
	if err != nil {
		log.Println(err)
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	return string(contents)
}

func isGdShortener(urlOrig string) (string, string) {
	escapedUrl := url.QueryEscape(urlOrig)
	isGdUrl := fmt.Sprintf("http://is.gd/create.php?url=%s&format=simple", escapedUrl)
	return getResponseData(isGdUrl), urlOrig
}

func (u *UrlShortener) short(urlOrig string) *UrlShortener {
	shortUrl, originalUrl := isGdShortener(urlOrig)
	u.ShortUrl = shortUrl
	u.OriginalUrl = originalUrl
	return u
}

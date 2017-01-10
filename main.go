// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"bytes"
	"math/rand"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client
var locmap = make(map[string]*linebot.LocationMessage)

func main() {
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	
	translate_init()
	yelp_init()
	
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				var outmsg bytes.Buffer

				switch {
					case strings.Compare(message.Text, "溫馨提醒") == 0:
						outmsg.WriteString("<<<溫馨提醒>>>\r\n因為這個群很吵 -->\r\n右上角 可以 關閉提醒\r\n\r\n[同學會] 投票進行中 -->\r\n右上角 筆記本 可以進行投票\r\n\r\n[通訊錄] 需要大家的協助 -->\r\n右上角 筆記本 請更新自己的聯絡方式")
					
					case strings.HasSuffix(message.Text, "麼帥"):
						outmsg.WriteString(GetHandsonText(message.Text))

					case strings.Compare(message.Text, "PPAP") == 0:
						outmsg.WriteString(GetPPAPText())

					case strings.HasPrefix(message.Text, "翻翻"):
						outmsg.WriteString(GetTransText(strings.TrimLeft(message.Text, "翻翻")))
						
					case strings.HasPrefix(message.Text, "吃吃"):
						yelp_parse(bot, event.ReplyToken, locmap[GetID(event.Source)], strings.TrimLeft(message.Text, "吃吃"))
						continue
					case strings.Compare(message.Text, "測試") == 0:
						outmsg.WriteString(message.ID)
						outmsg.WriteString("\r\n")
						outmsg.WriteString(event.Source.UserID)
						outmsg.WriteString("\r\n")
						outmsg.WriteString(event.Source.GroupID)
						outmsg.WriteString("\r\n")
						outmsg.WriteString(event.Source.RoomID)
						
					default:
						continue
				}
				
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(outmsg.String())).Do(); err != nil {
					log.Print(err)
				}
			case *linebot.LocationMessage:
				locmap[GetID(event.Source)] = message
			}
		}
	}
}

func GetID(source *linebot.EventSource) string {
	switch source.Type {
	case linebot.EventSourceTypeUser:
		return source.UserID
	case linebot.EventSourceTypeGroup:
		return source.GroupID
	case linebot.EventSourceTypeRoom:
		return source.RoomID
	}
	return source.UserID
}

func GetHandsonText(inText string) string {
	var outmsg bytes.Buffer	
	var outText bytes.Buffer
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(100)
	outmsg.WriteString("我覺得還是")
	switch i % 20 {
	case 0:
		outmsg.WriteString("我")
	case 1:
		outmsg.WriteString("你")
	default:
		outText.WriteString(inText)
		outText.WriteString("+1")
		return outText.String()
	}
	outmsg.WriteString("比較帥")
	return outmsg.String()	
}

func GetPPAPText() string {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(100)
	switch i % 5 {
	case 0:
		return "I have a pencil,\r\nI have an Apple,\r\nApple pencil.\r\nI have a watch,\r\nI have an Apple,\r\nApple watch."
	case 1:
		return "順帶一提，請不要把Apple Pencil刺進水果裡，不管是蘋果還是鳳梨。"
	case 2:
		return "我懂了，這是以書寫工具與種類食物為題的饒舌歌。"
	case 3:
		return "我不太清楚PPAP是什麼，但你可以問我AAPL的相關資訊。"
	case 4:
		return "我是不會接著唱的！"
	}
	return "去問 siri 啦"
}

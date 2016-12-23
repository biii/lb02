// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Translate uses the Google translate API from the command line to translate
// its arguments. By default it auto-detects the input language and translates
// to English.
// clone from "robpike/translate"

package main 

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"bytes"
	"strings"
	"os"
)

var key string

type Response struct {
	Data struct {
		Translations []Translation
	}
}

type Translation struct {
	TranslatedText         string
	DetectedSourceLanguage string
}

func translate_init() {
	key = os.Getenv("GOOGLEAPIKEY")
}

func GetTransText(text string) string {
	var target = "en"
	var intext = text
	switch {
		case strings.HasPrefix(text, "中"):
			target = "zh-TW"
			intext = strings.TrimLeft(text, "中")
		case strings.HasPrefix(text, "日"):
			target = "ja"
			intext = strings.TrimLeft(text, "日")
		case strings.HasPrefix(text, "法"):
			target = "fr"
			intext = strings.TrimLeft(text, "法")
		case strings.HasPrefix(text, "韓"):
			target = "ko"
			intext = strings.TrimLeft(text, "韓")
		case strings.HasPrefix(text, "英"):
			intext = strings.TrimLeft(text, "英")
	}
	return DoTrans(target, intext);
}

func DoTrans(target string, text string) string {
	var outmsg bytes.Buffer
	
	v := make(url.Values)
	v.Set("key", key)
	v.Set("target", target)
	v.Set("q", text)
	url := "https://www.googleapis.com/language/translate/v2?" + v.Encode()
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	var r Response
	if err := json.Unmarshal(data, &r); err != nil {
		log.Fatal(err)
	}
	for _, t := range r.Data.Translations {
		outmsg.WriteString(fmt.Sprintf("%s (%s)\n", html.UnescapeString(t.TranslatedText), t.DetectedSourceLanguage))
	}
	return outmsg.String()
}

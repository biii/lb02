package main

import (
	"fmt"
	"strings"
	"math/rand"
	"time"
)

const COMP_PREFIX = "問問"
const COMP_MORE = "誰比較"

var nPrefix = len(COMP_PREFIX)
var nMore = len(COMP_MORE)

var msgWrong = []string {
	"別問我，去問Google大神\r\n你只能問我 A，B，C 誰比較有錢",
	"來點會的好嗎?\r\n問問 彬彬，金城武，BigBang T.O.P 誰比較帥",
	"再這樣我報警了喔!! 我來教教你\r\n問問 AlphaGo,BetaGo 誰比較會Go",
}
var msgEmpty = []string {
	"跟鬼比的話，你是有比較%s",
	"你的程度只到這裡嗎?\r\n沒事跟空氣比%s幹嘛",
	"再這樣我報警了喔!!\r\n別跟我比%s，你是不可能會贏的\r\n我是AlphaGo II耶",
}
var msgSingle = []string {
	"小學有畢業嗎?\r\n只有[%s]，你是要比三小%s!?",
	"叫你讀書不讀書，語法懂嗎?\r\n[%s]誰比較%s\r\n你覺得這樣是對的嗎?",
	"再這樣我報警了喔!!\r\n[%s]誰比較%s，這樣你也敢問?",
}

func CompareGetString(objs []string) string {
	var seed = len(objs)
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(seed) 
	
	return objs[i]
}

func CompareNotSupport() string {
	return CompareGetString(msgWrong)
}

func CompareEmptyObject(comp string) string {
	return fmt.Sprintf(CompareGetString(msgEmpty), comp)
}

func CompareSingleObject(obj, comp string) string {
	return fmt.Sprintf(CompareGetString(msgSingle), obj, comp)
}

func CompareObjects(objs []string, comp string) string {
	var seed = len(objs)
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(seed + 1) 
	
	if i < seed {
		return fmt.Sprintf("根據AlphaGo II的計算結果，%s比較%s", objs[i], comp)
	}
	
	return fmt.Sprintf("不用爭了, 這%d個一樣%s", seed, comp)
}

func CompareSplitObjects(object string) []string {
	objs := strings.Split(object, ",")
	if len(objs) > 1 {
		return objs
	}
	
	return strings.Split(object, "，")
}

func CompareCheckTokens(input string) string {
	var idx = strings.LastIndex(input, COMP_MORE)
	if idx <= 0 {
		return CompareNotSupport()
	}
	
	comp := input[idx+nMore:]
	object := strings.TrimSpace(input[nPrefix:idx])
	if len(object) == 0 {
		return CompareEmptyObject(comp)
	}
	
	objs := CompareSplitObjects(object)
	if len(objs) == 1 {
		return CompareSingleObject(object, comp)
	}
	
	return CompareObjects(objs, comp)
}

func CompareCheckTokens2(input string) (string, bool) {
	var result = ""
	if !strings.HasPrefix(input, COMP_PREFIX) {
		return result, false
	}
	var idx = strings.LastIndex(input, COMP_MORE)
	if idx <= 0 {
		return CompareNotSupport(), true
	}
	
	comp := input[idx+nMore:]
	object := strings.TrimSpace(input[nPrefix:idx])
	if len(object) == 0 {
		return CompareEmptyObject(comp), true
	}
	
	objs := CompareSplitObjects(object)
	if len(objs) == 1 {
		return CompareSingleObject(object, comp), true
	}
	
	return CompareObjects(objs, comp), true
}

/*
// for test
func main() {
	var result string
	var err bool

	if result, err = CompareCheckTokens2("問問  123"); err {
		fmt.Println(result)
	}

	result, err = CompareCheckTokens2("問問  誰比較帥")
	fmt.Println(result)

	result, err = CompareCheckTokens2("問問 123;456 誰比較帥")
	fmt.Println(result)
	
	result = CompareCheckTokens("問問 123,456 誰比較帥")
	fmt.Println(result)

}
*/
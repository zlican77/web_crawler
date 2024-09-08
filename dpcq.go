//https://www.westnovel.com/wuxia/dpcq/137814.html
//拿到每一个title和cont
//<div class="inner" id="BookCon"> <h1>(.*?)</h1>
//<div id="BookText" style>(.*?)<div class="ads">
//合并写入一个txt

package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func getEachEssay(url string, i int, titlesAndConts []string, result chan int) { //拿到每一个文章的title和conts

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("getUrlerr", err)
		result <- i
		return
	}
	defer resp.Body.Close()

	buf := make([]byte, 4096)
	var cont string
	for {
		n, _ := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		cont += string(buf[:n]) //拿到每个essay中的所有html
	}

	reT := regexp.MustCompile(`<h1>(.*?)</h1>`)
	matchesT := reT.FindAllStringSubmatch(cont, -1)
	if len(matchesT) == 0 {
		titlesAndConts[i] = ""
		result <- i
		return
	}
	NewTitle := matchesT[0][1]
	NewTitle = strings.Replace(NewTitle, "[", "", -1)
	NewTitle = strings.Replace(NewTitle, "]", "", -1)
	NewTitle = strings.Replace(NewTitle, " ", "", -1)
	titlesAndConts[i] = NewTitle

	reC := regexp.MustCompile(`(?s)<div id="BookText" style>(.*?)<div class="ads">`)
	matchesC := reC.FindAllStringSubmatch(cont, -1)
	if len(matchesC) == 0 {
		titlesAndConts[i] = ""
		result <- i
		return
	}
	NewCont := matchesC[0][1]
	NewCont = strings.Replace(NewCont, "<p>", "", -1)
	NewCont = strings.Replace(NewCont, "</p>", "", -1)
	NewCont = strings.Replace(NewCont, "<br/>", "", -1)
	NewCont = strings.Replace(NewCont, "<h2>", "", -1)
	NewCont = strings.Replace(NewCont, "</h2>", "", -1)
	NewCont = strings.Replace(NewCont, "[", "", -1)
	NewCont = strings.Replace(NewCont, "]", "", -1)
	NewCont = strings.Replace(NewCont, " ", "", -1)
	titlesAndConts[i] += "\r\n======================" + NewCont + "\n\n\n\n\n"

	result <- i
	return
}

func DoWork(start int, end int) {
	var titlesAndConts = make([]string, 1624)
	var result = make(chan int)
	for i := start; i <= end; i++ {
		go getEachEssay("https://www.westnovel.com/wuxia/dpcq/"+strconv.Itoa(137813+i)+".html", i, titlesAndConts, result)
	}

	for i := start; i <= end; i++ {
		<-result
	}

	titlesAndContsStr := strings.Join(titlesAndConts, "\r\n")
	fileName := "斗破苍穹.txt"
	f, err := os.Create(fileName)
	if err != nil {
		return
	}

	f.WriteString(titlesAndContsStr)
}

func main() {
	var start, end int //控制爬取范围
	fmt.Printf("请输入起始章（ >= 1）：")
	fmt.Scan(&start)
	fmt.Printf("请输入终止章（ >= 起始章）：")
	fmt.Scan(&end)

	DoWork(start, end)

}

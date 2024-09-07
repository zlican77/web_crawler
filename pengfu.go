//1.启动爬虫目标页面
//https://www.pengfu.net/xiaohuaduanzi/index_2.html

//2.爬取该目标页面中的essay链接
//<h3 class="blogtitle"><a href="  链接  "  target="_blank"> 	25个

//在对应的每个essay链接中爬取title和content
//<h1 class="con_tilte"> title </h1>
//<div class="con_text"> content  <p class="share">

//将title + content 在essay 中返回
//将essay 返回的内容拼接 再返回到单页面

//将单页面返回的内容写入text

package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func GetEssayCont(linkUrl string, titlesAndConts chan string) { //处理单文章的title和content
	resp, err := http.Get(linkUrl)
	if err != nil {
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
		cont += string(buf[:n]) //拿到essay中的所有html

	}

	var titleAndCont string
	reT := regexp.MustCompile(`<h1 class="con_tilte">(.*?)</h1>`)
	matchesT := reT.FindAllStringSubmatch(cont, -1)
	if len(matchesT) == 0 {
		titlesAndConts <- ""
		return
	}
	NewTitle := matchesT[0][1]
	NewTitle = strings.Replace(NewTitle, "[", "", -1)
	NewTitle = strings.Replace(NewTitle, "]", "", -1)
	NewTitle = strings.Replace(NewTitle, " ", "", -1)
	titleAndCont = NewTitle

	reC := regexp.MustCompile(`(?s)<div class="con_text">(.*?)<p class="share">`)
	matchesC := reC.FindAllStringSubmatch(cont, -1)
	if len(matchesC) == 0 {
		titlesAndConts <- ""
		return
	}
	NewCont := matchesC[0][1]
	NewCont = strings.Replace(NewCont, "<p>", "", -1)
	NewCont = strings.Replace(NewCont, "</p>", "", -1)
	NewCont = strings.Replace(NewCont, "<h2>", "", -1)
	NewCont = strings.Replace(NewCont, "</h2>", "", -1)
	NewCont = strings.Replace(NewCont, "[", "", -1)
	NewCont = strings.Replace(NewCont, "]", "", -1)
	NewCont = strings.Replace(NewCont, " ", "", -1)
	titleAndCont += "\r\n======================" + NewCont + "\n\n\n\n\n"

	titlesAndConts <- titleAndCont
	return
}

func GetPageCont(url string, titlesAndConts chan string) (cont string, err error) { //负责拿到页面上的初始文章链接
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	buf := make([]byte, 4096)
	for {
		n, _ := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		cont += string(buf[:n])
	}

	re := regexp.MustCompile(`<h3 class="blogtitle"><a href="([^"]+)" target="_blank">`)
	matches := re.FindAllStringSubmatch(cont, -1)

	for _, link := range matches {
		go GetEssayCont(link[1], titlesAndConts)
	}

	return
}

func GetPage(url string, i int, pages chan int) { //处理单页面逻辑
	var titlesAndConts = make(chan string, 25)
	_, err := GetPageCont(url, titlesAndConts) //拿到一整个页面的所有title和content
	if err != nil {
		return
	}

	var cont string
	for i := 0; i < 25; i++ {
		cont += <-titlesAndConts
	}

	fileName := strconv.Itoa(i) + ".txt"
	f, err2 := os.Create(fileName)
	if err2 != nil {
		return
	}
	f.WriteString(cont)

	pages <- i
}

func DoWork(start int, end int) {
	var pages = make(chan int)
	for i := start; i <= end; i++ {
		go GetPage("https://www.pengfu.net/xiaohuaduanzi/index_"+strconv.Itoa(i+1)+".html", i, pages)
	}

	for i := start; i <= end; i++ {
		<-pages
	}

}

func main() {
	var start, end int //控制爬取范围
	fmt.Printf("请输入起始页（ >= 1）：")
	fmt.Scan(&start)
	fmt.Printf("请输入终止页（ >= 起始页）：")
	fmt.Scan(&end)

	DoWork(start, end)

}

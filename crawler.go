package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func HttpGet(url string) (cont string, err2 error) {
	resp, err := http.Get(url)
	if err != nil {
		err2 = err
		return
	}

	defer resp.Body.Close()

	buf := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(buf)
		if n == 0 {
			fmt.Println(err)
			break
		}
		cont += string(buf[:n])
	}
	return
}

func DoWork(start int, end int) {
	//https://tieba.baidu.com/f?kw=%E5%8C%BA%E5%9D%97%E9%93%BE&ie=utf-8&pn=50

	for i := start; i <= end; i++ {
		//在每次for循环中爬取单个网页
		go func() { //因为此处在for过程中，操作不能一次性立刻结束，需要引入时间维度等待http响应，故因采用go routine
			cont, err := HttpGet("https://tieba.baidu.com/f?kw=%E5%8C%BA%E5%9D%97%E9%93%BE&ie=utf-8&pn=" +
				strconv.Itoa(i*50))

			if err != nil {
				return
			}

			//把内容写入文件
			fileName := strconv.Itoa(i) + ".html"
			f, err1 := os.Create(fileName)
			if err1 != nil {
				fmt.Println(err1)
				return
			}
			f.WriteString(cont)
			f.Close()
		}()
	}
}

func main() {
	var start, end int //控制爬取范围
	fmt.Printf("请输入起始页（ >= 1）：")
	fmt.Scan(&start)
	fmt.Printf("请输入终止页（ >= 起始页）：")
	fmt.Scan(&end)

	DoWork(start, end)
	for {

	}
}

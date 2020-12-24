package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/viper"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"sync"
)

var wg sync.WaitGroup

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	config := viper.AllSettings()
	for name, password := range config {
		wg.Add(1)
		go login(name, password.(string))
	}
	wg.Wait()

}

func login(name, password string) {
	loginUrl := "https://www.hostloc.com/member.php?mod=logging&action=login&loginsubmit=yes&infloat=yes&lssubmit=yes&inajax=1"
	creditUrl := "https://www.hostloc.com/home.php?mod=spacecp&ac=credit"
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}
	_, err := client.PostForm(loginUrl, url.Values{"username": {name}, "password": {password}})
	if err != nil {
		panic(err)
	}

	resp, err := client.Post(creditUrl)
	if err != nil {
		panic(err)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 16; i++ {
		wg.Add(1)
		rand.Seed(int64(i))
		id := rand.Intn(60000)
		go viewSpace(client, id)
	}

	doc.Find(".creditl.mtm.bbda.cl").Each(func(i int, s *goquery.Selection) {
		money := s.Find("li").Eq(0).Text()
		credit := s.Find("li").Eq(1).Text()
		integral := s.Find("li").Eq(2).Text()
		abs := s.Find("li").Eq(2).Text()
		addr := s.Find("li").Eq(2).Text()
		article := s.Find("article").Eq(2).Text()
		log.Println(name, money, credit, integral, abs, addr)
	})
	wg.Done()
}

func viewSpace(client *http.Client, id int) {
	spaceUrl := "https://www.hostloc.com/space-uid-" + strconv.Itoa(id) + ".html"
	_, err := client.Get(spaceUrl)
	if err != nil {
		panic(err)
	}
	wg.Done()
}

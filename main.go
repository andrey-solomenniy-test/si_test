package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/andrey-solomenniy-test/si_test/config"
	"github.com/andrey-solomenniy-test/si_test/ipcache"
)

const (
	configFileName = "si_test.yaml"
)

var conf *config.Config
var CountryNameByCode map[[2]byte]string // карта для хранения названий страны по её коду ISO
var Cache *ipcache.Cache                 // кэш, в котором будут храниться ip-адреса

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	CountryNameByCode = make(map[[2]byte]string, 200)

	var err error
	conf, err = config.LoadConfig("./config.json")
	if err != nil {
		log.Println("Error reading config file:", err)
		return
	}
	//for _, v := range conf.ServerList {
	//	log.Println(v.URL)
	//}
	//log.Println()

	Cache = ipcache.NewCache(conf.Cache.Size, conf.Cache.ValidExpired)

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	// логирование в случае паники
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("%v\n%s", err, buf)
		}
	}()

	ip := getIPFromRequest(r) // получаем ip из заголовка запроса

	var code [2]byte
	var country string
	ok, code := Cache.Find(ip) // ищем ip-адрес в кэше
	if ok {                    // если нашли, то берём данные из кэша
		//log.Println("found in cache")
		country, ok = CountryNameByCode[code]
		if !ok {
			log.Println("Strange error: ")
		}
	} else { // если не нашли, то идём за данными на один из серверов
		//log.Println("not found in cache, go to server")
		code, country = getCountryByIP(ip)
		Cache.Save(ip, code) // после чего сохраняем полученные данные в  кэше
		// получаем название страны по её коду ISO
		if _, ok = CountryNameByCode[code]; !ok {
			CountryNameByCode[code] = country
		}
	}

	fmt.Fprintf(w, "ip = <%s>, country = <%s>", ip, country)
}

func getCountryByIP(ip string) ([2]byte, string) {
	// выбираем один из серверов
	curServer, err := getCurServer()
	if err != nil {
		return [2]byte{}, err.Error()
	}
	// выполняем запрос
	resp, err := http.Get(fmt.Sprintf(curServer.URL, ip))
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return [2]byte{}, "Cannot get URL"
	}
	if resp.StatusCode >= 400 {
		return [2]byte{}, "Request error: " + resp.Status
	}
	// копируем JSON в буфер
	buf := bytes.Buffer{}
	io.Copy(&buf, resp.Body)
	// получаем код
	code := getJsonField(buf.Bytes(), curServer.CodeField)
	b := [2]byte{}
	if len(code) >= 2 {
		b[0] = code[0]
		b[1] = code[1]
	}
	// возвращаем код и название страны
	return b, getJsonField(buf.Bytes(), curServer.NameField)
}

package config

import (
	"encoding/json"
	"os"
	"time"
)

type Config struct {
	TimeUnit   int      `json:"timeUnit"`   // количество минут в единице времени
	Cache      Cache    `json:"cache"`      // настройки для кеша
	ServerList []Server `json:"ServerList"` // список серверов
}

type Cache struct {
	Size         int `json:"size"`
	ValidExpired int `json:"exp"`
}

type Server struct {
	URL             string    `json:"URL"`       // форматная строка с одним параметром, куда подставляется ip-адрес
	Freq            int       `json:"freq"`      // разрешенное количество обращений к серверу за единицу времени
	CodeField       string    `json:"codeField"` // имя поля в json для кода страны, например code или country.iso
	NameField       string    `json:"nameField"` // имя поля в json для наименования страны по-английски, например country_name или country.name.english
	TimePeriodBegin time.Time // начало периода времени после последнего сброса счетчика
	CountForPeriod  int       // счетчик, количество обращений к серверу с момента последнего сброса
}

func LoadConfig(fileName string) (*Config, error) {
	// Читаем конфиг из файла
	config := &Config{}
	configFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	// Парсим json в структуру
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(config)
	if err != nil {
		return nil, err
	}

	// Инициализируем поля в списке серверов
	// начало периода текущим временем,
	// счётчик обращений к серверу за период времени нулём
	t := time.Now()
	for i := range config.ServerList {
		config.ServerList[i].TimePeriodBegin = t
		config.ServerList[i].CountForPeriod = 0
	}

	return config, nil
}

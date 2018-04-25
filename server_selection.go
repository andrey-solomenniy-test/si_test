package main

import (
	"errors"
	"time"

	"github.com/andrey-solomenniy-test/si_test/config"
)

// выбираем один сервер из списка
func getCurServer() (*config.Server, error) {
	found := false
	now := time.Now()
	var curServer *config.Server
	// в цикле по списку серверов начиная с первого проверяем не исчерпалось ли количество запросов в единицу времени
	// если все сервера исчерпали время, то возвращаем ошибку
	for i := range conf.ServerList {
		curServer = &conf.ServerList[i]
		diff := now.Sub(curServer.TimePeriodBegin)
		if diff >= (time.Duration(conf.TimeUnit) * time.Minute) {
			curServer.TimePeriodBegin = now
			curServer.CountForPeriod = 0
		}
		if curServer.CountForPeriod < curServer.Freq {
			found = true
			break
		}
	}
	if !found {
		return nil, errors.New("All servers expired time limits. Wait a minute and then try again.")
	}
	curServer.CountForPeriod++
	return curServer, nil
}

package main

import (
	"github.com/sirupsen/logrus.git"
	"time"
)

func main() {

	log := logrus.New()
	logrus.New()
	format := &logrus.TextFormatter{}
	format.TimestampFormat = "2006-01-02 15:04:05"
	format.ForceColors = true
	log.Formatter = format


	log.Infof("%s","hello")
	log.Infoln("a","b")

	log.WithFields(logrus.Fields{"animal":"dog","dog":"tom"}).Info("haha")

	log.Warnln("warning ......")


	go func() {
		log.Errorln("err.....")
	}()



	time.Sleep(100*time.Millisecond)


}

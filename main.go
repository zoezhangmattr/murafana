package main

import (
	"flag"
	"fmt"
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/superwomany/murafana/services"
	"github.com/superwomany/murafana/utils"
)

func main() {
	uid := flag.String("uid", "", "download specific dashboard uid to yaml file")
	dlist := flag.String("download-list", "", "download dashboard to yaml file if enabled")
	flag.Parse()
	if len(*uid) > 0 {
		c := services.New()
		err := c.GetDashboardByUID(*uid)
		if err != nil {
			logger.Error(err)
		}
		return
	}
	if len(*dlist) > 0 {
		err := DownloadDashboardList()
		if err != nil {
			logger.Error(err)
		}
		return
	}

	err := DownloadAllDashboards()
	if err != nil {
		logger.Error(err)
	}
}

func DownloadAllDashboards() error {
	c := services.New()
	uids, err := c.ListDashboards()
	if err != nil {
		logger.Error(err)
		return err
	}
	var pwg sync.WaitGroup
	pwg.Add(len(uids))
	errChan := make(chan error, len(uids))
	logger.Info(uids, len(uids))
	for _, v := range uids {
		go func(c services.ServiceIntf, uid string) {
			defer pwg.Done()

			err := c.GetDashboardByUID(uid)
			if err != nil {
				nerr := fmt.Errorf("uid: %v, error: %v", uid, err)
				errChan <- nerr
				return
			}
		}(c, v)
	}
	pwg.Wait()
	close(errChan)
	if err := utils.ParseErrorsFromChannel(errChan); err != nil {
		logger.Errorf("we have errors %v, but it is ok to continue.", err)
		// return err
	}
	return nil
}

func DownloadDashboardList() error {
	c := services.New()
	_, err := c.ListDashboards("enabled")
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

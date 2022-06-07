package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	logger "github.com/sirupsen/logrus"
	"github.com/superwomany/murafana/services"
	"github.com/superwomany/murafana/utils"
)

func main() {
	mode := flag.String("mode", "", "what task you want to run")
	uid := flag.String("uid", "", "download specific dashboard uid to yaml file")
	jfp := flag.String("file-path", "", "dashboard json file path")
	dn := flag.String("name", "", "dashboard name used for import")
	flag.Parse()
	c := services.New()
	switch *mode {
	case "download-dashboard":
		if len(*uid) > 0 {
			err := c.GetDashboardByUID(*uid)
			if err != nil {
				logger.Error(err)
			}
			return
		} else {
			logger.Error("uid is required")
		}
	case "download-list":
		err := DownloadDashboardList()
		if err != nil {
			logger.Error(err)
		}
	case "download-dashboards":
		err := DownloadAllDashboards()
		if err != nil {
			logger.Error(err)
		}
	case "import-dashboard":
		if len(*jfp) == 0 {
			logger.Error("dashboard json file path is required")
		}
		d, err := os.Open(*jfp)
		if err != nil {
			logger.Error("read file data failed")
		}
		defer d.Close()
		db, err := ioutil.ReadAll(d)
		if err != nil {
			logger.Error("read bytes failed")
		}
		err = c.ImportDashboardFromJson(*dn, db)
		if err != nil {
			logger.Error("import failed")
		}
	default:
		logger.Error("unknown mode")
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

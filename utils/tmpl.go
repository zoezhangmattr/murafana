package utils

import (
	"html/template"
	"os"

	logger "github.com/sirupsen/logrus"
)

func PrintOutDashboardObject(dl interface{}, tmpl string) error {
	tpl, err := template.New("").Parse(tmpl)
	if err != nil {
		logger.Errorf("parse template failure %v", err)
		return err
	}
	err = tpl.Execute(os.Stdout, dl)
	if err != nil {
		logger.Errorf("render template failure %v", err)
		return err
	}
	return nil
}

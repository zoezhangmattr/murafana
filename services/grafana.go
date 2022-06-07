package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"github.com/superwomany/murafana/utils"
)

// HTTPDoer describes something that performs an http request.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}
type Config struct {
	GrafanaHost   string
	GrafanaApiKey string
}
type client struct {
	config *Config
	doer   HTTPDoer
}

type ServiceIntf interface {
	ListDashboards(...string) ([]string, error)
	GetDashboardByUID(string) error
	ImportDashboardFromJson(string, []byte) error
}

func New() ServiceIntf {
	c := &Config{
		GrafanaHost:   os.Getenv("GRAFANA_URL"),
		GrafanaApiKey: os.Getenv("GRAFANA_CLOUD_API_KEY"),
	}
	return &client{
		config: c,
		doer:   &http.Client{},
	}
}

func NewWith(c *Config, h HTTPDoer) ServiceIntf {
	return &client{
		config: c,
		doer:   h,
	}
}

func (c *client) GetDashboardByUID(uid string) error {
	uri := c.config.GrafanaHost + "/api/dashboards/uid/" + uid
	req, err := http.NewRequestWithContext(context.Background(), "GET", uri, nil)

	result := DashboardMeta{}
	if err != nil {
		logger.WithFields(logger.Fields{
			"url":   uri,
			"error": err,
		}).Errorf("request failed.")
		return err
	}
	err = c.DoRequest(c.doer, req, &result)
	if err != nil {
		logger.WithFields(logger.Fields{
			"url":   uri,
			"error": err,
		}).Errorf("request failed.")
		return err
	}
	logger.WithFields(logger.Fields{
		"url": uri,
	}).Info("request succeed")
	utils.LoadToJson(result.Dashboard, "dashboard-"+uid+".json")
	return nil
}

// ListDashboards retrieves a list of dashboard and return a list of uid
func (c *client) ListDashboards(download ...string) ([]string, error) {

	uri := c.config.GrafanaHost + "/api/search?query=%"

	req, err := http.NewRequestWithContext(context.Background(), "GET", uri, nil)
	if err != nil {
		logger.WithFields(logger.Fields{
			"url":   uri,
			"error": err,
		}).Errorf("request failed.")
		return nil, err
	}

	result := []DashboardObject{}
	err = c.DoRequest(c.doer, req, &result)
	if err != nil {
		logger.WithFields(logger.Fields{
			"url":   uri,
			"error": err,
		}).Errorf("request failed.")
		return nil, err
	}
	logger.WithFields(logger.Fields{
		"url":   uri,
		"count": len(result),
	}).Info("request succeed.")
	if len(download) > 0 && download[0] == "enabled" {
		utils.LoadToYaml(result, "dashboards.yaml")
	}

	uids := []string{}
	for _, res := range result {
		uids = append(uids, res.UID)
	}

	return uids, nil
}

// ImportDashboardFromJson use json data to import dashboard
func (c *client) ImportDashboardFromJson(name string, data []byte) error {
	uri := c.config.GrafanaHost + "/api/dashboards/import"

	ds := Dashboard{}
	err := json.Unmarshal(data, &ds)
	// customize dashboard name
	if len(name) > 0 {
		ds.Title = name
	}
	// we want to import so set id to null
	ds.Id = nil
	if err != nil {
		logger.Error("unmarshal failed")
		return err
	}

	payload := DashboardImportBody{
		Dashboard: ds,
	}

	pm, err := json.Marshal(payload)
	if err != nil {
		logger.Error("payload marshal failed")
		return err
	}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, uri, bytes.NewBuffer(pm))
	if err != nil {
		logger.WithFields(logger.Fields{
			"url":   uri,
			"error": err,
		}).Errorf("request failed.")
		return err
	}

	var result interface{}
	err = c.DoRequest(c.doer, req, &result)
	if err != nil {
		logger.WithFields(logger.Fields{
			"url":   uri,
			"error": err,
		}).Errorf("request failed.")
		return err
	}
	logger.WithFields(logger.Fields{
		"url": uri,
	}).Info("request succeed.")

	return nil

}

// DoRequest sends http request
func (c *client) DoRequest(doer HTTPDoer, req *http.Request, dst interface{}) error {
	req.Header.Set("Authorization", "Bearer "+c.config.GrafanaApiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := doer.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-200 OK status code: %v response body: %q", resp.Status, raw)
	}

	if err := json.Unmarshal(raw, dst); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

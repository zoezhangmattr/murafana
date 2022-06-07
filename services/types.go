package services

type DashboardObject struct {
	ID            int32    `json:"id" yaml:"id"`
	UID           string   `json:"uid" yaml:"uid"`
	Title         string   `json:"title" yaml:"title"`
	Uri           string   `json:"uri" yaml:"uri"`
	Url           string   `json:"url" yaml:"url"`
	Slug          string   `json:"slug" yaml:"slug"`
	DashboardType string   `json:"type" yaml:"type"`
	Tags          []string `json:"tags" yaml:"tags"`
	IsStarred     bool     `json:"isStarred" yaml:"isStarred"`
	SortMeta      int32    `json:"sortMeta" yaml:"sortMeta"`
}

type DashboardMeta struct {
	Dashboard Dashboard              `json:"dashboard"`
	Meta      map[string]interface{} `json:"meta"`
}

type Dashboard struct {
	Id                   *int64                 `json:"id"`
	Uid                  string                 `json:"uid"`
	Title                string                 `json:"title"`
	Tags                 []string               `json:"tags"`
	Style                string                 `json:"style"`
	Timezone             string                 `json:"timezone"`
	Time                 map[string]interface{} `json:"time"`
	Editable             bool                   `json:"editable"`
	Annotations          map[string]interface{} `json:"annotations"`
	Refresh              string                 `json:"refresh"`
	Version              int64                  `json:"version"`
	SchemaVersion        int64                  `json:"schemaVersion"`
	Panels               []interface{}          `json:"panels"`
	Templating           map[string]interface{} `json:"templating"`
	Timepicker           map[string]interface{} `json:"timepicker"`
	GraphTooltip         int64                  `json:"graphTooltip"`
	WeekStart            string                 `json:"weekStart"`
	Description          string                 `json:"description"`
	GnetId               int64                  `json:"gnetId"`
	Links                []interface{}          `json:"links"`
	LiveNow              bool                   `json:"liveNow"`
	FiscalYearStartMonth int64                  `json:"fiscalYearStartMonth"`
	Iteration            int64                  `json:"iteration"`
}

type DashboardImportBody struct {
	Dashboard Dashboard      `json:"dashboard"`
	Overwrite bool           `json:"overwrite"`
	Inputs    []ImportInputs `json:"inputs"`
}

type ImportInputs struct {
	Type     string `json:"type"`
	PluginId string `json:"pluginId"`
	Name     string `json:"name"`
	Value    string `json:"value"`
}

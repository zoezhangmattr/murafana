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

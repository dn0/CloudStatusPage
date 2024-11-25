package data

import (
	"cspage/pkg/db"
)

//nolint:gochecknoglobals // These are constants.
var (
	cloudAWS   = &Cloud{Id: "aws", Name: "AWS", Color: "#9333ea", Symbol: "rect"}
	cloudAzure = &Cloud{Id: "azure", Name: "Azure", Color: "#06b6d4", Symbol: "triangle"}
	cloudGCP   = &Cloud{Id: "gcp", Name: "GCP", Color: "#6366f1", Symbol: "diamond"}

	Clouds = []*Cloud{
		cloudAWS,
		cloudAzure,
		cloudGCP,
	}

	CloudIds = []string{
		cloudAWS.Id,
		cloudAzure.Id,
		cloudGCP.Id,
	}

	CloudMap = map[string]*Cloud{
		cloudAWS.Id:   cloudAWS,
		cloudAzure.Id: cloudAzure,
		cloudGCP.Id:   cloudGCP,
	}
)

type Cloud struct {
	Id     string
	Name   string
	Color  string
	Symbol string
}

func (c *Cloud) URLPrefix() string {
	if c == nil {
		return ""
	}
	return "/cloud/" + c.Id
}

func GetCloud(id string) (*Cloud, error) {
	if cloud, ok := CloudMap[id]; ok {
		return cloud, nil
	}
	return nil, &db.ObjectNotFoundError{Object: "cloud=" + id}
}

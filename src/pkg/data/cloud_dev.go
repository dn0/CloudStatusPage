//go:build dev

package data

var (
	cloudDummy = &Cloud{Id: "dummy", Name: "ABC", Color: "#00b361", Symbol: "circle"}
)

func init() {
	Clouds = append(Clouds, cloudDummy)
	CloudIds = append(CloudIds, cloudDummy.Id)
	CloudMap[cloudDummy.Id] = cloudDummy
}

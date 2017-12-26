package model

const (
	SYSTEM_RAN_UNIT_TESTS = "RanUnitTests"
)

type System struct {
	Name  string `bson:"_id" json:"name"`
	Value string `bson:"value" json:"value"`
}

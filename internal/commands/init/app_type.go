package generator

type AppType string

const (
	Web AppType = "web"
	API AppType = "api"
)

var AllAppTypes = []AppType{
	Web,
	API,
}

func (t AppType) IsValid() bool {
	for _, appType := range AllAppTypes {
		if appType == t {
			return true
		}
	}
	return false
}

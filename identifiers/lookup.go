package identifiers

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed guids.json
var guidDatabaseData []byte

var guidDatabase map[string]string

func init() {
	if err := json.Unmarshal(guidDatabaseData, &guidDatabase); err != nil {
		panic(fmt.Errorf("unable to load GUID database: %w", err))
	}
}

func LookupGUID(g GUID) (string, bool) {
	description, ok := guidDatabase[g.String()]
	return description, ok
}

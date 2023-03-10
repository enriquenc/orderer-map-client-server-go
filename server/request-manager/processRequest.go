package requestmanager

import (
	"encoding/json"
	"fmt"
	logger "server/logger"

	orderermap "server/orderer-map"

	types "github.com/enriquenc/orderer-map-client-server-go/shared"
)

func ProcessRequests(reqs <-chan types.Request, logger *logger.Logger) {
	// Create an OrderedMap to store the processed requests
	dataStorage := orderermap.NewOrderedMap()

	for req := range reqs {
		// Processing of requests
		switch req.Action {
		case types.AddItem:
			dataStorage.Add(req.Key, req.Value)
			logger.Log(fmt.Sprintf("[add] Added key %s with value %s", req.Key, req.Value))

		case types.RemoveItem:
			exists := dataStorage.Remove(req.Key)
			if exists {
				logger.Log(fmt.Sprintf("[remove] key %s", req.Key))
			} else {
				logger.Log(fmt.Sprintf("[remove] key %s doesn't exist", req.Key))
			}
		case types.GetItem:
			value, exists := dataStorage.Get(req.Key)
			if exists {
				logger.Log(fmt.Sprintf("[get] Got key %s with value %s", req.Key, value))
			} else {
				logger.Log(fmt.Sprintf("[get] Key %s doesn't exist", req.Key))
			}
		case types.GetAll:
			items := dataStorage.GetAll()
			b, _ := json.Marshal(items)
			logger.Log(fmt.Sprintf("[getAll] All values %s", string(b)))
		}
	}
}

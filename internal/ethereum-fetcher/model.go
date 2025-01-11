package ethereumfetcher

import "time"

type SmartContractData struct {
	Source       string `json:"source"`
	Pair         string `json:"pair"`
	Pair0        string `json:"pair0"`
	Pair1        string `json:"pair1"`
	LastUpdated  int64  `json:"last_updated"`
	ExchangeRate string `json:"exchange_rate"`
}

type RedisInfuraData struct {
	Infura SmartContractData `json:"infura"`
}

func NewSmartContractData(source, pair, pair0, pair1, exchangeRate string) SmartContractData {
	return SmartContractData{
		Source:       source,
		Pair:         pair,
		Pair0:        pair0,
		Pair1:        pair1,
		LastUpdated:  time.Now().Unix(),
		ExchangeRate: exchangeRate,
	}
}

func NewRedisInfuraData(infura SmartContractData) RedisInfuraData {
	return RedisInfuraData{
		Infura: infura,
	}
}

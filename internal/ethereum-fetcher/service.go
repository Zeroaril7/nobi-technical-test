package ethereumfetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strings"
	"sync"

	"github.com/Zeroaril7/nobi-technical-test/config"
	"github.com/Zeroaril7/nobi-technical-test/pkg/constant"
	redissdk "github.com/Zeroaril7/nobi-technical-test/pkg/redis-sdk"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var ethereumClient *ethclient.Client
var ethereumClientMutex sync.Mutex

func ConnectToEthereum() (*ethclient.Client, error) {
	ethereumClientMutex.Lock()
	defer ethereumClientMutex.Unlock()

	if ethereumClient == nil {
		api := fmt.Sprintf("%s%v", config.Config("ETHEREUM_RPC_URL"), config.Config("ETHEREUM_RPC_API_KEY"))
		client, err := ethclient.Dial(api)
		if err != nil {
			return nil, err
		}
		ethereumClient = client
	}
	return ethereumClient, nil
}

func FetchExchangeRate() (*big.Float, error) {
	client, err := ConnectToEthereum()
	if err != nil {
		log.Println("Error connecting to Ethereum:", err)
		return nil, err
	}

	contractAddress := common.HexToAddress(config.Config("CONTRACT_ADDRESS"))
	data, err := abi.JSON(strings.NewReader(constant.CONTRACT_ABI))
	if err != nil {
		log.Println("Error parsing contract ABI:", err)
		return nil, err
	}

	callData, err := data.Pack("exchangeRate")
	if err != nil {
		log.Println("Error packing ABI data:", err)
		return nil, err
	}

	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: callData,
	}

	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		log.Println("Error calling contract:", err)
		return nil, err
	}

	exchangeRate := new(big.Int).SetBytes(result)
	decimalShift := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(constant.CONTRACT_DECIMAL), nil))
	finalRate := new(big.Float).Quo(new(big.Float).SetInt(exchangeRate), decimalShift)

	log.Println("Fetched exchange rate:", finalRate.Text('f', 18))
	return finalRate, nil
}

func SaveExchangeRateIntoRedis(rate *big.Float) error {
	key := "price:Crypto:APEETH/ETH"
	smartContract := NewSmartContractData("infura", "Crypto:ALL:APEETH/ETH", "APEETH", "ETH", rate.Text('f', 18))
	redisData := NewRedisInfuraData(smartContract)

	dataJSON, err := json.Marshal(redisData)
	if err != nil {
		log.Println("Error marshaling Redis data:", err)
		return err
	}

	err = redissdk.RedisClient.Set(redissdk.Ctx, key, dataJSON, 0).Err()
	if err != nil {
		log.Println("Error set data in Redis:", err)
		return err
	}

	return nil
}

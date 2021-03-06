package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Stats struct {
	NumBlocks         int `json:"NumBlocks,omitempty"`
	NumTx             int `json:"NumTx,omitempty"`
	NumAuthRequests   int `json:"NumAuthRequests,omitempty"`
	NumUnAuthRequests int `json:"NumUnAuthRequests,omitempty"`
	NumSystemRequests int `json:"NumSystemRequests,omitempty"`
	NumApiConns       int `json:"NumApiConns,omitempty"`
}

var allStats = Stats{NumBlocks: 0, NumTx: 0, NumAuthRequests: 0, NumUnAuthRequests: 0, NumSystemRequests: 0, NumApiConns: 0}
var (
	eventsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "snoopy_processed_events_total",
		Help: "The total number of processed events",
	})
	blocksProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "snoopy_processed_blocks_total",
		Help: "The total number of processed blocks",
	})
	txProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "snoopy_processed_transactions_total",
		Help: "The total number of processed transactions",
	})
	requestsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "snoopy_processed_requests_total",
		Help: "The total number of processed requests",
	})
	apiCallsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "snoopy_processed_apicalls_total",
		Help: "The total number of processed api calls",
	})
)

type Block struct {
	Id                   int    `json:"Id,omitempty"`
	BlockHash            string `json:"BlockHash,omitempty"`
	BlockNumber          uint64 `json:"BlockNumber,omitempty"`
	BlockTime            uint64 `json:"BlockTime,omitempty"`
	BlockNonce           uint64 `json:"BlockNonce,omitempty"`
	BlockNumTransactions int    `json:"BlockNumTransactions,omitempty"`
}

var BlockById map[int]*Block = make(map[int]*Block)
var BlockByNumber map[uint64][]*Block = make(map[uint64][]*Block)
var BlockByHash map[string][]*Block = make(map[string][]*Block)

type Tx struct {
	Id              int    `json:"Id,omitempty"`
	TxBlockId       int    `json:"TxBlockId,omitempty"`
	TxBlockNumber   uint64 `json:"TxBlockNumber,omitempty"`
	TxHash          string `json:"TxHash,omitempty"`
	TxValue         uint64 `json:"TxValue,omitempty"`
	TxGas           uint64 `json:"TxGas,omitempty"`
	TxGasPrice      uint64 `json:"TxGasPrice,omitempty"`
	TxCost          uint64 `json:"TxCost,omitempty"`
	TxNonce         uint64 `json:"TxNonce,omitempty"`
	TxData          string `json:"TxData,omitempty"`
	TxTo            string `json:"TxTo,omitempty"`
	TxReceiptStatus uint64 `json:"TxReceiptStatus,omitempty"`
}

var TxById map[int]*Tx = make(map[int]*Tx)
var TxByTo map[string][]*Tx = make(map[string][]*Tx)
var TxByHash map[string][]*Tx = make(map[string][]*Tx)
var TxByBlockId map[int][]*Tx = make(map[int][]*Tx)
var TxByBlockNumber map[uint64][]*Tx = make(map[uint64][]*Tx)

type Filters struct {
	Id   int    `json:"Id,omitempty"`
	TxTo string `json:"TxTo,omitempty"`
}

var FilterById map[int]*Filters = make(map[int]*Filters)
var FilterByTxTo map[string][]*Filters = make(map[string][]*Filters)

func BlockStore(block Block) {
	BlockById[block.Id] = &block
	BlockByNumber[block.BlockNumber] = append(BlockByNumber[block.BlockNumber], &block)
	BlockByHash[fmt.Sprint(block.BlockHash)] = append(BlockByHash[fmt.Sprint(block.BlockHash)], &block)
}
func TxStore(tx Tx) {
	TxById[tx.Id] = &tx
	TxByTo[tx.TxTo] = append(TxByTo[tx.TxTo], &tx)
	TxByBlockId[tx.TxBlockId] = append(TxByBlockId[tx.TxBlockId], &tx)
	TxByBlockNumber[tx.TxBlockNumber] = append(TxByBlockNumber[tx.TxBlockNumber], &tx)
	TxByHash[fmt.Sprint(tx.TxHash)] = append(TxByHash[fmt.Sprint(tx.TxHash)], &tx)
}
func FilterStore(filter Filters) {
	FilterById[filter.Id] = &filter
	FilterByTxTo[fmt.Sprint(filter.TxTo)] = append(FilterByTxTo[fmt.Sprint(filter.TxTo)], &filter)
}

func check_connect(projectID string, networkName string) bool {
	if projectID == "" {
		log.Fatalf("No projectID found.")
	}
	if networkName == "" {
		log.Fatalf("No networkName found.")
	}

	_, err := ethclient.Dial("wss://" + networkName + ".infura.io/ws/v3/" + projectID)

	if err != nil {
		log.Fatal("Oops! There was a problem", err)
		return false
	} else {
		log.Println("Success! you connected to the " + networkName + " Network")
		allStats.NumApiConns++
		apiCallsProcessed.Inc()
		return true
	}
}

func snoop(wg *sync.WaitGroup, maxBlocks int, ch1 chan bool) bool {
	defer wg.Done()
	projectID := os.Getenv("SNOOPY_PROJECT_ID")
	networkName := os.Getenv("SNOOPY_NETWORK_NAME")

	if check_connect(projectID, networkName) {
		client, err := ethclient.Dial("wss://" + networkName + ".infura.io/ws/v3/" + projectID)
		if err != nil {
			log.Fatal(err)
		}
		headers := make(chan *types.Header)
		sub, err := client.SubscribeNewHead(context.Background(), headers)
		if err != nil {
			log.Fatal(err)
		}
		apiCallsProcessed.Inc()
		allStats.NumApiConns++
		var wgb sync.WaitGroup
		ch2 := make(chan bool)
		var i = 0
		for {
			select {
			case err := <-sub.Err():
				log.Print(err) // Log error and continue
			case header := <-headers:
				i++
				if i >= maxBlocks && maxBlocks > 0 {
					ch1 <- true
					wg.Done()
					return true
				} else {
					wgb.Add(1)
					go snoopProcessEvent(&wgb, i, client, sub, header, maxBlocks, ch2)
					wgb.Wait() // Enable breakout
				}
			}
		}
	}
	return true
}
func snoopProcessEvent(wgb *sync.WaitGroup, i int, client *ethclient.Client, sub ethereum.Subscription, header *types.Header, maxBlocks int, ch2 chan bool) {
	defer wgb.Done()
	// log.Println(header.Hash().Hex()) // 0xbc10defa8dda384c96a17640d84de5578804945d347072e091b4e5f390ddea7f
	eventsProcessed.Inc()
	block, err := client.BlockByHash(context.Background(), header.Hash())
	if err != nil {
		log.Print(err) // Log error and continue
	}
	cBlock := Block{Id: i, BlockHash: block.Hash().Hex(), BlockNumber: block.Number().Uint64(), BlockTime: block.Time(), BlockNonce: block.Nonce(), BlockNumTransactions: len(block.Transactions())}
	BlockStore(cBlock)
	allStats.NumBlocks++
	allStats.NumTx += len(block.Transactions())
	// Combine Prometheus metrics
	blocksProcessed.Inc()
	var txInBlock float64 = float64(len(block.Transactions()))
	txProcessed.Add(txInBlock)
	// Reply with Block Data
	s, err := json.Marshal(BlockById[i])
	if err != nil {
		log.Print(err)
	}
	log.Println("Got: " + string(s))

	blockT, errT := client.BlockByNumber(context.Background(), block.Number())
	if errT != nil {
		log.Print(err) // Log error and continue
	}
	// Id         int    `json:"Id,omitempty"`
	// TxBlockId  uint64 `json:"TxBlockId,omitempty"`
	// TxBlockNumber  uint64 `json:"TxBlockId,omitempty"`
	// TxHash     string `json:"TxHash,omitempty"`
	// TxValue    uint64 `json:"TxValue,omitempty"`
	// TxGas      uint64 `json:"TxGas,omitempty"`
	// TxGasPrice uint64 `json:"TxGasPrice,omitempty"`
	// TxCost 	  uint64 `json:"TxCost,omitempty"`
	// TxNonce    uint64 `json:"TxNonce,omitempty"`
	// TxData     string `json:"TxData,omitempty"`
	// TxTo       string `json:"TxTo,omitempty"`
	// TxReceiptStatus  uint64 `json:"TxTo,omitempty"`
	var ti = 0
	log.Println("Processing #" + block.Number().String())
	for _, tx := range blockT.Transactions() {
		ti++
		// fmt.Println(tx.Hash().Hex())        // 0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2
		// fmt.Println(tx.Value().String())    // 10000000000000000
		// fmt.Println(tx.Gas())               // 105000
		// fmt.Println(tx.GasPrice().Uint64()) // 102000000000
		// fmt.Println(tx.Nonce())             // 110644
		// fmt.Println(tx.Data())              // []
		var TxTo string = "0x0"
		if tx.To() != nil {
			TxTo = tx.To().String()
		}
		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Print(err) // Log error and continue
			continue
		}
		//fmt.Println(receipt.Status) // 1
		var gotTx = 0
		filters := FilterByTxTo
		if len(filters) > 0 {
			// Filters Exists
			filter := FilterByTxTo[TxTo]
			if len(filter) > 0 {
				log.Println("Matched: " + string(TxTo))
				cTx := Tx{Id: ti, TxBlockId: i, TxBlockNumber: block.Number().Uint64(), TxHash: tx.Hash().Hex(), TxValue: tx.Value().Uint64(), TxGas: tx.Gas(), TxGasPrice: tx.GasPrice().Uint64(), TxCost: tx.Cost().Uint64(), TxNonce: tx.Nonce(), TxTo: TxTo, TxReceiptStatus: receipt.Status}
				TxStore(cTx)
				gotTx = 1
			}
		} else {
			// No Filters Store everything
			cTx := Tx{Id: ti, TxBlockId: i, TxBlockNumber: block.Number().Uint64(), TxHash: tx.Hash().Hex(), TxValue: tx.Value().Uint64(), TxGas: tx.Gas(), TxGasPrice: tx.GasPrice().Uint64(), TxCost: tx.Cost().Uint64(), TxNonce: tx.Nonce(), TxTo: TxTo, TxReceiptStatus: receipt.Status}
			TxStore(cTx)
			gotTx = 1
		}
		if gotTx == 1 {
			s, err := json.Marshal(TxById[ti])
			if err != nil {
				log.Print(err)
				continue
			}
			log.Println("Tx: " + string(s))
		}
	}
	log.Println("Done #" + block.Number().String())
}

type App struct {
	Router *mux.Router
	// DB     *sql.DB // If DB backend is required add it here
}

type ProcessSnoopBlockIdRequest struct {
	Id int `json:"id,omitempty"`
}
type ProcessSnoopBlockHashRequest struct {
	Hash string `json:"hash,omitempty"`
}
type ProcessSnoopBlockNumberRequest struct {
	Number uint64 `json:"number,omitempty"`
}
type ProcessSnoopTxIdRequest struct {
	Id int `json:"id,omitempty"`
}
type ProcessSnoopTxNumberRequest struct {
	Number uint64 `json:"number,omitempty"`
}
type ProcessSnoopTxToRequest struct {
	To string `json:"to,omitempty"`
}
type ProcessSnoopFilterToRequest struct {
	To string `json:"to,omitempty"`
}
type ProcessSnoopFilterIdRequest struct {
	Id int `json:"to,omitempty"`
}

// Define our auth struct
type authenticationMiddleware struct {
	tokenUsers map[string]string
}

func (a *App) Initialize() {
	a.Router = mux.NewRouter()
	a.setupRoutes()
}
func (a *App) Run(addr string, wg *sync.WaitGroup) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) setupRoutes() {

	api := a.Router.PathPrefix("/").Subrouter()
	// Protected Routes
	api.HandleFunc("/", a.snoopStatsRequest).Methods("GET")
	api.HandleFunc("/blocks", a.snoopBlocksRequest).Methods("GET")
	api.HandleFunc("/blockid", a.snoopBlockIdRequest).Methods("POST")
	api.HandleFunc("/blockhash", a.snoopBlockHashRequest).Methods("POST")
	api.HandleFunc("/blocknumber", a.snoopBlockNumberRequest).Methods("POST")
	api.HandleFunc("/txs", a.snoopTxRequest).Methods("GET")
	api.HandleFunc("/txid", a.snoopTxIdRequest).Methods("POST")
	api.HandleFunc("/txnumber", a.snoopTxNumberRequest).Methods("POST")
	api.HandleFunc("/filters", a.snoopFiltersRequest).Methods("GET")
	api.HandleFunc("/filterid", a.snoopFilterIdRequest).Methods("POST")
	api.HandleFunc("/filterto", a.snoopFilterToRequest).Methods("POST")
	api.HandleFunc("/filteradd", a.snoopFilterAddToRequest).Methods("POST")
	api.HandleFunc("/filterdelete", a.snoopFilterDeleteIdRequest).Methods("POST")
	// Non Authenticated Routes
	a.Router.HandleFunc("/ping", a.pingRoute).Methods("GET")
	a.Router.HandleFunc("/health", a.healthCheck).Methods("GET")
	// Setup MiddleWare for Auth
	amw := authenticationMiddleware{make(map[string]string)}
	amw.PopulateAllowedTokens()
	api.Use(amw.Middleware)
}

// Loads allowed tokens
func (amw *authenticationMiddleware) PopulateAllowedTokens() {
	// Picks a single predefined token, can be linked to a database of allowed users.
	token := os.Getenv("SNOOPY_API_TOKEN")
	amw.tokenUsers[token] = "infura"
}

// Middleware function, which will be called for each request
func (amw *authenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestsProcessed.Inc()
		token := r.Header.Get("X-Token")
		if user, found := amw.tokenUsers[token]; found {
			// We found the token in our map
			log.Printf("Authenticated user %s\n", user)
			allStats.NumAuthRequests++
			next.ServeHTTP(w, r)
		} else {
			log.Printf("Unauthenticated user\n")
			allStats.NumUnAuthRequests++
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}

func (a *App) pingRoute(w http.ResponseWriter, r *http.Request) {
	log.Printf("ping received\n")
	allStats.NumSystemRequests++
	respondWithJSON(w, http.StatusOK, map[string]string{"ping": "pong"})
}

func (a *App) healthCheck(w http.ResponseWriter, r *http.Request) {
	log.Printf("healthcheck received\n")
	allStats.NumSystemRequests++
	respondWithJSON(w, http.StatusOK, map[string]string{"alive": "true"})
}

func (a *App) snoopStatsRequest(w http.ResponseWriter, r *http.Request) {
	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request"})
		return
	}
	log.Println("Request: /")
	// Reply with Stats
	s, err := json.Marshal(allStats)
	if err != nil {
		log.Print(err)
	}
	log.Println("Sending: " + string(s))
	respondWithJSON(w, http.StatusOK, allStats)
}

func (a *App) snoopBlockHashRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request"})
		return
	}
	var pr ProcessSnoopBlockHashRequest
	err = json.Unmarshal(body, &pr)
	if err != nil {
		println(err.Error())
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request"})
		return
	}

	if pr.Hash == "" {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request. Missing or faulty fields"})
		return
	}
	// Reply with Block Data
	s, err := json.Marshal(BlockByHash[pr.Hash])
	if err != nil {
		log.Print(err)
	}
	log.Println("Sending: " + string(s))
	respondWithJSON(w, http.StatusOK, BlockByHash[pr.Hash])
}

func (a *App) snoopBlockIdRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request"})
		return
	}
	log.Println("Request: " + string(body))
	var pr ProcessSnoopBlockIdRequest
	err = json.Unmarshal(body, &pr)
	if err != nil {
		log.Println(err.Error())
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request"})
		return
	}

	if pr.Id < 1 || pr.Id > allStats.NumBlocks {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request. Missing or faulty fields or out of bounds"})
		return
	}

	// Reply with Block Data
	s, err := json.Marshal(BlockById[pr.Id])
	if err != nil {
		log.Print(err)
	}
	log.Println("Sending: " + string(s))
	respondWithJSON(w, http.StatusOK, BlockById[pr.Id])
}

func (a *App) snoopBlockNumberRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request"})
		return
	}
	log.Println("Request: " + string(body))
	var pr ProcessSnoopBlockNumberRequest
	err = json.Unmarshal(body, &pr)
	if err != nil {
		log.Println(err.Error())
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request"})
		return
	}

	if pr.Number < 1 {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request. Missing or faulty fields or out of bounds"})
		return
	}

	// Reply with Block Data
	s, err := json.Marshal(BlockByNumber[pr.Number])
	if err != nil {
		log.Print(err)
	}
	log.Println("Sending: " + string(s))
	respondWithJSON(w, http.StatusOK, BlockByNumber[pr.Number])
}

func (a *App) snoopBlocksRequest(w http.ResponseWriter, r *http.Request) {
	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request"})
		return
	}
	log.Println("Request: /blocks")
	// Reply with All Blocks
	s, err := json.Marshal(BlockById)
	if err != nil {
		log.Print(err)
	}
	log.Println("Sending: " + string(s))
	respondWithJSON(w, http.StatusOK, BlockById)
}

func (a *App) snoopTxRequest(w http.ResponseWriter, r *http.Request) {
	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request"})
		return
	}
	log.Println("Request: /tx")
	// Reply with All Blocks
	s, err := json.Marshal(TxById)
	if err != nil {
		log.Print(err)
	}
	log.Println("Sending: " + string(s))
	respondWithJSON(w, http.StatusOK, TxById)
}
func (a *App) snoopTxIdRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request"})
		return
	}
	log.Println("Request: " + string(body))
	var pr ProcessSnoopTxIdRequest
	err = json.Unmarshal(body, &pr)
	if err != nil {
		log.Println(err.Error())
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request"})
		return
	}

	if pr.Id < 1 {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request. Missing or faulty fields or out of bounds"})
		return
	}

	// Reply with Block Data
	s, err := json.Marshal(TxById[pr.Id])
	if err != nil {
		log.Print(err)
	}
	log.Println("Sending: " + string(s))
	respondWithJSON(w, http.StatusOK, TxById[pr.Id])
}
func (a *App) snoopTxNumberRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request"})
		return
	}
	log.Println("Request: " + string(body))
	var pr ProcessSnoopTxNumberRequest
	err = json.Unmarshal(body, &pr)
	if err != nil {
		log.Println(err.Error())
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request"})
		return
	}

	if pr.Number < 1 {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request. Missing or faulty fields or out of bounds"})
		return
	}

	// Reply with Block Data
	s, err := json.Marshal(TxByBlockNumber[pr.Number])
	if err != nil {
		log.Print(err)
	}
	log.Println("Sending: " + string(s))
	respondWithJSON(w, http.StatusOK, TxByBlockNumber[pr.Number])
}
func (a *App) snoopFilterIdRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request"})
		return
	}
	log.Println("Request: " + string(body))
	var pr ProcessSnoopFilterIdRequest
	err = json.Unmarshal(body, &pr)
	if err != nil {
		log.Println(err.Error())
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request"})
		return
	}

	if pr.Id > 0 {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request. Missing or faulty fields or out of bounds"})
		return
	}

	// Reply with Block Data
	s, err := json.Marshal(FilterById[pr.Id])
	if err != nil {
		log.Print(err)
	}
	log.Println("Sending: " + string(s))
	respondWithJSON(w, http.StatusOK, FilterById[pr.Id])
}
func (a *App) snoopFilterToRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request"})
		return
	}
	log.Println("Request: " + string(body))
	var pr ProcessSnoopFilterToRequest
	err = json.Unmarshal(body, &pr)
	if err != nil {
		log.Println(err.Error())
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request"})
		return
	}

	if pr.To == "" {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request. Missing or faulty fields or out of bounds"})
		return
	}

	// Reply with Block Data
	s, err := json.Marshal(FilterByTxTo[pr.To])
	if err != nil {
		log.Print(err)
	}
	log.Println("Sending: " + string(s))
	respondWithJSON(w, http.StatusOK, FilterByTxTo[pr.To])
}
func (a *App) snoopFilterAddToRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request"})
		return
	}
	log.Println("Request: " + string(body))
	var pr ProcessSnoopFilterToRequest
	err = json.Unmarshal(body, &pr)
	if err != nil {
		log.Println(err.Error())
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request"})
		return
	}

	if pr.To == "" {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request. Missing or faulty fields or out of bounds"})
		return
	}
	// Add
	AddFilter(pr.To)
	// Reply with Block Data
	s, err := json.Marshal(FilterByTxTo[pr.To])
	if err != nil {
		log.Print(err)
	}
	log.Println("Added Filter: " + string(s))
	respondWithJSON(w, http.StatusOK, FilterByTxTo[pr.To])
}
func (a *App) snoopFilterDeleteIdRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request"})
		return
	}
	log.Println("Request: " + string(body))
	var pr ProcessSnoopFilterIdRequest
	err = json.Unmarshal(body, &pr)
	if err != nil {
		log.Println(err.Error())
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request"})
		return
	}

	if pr.Id > 0 {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request. Missing or faulty fields or out of bounds"})
		return
	}
	// Delete
	DeleteFilter(pr.Id)

	// Reply with Block Data
	log.Println("Deleted Filter " + fmt.Sprint(pr.Id))
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "true"})
}

func (a *App) snoopFiltersRequest(w http.ResponseWriter, r *http.Request) {
	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request"})
		return
	}
	log.Println("Request: /blocks")
	// Reply with All Blocks
	s, err := json.Marshal(FilterById)
	if err != nil {
		log.Print(err)
	}
	log.Println("Sending: " + string(s))
	respondWithJSON(w, http.StatusOK, FilterById)
}
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
func AddFilter(to string) bool {
	cRows := len(FilterById)
	cFIlterRow := Filters{Id: cRows, TxTo: to}
	FilterStore(cFIlterRow)
	return true
}
func DeleteFilter(id int) bool {
	cFilterRow := FilterById[id]
	delete(FilterById, cFilterRow.Id)
	delete(FilterByTxTo, cFilterRow.TxTo)
	return true
}
func prometheusRun(port string, wg *sync.WaitGroup) bool {
	defer wg.Done()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(port, nil)
	return true
}

func main() {
	var wg sync.WaitGroup
	ch1 := make(chan bool)
	a := App{}
	a.Initialize()
	// WaitGroups
	wg.Add(3)
	go a.Run(":9080", &wg)
	go prometheusRun(":2112", &wg)
	go snoop(&wg, 0, ch1)
	wg.Wait()
	close(ch1)
}

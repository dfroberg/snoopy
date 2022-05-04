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

type Transaction struct {
	Id                   int    `json:"Id,omitempty"`
	BlockHash            string `json:"BlockHash,omitempty"`
	BlockNumber          string `json:"BlockNumber,omitempty"`
	BlockTime            uint64 `json:"BlockTime,omitempty"`
	BlockNonce           uint64 `json:"BlockNonce,omitempty"`
	BlockNumTransactions int    `json:"BlockNumTransactions,omitempty"`
}

type TxData struct {
	Id         int    `json:"Id,omitempty"`
	TxHash     string `json:"TxHash,omitempty"`
	TxValue    uint64 `json:"TxValue,omitempty"`
	TxGas      uint64 `json:"TxGas,omitempty"`
	TxGasPrice uint64 `json:"TxGasPrice,omitempty"`
	TxNonce    uint64 `json:"TxNonce,omitempty"`
	TxData     string `json:"TxData,omitempty"`
	TxTo       string `json:"TxTo,omitempty"`
}

var TransactionById map[int]*Transaction
var TransactionByNumber map[string][]*Transaction
var TransactionByHash map[string][]*Transaction

func LocalStore(transaction Transaction) {
	TransactionById[transaction.Id] = &transaction
	TransactionByNumber[transaction.BlockNumber] = append(TransactionByNumber[transaction.BlockNumber], &transaction)
	TransactionByHash[fmt.Sprint(transaction.BlockHash)] = append(TransactionByHash[fmt.Sprint(transaction.BlockHash)], &transaction)
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
		TransactionById = make(map[int]*Transaction)
		TransactionByHash = make(map[string][]*Transaction)
		TransactionByNumber = make(map[string][]*Transaction)

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
		var i = 0
		for {
			select {
			case err := <-sub.Err():
				log.Print(err) // Log error and continue
			case header := <-headers:
				// log.Println(header.Hash().Hex()) // 0xbc10defa8dda384c96a17640d84de5578804945d347072e091b4e5f390ddea7f
				eventsProcessed.Inc()
				block, err := client.BlockByHash(context.Background(), header.Hash())
				if err != nil {
					log.Print(err) // Log error and continue
					continue
				}
				// Increment internal id
				i++
				cTransaction := Transaction{Id: i, BlockHash: block.Hash().Hex(), BlockNumber: fmt.Sprint(block.Number().Uint64()), BlockTime: block.Time(), BlockNonce: block.Nonce(), BlockNumTransactions: len(block.Transactions())}
				LocalStore(cTransaction)
				allStats.NumBlocks++
				allStats.NumTx += len(block.Transactions())
				// Combine Prometheus metrics
				blocksProcessed.Inc()
				var txInBlock float64 = float64(len(block.Transactions()))
				txProcessed.Add(txInBlock)
				// Reply with Transaction Data
				s, err := json.Marshal(TransactionById[i])
				if err != nil {
					log.Print(err)
					continue
				}
				log.Println("Got: " + string(s))

				// blockT, errT := client.BlockByNumber(context.Background(), block.Number())
				// if errT != nil {
				// 	log.Print(err) // Log error and continue
				// 	continue
				// }
				// for _, tx := range blockT.Transactions() {
				// 	fmt.Println(tx.Hash().Hex())        // 0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2
				// 	fmt.Println(tx.Value().String())    // 10000000000000000
				// 	fmt.Println(tx.Gas())               // 105000
				// 	fmt.Println(tx.GasPrice().Uint64()) // 102000000000
				// 	fmt.Println(tx.Nonce())             // 110644
				// 	//fmt.Println(tx.Data())              // []
				// 	if tx.To() != nil {
				// 		fmt.Println(tx.To().Hex()) // 0x55fE59D8Ad77035154dDd0AD0388D09Dd4047A8e
				// 	} else {
				// 		fmt.Println("No to")
				// 	}
				// 	receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
				// 	if err != nil {
				// 		log.Print(err) // Log error and continue
				// 		continue
				// 	}

				// 	fmt.Println(receipt.Status) // 1
				// }
				// Enable breakout
				if i >= maxBlocks && maxBlocks > 0 {
					wg.Done()
					ch1 <- true
				}
			}
		}

	}
	return true
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
	Number string `json:"number,omitempty"`
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
	api.HandleFunc("/", a.processSnoopStatsRequest).Methods("GET")
	api.HandleFunc("/blocks", a.processSnoopBlocksRequest).Methods("GET")
	api.HandleFunc("/blockid", a.processSnoopBlockIdRequest).Methods("POST")
	api.HandleFunc("/blockhash", a.processSnoopBlockHashRequest).Methods("POST")
	api.HandleFunc("/blocknumber", a.processSnoopBlockNumberRequest).Methods("POST")
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

func (a *App) processSnoopStatsRequest(w http.ResponseWriter, r *http.Request) {
	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request"})
		return
	}
	// Reply with Stats
	s, err := json.Marshal(allStats)
	if err != nil {
		log.Print(err)
	}
	log.Println("Sending: " + string(s))
	respondWithJSON(w, http.StatusOK, allStats)
}

func (a *App) processSnoopBlockHashRequest(w http.ResponseWriter, r *http.Request) {
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
	// Reply with Transaction Data
	s, err := json.Marshal(TransactionByHash[pr.Hash])
	if err != nil {
		log.Print(err)
	}
	log.Println("Sending: " + string(s))
	respondWithJSON(w, http.StatusOK, TransactionByHash[pr.Hash])
}

func (a *App) processSnoopBlockIdRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request"})
		return
	}
	println(string(body))
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

	// Reply with Transaction Data
	s, err := json.Marshal(TransactionById[pr.Id])
	if err != nil {
		log.Print(err)
	}
	log.Println("Sending: " + string(s))
	respondWithJSON(w, http.StatusOK, TransactionById[pr.Id])
}

func (a *App) processSnoopBlockNumberRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request"})
		return
	}
	println(string(body))
	var pr ProcessSnoopBlockNumberRequest
	err = json.Unmarshal(body, &pr)
	if err != nil {
		log.Println(err.Error())
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request"})
		return
	}

	if pr.Number == "" {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request. Missing or faulty fields or out of bounds"})
		return
	}

	// Reply with Transaction Data
	s, err := json.Marshal(TransactionByNumber[pr.Number])
	if err != nil {
		log.Print(err)
	}
	log.Println("Sending: " + string(s))
	respondWithJSON(w, http.StatusOK, TransactionByNumber[pr.Number])
}

func (a *App) processSnoopBlocksRequest(w http.ResponseWriter, r *http.Request) {
	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"result": "false", "error": "Invalid Request"})
		return
	}

	// Reply with All Transaction Blocks
	s, err := json.Marshal(TransactionById)
	if err != nil {
		log.Print(err)
	}
	log.Println("Sending: " + string(s))
	respondWithJSON(w, http.StatusOK, TransactionById)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
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
}

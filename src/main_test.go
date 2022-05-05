package main

import (
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const mockTransactions = `{
	"1": {
	  "Id": 1,
	  "TxBlockId": 22,
	  "TxBlockNumber": 14717116,
	  "TxHash": "0xda6a3cce91baf1a4d8648fb659c6630078fae6a72ea12caff6aff975d92617c7",
	  "TxValue": 1.5e+19,
	  "TxGas": 393786,
	  "TxGasPrice": 35792179968,
	  "TxCost": 15014094459380880000,
	  "TxNonce": 1018,
	  "TxTo": "0x75A6787C7EE60424358B449B539A8b774c9B4862",
	  "TxReceiptStatus": 1
	},
	"10": {
	  "Id": 10,
	  "TxBlockId": 22,
	  "TxBlockNumber": 14717116,
	  "TxHash": "0x8f2632a046d96c240b869b596bf750c427fc4d4808878e13ee5b33f7ca3dfbd3",
	  "TxValue": 7e+17,
	  "TxGas": 27938,
	  "TxGasPrice": 29770083149,
	  "TxCost": 700831716583016700,
	  "TxNonce": 68,
	  "TxTo": "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
	  "TxReceiptStatus": 1
	},
	"100": {
	  "Id": 100,
	  "TxBlockId": 22,
	  "TxBlockNumber": 14717116,
	  "TxHash": "0xdf9d697b90c6d196d34359b0bdcccdaad05b57a460bbd66bdec88f9dd7227fb3",
	  "TxValue": 15375963677623800,
	  "TxGas": 21000,
	  "TxGasPrice": 56000000000,
	  "TxCost": 16551963677623800,
	  "TxNonce": 1,
	  "TxTo": "0xA090e606E30bD747d4E6245a1517EbE430F0057e",
	  "TxReceiptStatus": 1
	},
	"101": {
	  "Id": 101,
	  "TxBlockId": 22,
	  "TxBlockNumber": 14717116,
	  "TxHash": "0x1b54d0f0f28dd379ae7fa276011721c07c4f4470937db586a60b169059afb91f",
	  "TxGas": 99367,
	  "TxGasPrice": 36717479710,
	  "TxCost": 3648505806343570,
	  "TxNonce": 32034,
	  "TxTo": "0x31eFc4AeAA7c39e54A33FDc3C46ee2Bd70ae0A09",
	  "TxReceiptStatus": 1
	},
	"102": {
	  "Id": 102,
	  "TxBlockId": 22,
	  "TxBlockNumber": 14717116,
	  "TxHash": "0x895621d07cae4e66d98927fb972f69b41224ca637be6c9b3bfcc934d7e8c6ae0",
	  "TxValue": 15376311699504468,
	  "TxGas": 21000,
	  "TxGasPrice": 56000000000,
	  "TxCost": 16552311699504468,
	  "TxNonce": 1,
	  "TxTo": "0xA090e606E30bD747d4E6245a1517EbE430F0057e",
	  "TxReceiptStatus": 1
	},
	"103": {
	  "Id": 103,
	  "TxBlockId": 22,
	  "TxBlockNumber": 14717116,
	  "TxHash": "0x6e1052bd5745a300756cdfaa7fa63c75c672be57bd97898f40bf9c03e276053b",
	  "TxGas": 60760,
	  "TxGasPrice": 38530000000,
	  "TxCost": 2341082800000000,
	  "TxNonce": 65,
	  "TxTo": "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
	  "TxReceiptStatus": 1
	},
	"104": {
	  "Id": 104,
	  "TxBlockId": 22,
	  "TxBlockNumber": 14717116,
	  "TxHash": "0xd5daa38167d77104e3ddd3d00d09726a59f361c34f68a13454d895f51ba0df8d",
	  "TxValue": 8477520000000000,
	  "TxGas": 21000,
	  "TxGasPrice": 56000000000,
	  "TxCost": 9653520000000000,
	  "TxNonce": 7015864,
	  "TxTo": "0xF3929C9eED603E4617CF02908565f06963492B5d",
	  "TxReceiptStatus": 1
	},
	"105": {
	  "Id": 105,
	  "TxBlockId": 22,
	  "TxBlockNumber": 14717116,
	  "TxHash": "0xa3ba2f4e45f63d4b0858d823227da50d923713844f56f4e89677dfd68d38b5ef",
	  "TxGas": 81986,
	  "TxGasPrice": 36717479710,
	  "TxCost": 3010319291504060,
	  "TxNonce": 79,
	  "TxTo": "0x7f268357A8c2552623316e2562D90e642bB538E5",
	  "TxReceiptStatus": 1
	}`

func loadMockTransactions() {
	raw := []byte(mockTransactions)
	rec := &Receipt{
		Receipt:     &types.Receipt{},
		BlockNumber: 0,
	}
	err := rec.UnmarshalJSON(mockTransactions)
	require.NoError(t, err)
	assert.Equal(t, int64(255), rec.BlockNumber)
}
func loadMockBlocks() {

	var cBlockRow Block
	cBlockRow = Block{Id: 1, BlockHash: "0x547bd8bd5f9c8eee5d2be941f275ee95672632159b8981df7917335963642fbe", BlockNumber: 12232752, BlockTime: 1651499015, BlockNonce: 4627854504322470268, BlockNumTransactions: 8}
	BlockStore(cBlockRow)
	cBlockRow = Block{Id: 3, BlockHash: "0xe441ec0412436c460e4430881ba24a6b1fc8cdb35e3d462a77bfd616021b79b1", BlockNumber: 12232754, BlockTime: 1651499051, BlockNonce: 8148927535907424638, BlockNumTransactions: 7}
	BlockStore(cBlockRow)
	cBlockRow = Block{Id: 4, BlockHash: "0x74c13672c717b3651f058d3f8a45cd0abd58c4bc9f4c33745f57ce37541062de", BlockNumber: 12232755, BlockTime: 1651499052, BlockNonce: 92853587781942119, BlockNumTransactions: 24}
	BlockStore(cBlockRow)
	cBlockRow = Block{Id: 5, BlockHash: "0x370b543677498bb45d0c64f3dd8734a2ac8e3400e6ec1da92eda62637fb456fe", BlockNumber: 12232756, BlockTime: 1651499109, BlockNonce: 1199686451943716732, BlockNumTransactions: 43}
	BlockStore(cBlockRow)
	cBlockRow = Block{Id: 6, BlockHash: "0xd1edce80576d1c8d131c688d85b31dcbeabbbe49d40c730588a857ef1b8d3656", BlockNumber: 12232757, BlockTime: 1651499189, BlockNonce: 5538472040065513350, BlockNumTransactions: 11}
	BlockStore(cBlockRow)
	cBlockRow = Block{Id: 7, BlockHash: "0x7f5ffdb0c52c04ecdd37ed695ce9c8d324a00e9f89a57335aeb357f121cdc722", BlockNumber: 12232758, BlockTime: 1651499227, BlockNonce: 4300166683628318943, BlockNumTransactions: 7}
	BlockStore(cBlockRow)
	cBlockRow = Block{Id: 8, BlockHash: "0x8d802579baa0426935bf38308635960891058c711bd3745adaf2c38f6f23404b", BlockNumber: 12232759, BlockTime: 1651499234, BlockNonce: 3420106972170463577, BlockNumTransactions: 21}
	BlockStore(cBlockRow)
	cBlockRow = Block{Id: 9, BlockHash: "0xe6945ef9ad3c0d2b85fe6941a04bf89449a9e2f3d0684ebd52819eaef8a9f291", BlockNumber: 12232760, BlockTime: 1651499253, BlockNonce: 131136441185041252, BlockNumTransactions: 4}
	BlockStore(cBlockRow)
	cBlockRow = Block{Id: 10, BlockHash: "0xf156931366044f66f8dbfa42233c497babae15984e1e05eb330f120383562ab5", BlockNumber: 12232761, BlockTime: 1651499270, BlockNonce: 8425941817338320609, BlockNumTransactions: 23}
	BlockStore(cBlockRow)
	cBlockRow = Block{Id: 11, BlockHash: "0xf0d2dd4da8e3a7d6181ad78e0740989af4dba18534f102d4c1e01e3084cf0019", BlockNumber: 12232762, BlockTime: 1651499274, BlockNonce: 1199686451756006312, BlockNumTransactions: 4}
	BlockStore(cBlockRow)
	cBlockRow = Block{Id: 12, BlockHash: "0x5c382af5ac37acc7b6099cac522b7a40b084d21b3f81240b2e84982e1a483403", BlockNumber: 12232763, BlockTime: 1651499285, BlockNonce: 9289304408881137609, BlockNumTransactions: 4}
	BlockStore(cBlockRow)
	cBlockRow = Block{Id: 13, BlockHash: "0x17631974cc5e0458e39750ef013bba2ff0f3dca5d4a31c62698faea6befa4423", BlockNumber: 12232764, BlockTime: 1651499288, BlockNonce: 3977542852021753102, BlockNumTransactions: 38}
	BlockStore(cBlockRow)
	cBlockRow = Block{Id: 14, BlockHash: "0x0f683615911992f62e18b5aa3629fe35d58bec6610ff7c813fd8df19c44ef13c", BlockNumber: 12232765, BlockTime: 1651499339, BlockNonce: 5435131497541826742, BlockNumTransactions: 25}
	BlockStore(cBlockRow)
	cBlockRow = Block{Id: 15, BlockHash: "0x8403c1946e7fd1922249bdee11f61a4f516a32eea5595d03b2f1270060be5268", BlockNumber: 12232766, BlockTime: 1651499384, BlockNonce: 8723929707337473946, BlockNumTransactions: 6}
	BlockStore(cBlockRow)
	cBlockRow = Block{Id: 16, BlockHash: "0xa43dc26b3dfe7c57717bacfa8dd6536a2c70e55ce669e1be05431f203f9f6be3", BlockNumber: 12232767, BlockTime: 1651499411, BlockNonce: 9289304408917954889, BlockNumTransactions: 34}
	BlockStore(cBlockRow)
	cBlockRow = Block{Id: 17, BlockHash: "0xead7857d44d244377b651799c23d900d91a1eb66d27ee4e4e33c9e590a24c740", BlockNumber: 12232768, BlockTime: 1651499499, BlockNonce: 9289304408840996491, BlockNumTransactions: 13}
	BlockStore(cBlockRow)
	cBlockRow = Block{Id: 18, BlockHash: "0x712a95c2c07f3de120edef98b44def573cd186486246e1e16643df5b8ab9042f", BlockNumber: 12232769, BlockTime: 1651499540, BlockNonce: 7562348827223693395, BlockNumTransactions: 10}
	BlockStore(cBlockRow)
	cBlockRow = Block{Id: 19, BlockHash: "0xa8341a8747775e23c70a96308e756538fd2a6b5a5603882b5f6aaeacc03eeb27", BlockNumber: 12232770, BlockTime: 1651499552, BlockNonce: 3948699871155854888, BlockNumTransactions: 7}
	BlockStore(cBlockRow)
	cBlockRow = Block{Id: 20, BlockHash: "0x45fb418084fd5862ec4226bec5c5f850b0d204e8aec6e05f42d8f0c7b8c11b24", BlockNumber: 12232771, BlockTime: 1651499568, BlockNonce: 6354333488874994072, BlockNumTransactions: 2}
	BlockStore(cBlockRow)
	cBlockRow = Block{Id: 21, BlockHash: "0xa7c501947a673f7e2de41f4f5fd0993c78c516c5bf230d706e0f2e7a38c3b11a", BlockNumber: 12232772, BlockTime: 1651499570, BlockNonce: 4300166683577701297, BlockNumTransactions: 10}
	BlockStore(cBlockRow)
	cBlockRow = Block{Id: 22, BlockHash: "0x9069fee4bb0d739fbcbc709ac874c65ee3dd252b9eceb7c4417f4689eed1491a", BlockNumber: 12232773, BlockTime: 1651499596, BlockNonce: 6971670793103091897, BlockNumTransactions: 12}
	BlockStore(cBlockRow)
	cBlockRow = Block{Id: 23, BlockHash: "0x5a4bda11cd76a1eebd78d666ec7b1676f28d790b24b9249d8b724e51cf1b4c7c", BlockNumber: 12232774, BlockTime: 1651499607, BlockNonce: 5006083717830328476, BlockNumTransactions: 7}
	BlockStore(cBlockRow)
	cBlockRow = Block{Id: 24, BlockHash: "0xbf4b987cfa926be430b1fe7669f6364910ab379289b8985ba63c8d8c27d656be", BlockNumber: 12232775, BlockTime: 1651499639, BlockNonce: 8692680990958600297, BlockNumTransactions: 1}
	BlockStore(cBlockRow)
	cBlockRow = Block{Id: 25, BlockHash: "0xf1e4b90e7a708a3e4e48762c7fa7dc2a9f6e01dcf70bd8ab36ee545a3c7c0b55", BlockNumber: 12232776, BlockTime: 1651499650, BlockNonce: 5536629777038220738, BlockNumTransactions: 8}
	BlockStore(cBlockRow)
	cBlockRow = Block{Id: 26, BlockHash: "0x4fa272f0a9ae5edf8b5d35a2075dc22ec25f6889f495d236d535b295146f3542", BlockNumber: 12232777, BlockTime: 1651499682, BlockNonce: 8166034474651813189, BlockNumTransactions: 172}
	BlockStore(cBlockRow)
	cBlockRow = Block{Id: 27, BlockHash: "0xd1ee549faee24058432f750a6f3aa5e5a96789b4bed29914da95e3b8c98b9ee0", BlockNumber: 12232778, BlockTime: 1651499823, BlockNonce: 1199686451859900871, BlockNumTransactions: 17}
	BlockStore(cBlockRow)

}
func TestLocalStorage(t *testing.T) {
	loadMockBlocks()
	// Check loaded Blocks
	assert.Equal(t, int(1), BlockById[1].Id)
	assert.Equal(t, int(22), BlockById[22].Id)
}

func TestBlockIdAccess(t *testing.T) {
	loadMockBlocks()
	assert.Equal(t, int(4), BlockById[4].Id)
}

func TestBlockHashAccess(t *testing.T) {
	loadMockBlocks()
	cBlock := BlockByHash["0xd1ee549faee24058432f750a6f3aa5e5a96789b4bed29914da95e3b8c98b9ee0"]
	assert.Equal(t, string("0xd1ee549faee24058432f750a6f3aa5e5a96789b4bed29914da95e3b8c98b9ee0"), cBlock.BlockHash)
}

func TestBlockNumberAccess(t *testing.T) {
	loadMockBlocks()
	cBlock := BlockByNumber[12232778]
	assert.Equal(t, uint64(12232778), cBlock.BlockNumber)
}

func TestTxIdAccess(t *testing.T) {
	loadMockBlocks()
	assert.Equal(t, int(4), BlockById[4].Id)
}
func TestNetworkAccess(t *testing.T) {
	projectID := os.Getenv("SNOOPY_PROJECT_ID")
	networkName := os.Getenv("SNOOPY_NETWORK_NAME")
	assert.Equal(t, bool(true), check_connect(projectID, networkName))
}
func TestAddFilter(t *testing.T) {
	assert.Equal(t, bool(true), AddFilter("0xE592427A0AEce92De3Edee1F18E0157C05861564"))
}
func TestDeleteFilter(t *testing.T) {
	assert.Equal(t, bool(true), AddFilter("0xE592427A0AEce92De3Edee1F18E0157C05861564"))
	assert.Equal(t, bool(true), DeleteFilter(1))
}

// func TestSnoop(t *testing.T) {
// 	var wg sync.WaitGroup
// 	wg.Add(2)
// 	ch1 := make(chan bool)
// 	// Run Snoop, collect 1 block and return
// 	go snoop(&wg, 1, ch1)
// 	var r bool = <-ch1
// 	assert.Equal(t, bool(true), r)
// 	close(ch1)
// }

// func TestSnoopWithFilter(t *testing.T) {
// 	var wg sync.WaitGroup
// 	wg.Add(2)
// 	ch1 := make(chan bool)
// 	// Run Snoop, collect 1 block and return
// 	assert.Equal(t, bool(true), AddFilter("0x00000000219ab540356cbb839cbe05303d7705fa"))
// 	go snoop(&wg, 10, ch1)
// 	var r bool = <-ch1
// 	assert.Equal(t, bool(true), r)
// 	close(ch1)
// }
func TestMain(t *testing.T) {
	a := App{}
	a.Initialize()
}

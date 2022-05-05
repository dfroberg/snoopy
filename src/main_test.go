package main

import (
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalStorage(t *testing.T) {

	var cTransactionRow Transaction
	cTransactionRow = Transaction{Id: 1, BlockHash: "0x547bd8bd5f9c8eee5d2be941f275ee95672632159b8981df7917335963642fbe", BlockNumber: "12232752", BlockTime: 1651499015, BlockNonce: 4627854504322470268, BlockNumTransactions: 8}
	LocalStore(cTransactionRow)
	assert.Equal(t, int(1), TransactionById[1].Id)
	cTransactionRow = Transaction{Id: 3, BlockHash: "0xe441ec0412436c460e4430881ba24a6b1fc8cdb35e3d462a77bfd616021b79b1", BlockNumber: "12232754", BlockTime: 1651499051, BlockNonce: 8148927535907424638, BlockNumTransactions: 7}
	LocalStore(cTransactionRow)
	cTransactionRow = Transaction{Id: 4, BlockHash: "0x74c13672c717b3651f058d3f8a45cd0abd58c4bc9f4c33745f57ce37541062de", BlockNumber: "12232755", BlockTime: 1651499052, BlockNonce: 92853587781942119, BlockNumTransactions: 24}
	LocalStore(cTransactionRow)
	cTransactionRow = Transaction{Id: 5, BlockHash: "0x370b543677498bb45d0c64f3dd8734a2ac8e3400e6ec1da92eda62637fb456fe", BlockNumber: "12232756", BlockTime: 1651499109, BlockNonce: 1199686451943716732, BlockNumTransactions: 43}
	LocalStore(cTransactionRow)
	cTransactionRow = Transaction{Id: 6, BlockHash: "0xd1edce80576d1c8d131c688d85b31dcbeabbbe49d40c730588a857ef1b8d3656", BlockNumber: "12232757", BlockTime: 1651499189, BlockNonce: 5538472040065513350, BlockNumTransactions: 11}
	LocalStore(cTransactionRow)
	cTransactionRow = Transaction{Id: 7, BlockHash: "0x7f5ffdb0c52c04ecdd37ed695ce9c8d324a00e9f89a57335aeb357f121cdc722", BlockNumber: "12232758", BlockTime: 1651499227, BlockNonce: 4300166683628318943, BlockNumTransactions: 7}
	LocalStore(cTransactionRow)
	cTransactionRow = Transaction{Id: 8, BlockHash: "0x8d802579baa0426935bf38308635960891058c711bd3745adaf2c38f6f23404b", BlockNumber: "12232759", BlockTime: 1651499234, BlockNonce: 3420106972170463577, BlockNumTransactions: 21}
	LocalStore(cTransactionRow)
	cTransactionRow = Transaction{Id: 9, BlockHash: "0xe6945ef9ad3c0d2b85fe6941a04bf89449a9e2f3d0684ebd52819eaef8a9f291", BlockNumber: "12232760", BlockTime: 1651499253, BlockNonce: 131136441185041252, BlockNumTransactions: 4}
	LocalStore(cTransactionRow)
	cTransactionRow = Transaction{Id: 10, BlockHash: "0xf156931366044f66f8dbfa42233c497babae15984e1e05eb330f120383562ab5", BlockNumber: "12232761", BlockTime: 1651499270, BlockNonce: 8425941817338320609, BlockNumTransactions: 23}
	LocalStore(cTransactionRow)
	cTransactionRow = Transaction{Id: 11, BlockHash: "0xf0d2dd4da8e3a7d6181ad78e0740989af4dba18534f102d4c1e01e3084cf0019", BlockNumber: "12232762", BlockTime: 1651499274, BlockNonce: 1199686451756006312, BlockNumTransactions: 4}
	LocalStore(cTransactionRow)
	cTransactionRow = Transaction{Id: 12, BlockHash: "0x5c382af5ac37acc7b6099cac522b7a40b084d21b3f81240b2e84982e1a483403", BlockNumber: "12232763", BlockTime: 1651499285, BlockNonce: 9289304408881137609, BlockNumTransactions: 4}
	LocalStore(cTransactionRow)
	cTransactionRow = Transaction{Id: 13, BlockHash: "0x17631974cc5e0458e39750ef013bba2ff0f3dca5d4a31c62698faea6befa4423", BlockNumber: "12232764", BlockTime: 1651499288, BlockNonce: 3977542852021753102, BlockNumTransactions: 38}
	LocalStore(cTransactionRow)
	cTransactionRow = Transaction{Id: 14, BlockHash: "0x0f683615911992f62e18b5aa3629fe35d58bec6610ff7c813fd8df19c44ef13c", BlockNumber: "12232765", BlockTime: 1651499339, BlockNonce: 5435131497541826742, BlockNumTransactions: 25}
	LocalStore(cTransactionRow)
	cTransactionRow = Transaction{Id: 15, BlockHash: "0x8403c1946e7fd1922249bdee11f61a4f516a32eea5595d03b2f1270060be5268", BlockNumber: "12232766", BlockTime: 1651499384, BlockNonce: 8723929707337473946, BlockNumTransactions: 6}
	LocalStore(cTransactionRow)
	cTransactionRow = Transaction{Id: 16, BlockHash: "0xa43dc26b3dfe7c57717bacfa8dd6536a2c70e55ce669e1be05431f203f9f6be3", BlockNumber: "12232767", BlockTime: 1651499411, BlockNonce: 9289304408917954889, BlockNumTransactions: 34}
	LocalStore(cTransactionRow)
	cTransactionRow = Transaction{Id: 17, BlockHash: "0xead7857d44d244377b651799c23d900d91a1eb66d27ee4e4e33c9e590a24c740", BlockNumber: "12232768", BlockTime: 1651499499, BlockNonce: 9289304408840996491, BlockNumTransactions: 13}
	LocalStore(cTransactionRow)
	cTransactionRow = Transaction{Id: 18, BlockHash: "0x712a95c2c07f3de120edef98b44def573cd186486246e1e16643df5b8ab9042f", BlockNumber: "12232769", BlockTime: 1651499540, BlockNonce: 7562348827223693395, BlockNumTransactions: 10}
	LocalStore(cTransactionRow)
	cTransactionRow = Transaction{Id: 19, BlockHash: "0xa8341a8747775e23c70a96308e756538fd2a6b5a5603882b5f6aaeacc03eeb27", BlockNumber: "12232770", BlockTime: 1651499552, BlockNonce: 3948699871155854888, BlockNumTransactions: 7}
	LocalStore(cTransactionRow)
	cTransactionRow = Transaction{Id: 20, BlockHash: "0x45fb418084fd5862ec4226bec5c5f850b0d204e8aec6e05f42d8f0c7b8c11b24", BlockNumber: "12232771", BlockTime: 1651499568, BlockNonce: 6354333488874994072, BlockNumTransactions: 2}
	LocalStore(cTransactionRow)
	cTransactionRow = Transaction{Id: 21, BlockHash: "0xa7c501947a673f7e2de41f4f5fd0993c78c516c5bf230d706e0f2e7a38c3b11a", BlockNumber: "12232772", BlockTime: 1651499570, BlockNonce: 4300166683577701297, BlockNumTransactions: 10}
	LocalStore(cTransactionRow)
	cTransactionRow = Transaction{Id: 22, BlockHash: "0x9069fee4bb0d739fbcbc709ac874c65ee3dd252b9eceb7c4417f4689eed1491a", BlockNumber: "12232773", BlockTime: 1651499596, BlockNonce: 6971670793103091897, BlockNumTransactions: 12}
	LocalStore(cTransactionRow)
	cTransactionRow = Transaction{Id: 23, BlockHash: "0x5a4bda11cd76a1eebd78d666ec7b1676f28d790b24b9249d8b724e51cf1b4c7c", BlockNumber: "12232774", BlockTime: 1651499607, BlockNonce: 5006083717830328476, BlockNumTransactions: 7}
	LocalStore(cTransactionRow)
	cTransactionRow = Transaction{Id: 24, BlockHash: "0xbf4b987cfa926be430b1fe7669f6364910ab379289b8985ba63c8d8c27d656be", BlockNumber: "12232775", BlockTime: 1651499639, BlockNonce: 8692680990958600297, BlockNumTransactions: 1}
	LocalStore(cTransactionRow)
	cTransactionRow = Transaction{Id: 25, BlockHash: "0xf1e4b90e7a708a3e4e48762c7fa7dc2a9f6e01dcf70bd8ab36ee545a3c7c0b55", BlockNumber: "12232776", BlockTime: 1651499650, BlockNonce: 5536629777038220738, BlockNumTransactions: 8}
	LocalStore(cTransactionRow)
	cTransactionRow = Transaction{Id: 26, BlockHash: "0x4fa272f0a9ae5edf8b5d35a2075dc22ec25f6889f495d236d535b295146f3542", BlockNumber: "12232777", BlockTime: 1651499682, BlockNonce: 8166034474651813189, BlockNumTransactions: 172}
	LocalStore(cTransactionRow)
	cTransactionRow = Transaction{Id: 27, BlockHash: "0xd1ee549faee24058432f750a6f3aa5e5a96789b4bed29914da95e3b8c98b9ee0", BlockNumber: "12232778", BlockTime: 1651499823, BlockNonce: 1199686451859900871, BlockNumTransactions: 17}
	LocalStore(cTransactionRow)

	// Check loaded transactions
	assert.Equal(t, int(22), TransactionById[22].Id)
}

func TestIdAccess(t *testing.T) {
	cTransactionRow := Transaction{Id: 4, BlockHash: "0xd1ee549faee24058432f750a6f3aa5e5a96789b4bed29914da95e3b8c98b9ee0", BlockNumber: "12232778", BlockTime: 1651499823, BlockNonce: 1199686451859900871, BlockNumTransactions: 17}
	LocalStore(cTransactionRow)
	assert.Equal(t, int(4), cTransactionRow.Id)
}

func TestHashAccess(t *testing.T) {
	cTransactionRow := Transaction{Id: 4, BlockHash: "0xd1ee549faee24058432f750a6f3aa5e5a96789b4bed29914da95e3b8c98b9ee0", BlockNumber: "12232778", BlockTime: 1651499823, BlockNonce: 1199686451859900871, BlockNumTransactions: 17}
	LocalStore(cTransactionRow)
	assert.Equal(t, string("0xd1ee549faee24058432f750a6f3aa5e5a96789b4bed29914da95e3b8c98b9ee0"), cTransactionRow.BlockHash)
}

func TestBlockNumberAccess(t *testing.T) {
	cTransactionRow := Transaction{Id: 4, BlockHash: "0xd1ee549faee24058432f750a6f3aa5e5a96789b4bed29914da95e3b8c98b9ee0", BlockNumber: "12232778", BlockTime: 1651499823, BlockNonce: 1199686451859900871, BlockNumTransactions: 17}
	LocalStore(cTransactionRow)
	assert.Equal(t, string("12232778"), cTransactionRow.BlockNumber)
}

func TestNetworkAccess(t *testing.T) {
	projectID := os.Getenv("SNOOPY_PROJECT_ID")
	networkName := os.Getenv("SNOOPY_NETWORK_NAME")
	assert.Equal(t, bool(true), check_connect(projectID, networkName))
}

func TestSnoop(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	ch1 := make(chan bool)
	// Run Snoop, collect 1 block and return
	go snoop(&wg, 1, ch1)
	var r bool = <-ch1
	assert.Equal(t, bool(true), r)
	close(ch1)
}

func TestMain(t *testing.T) {
	a := App{}
	a.Initialize()
}

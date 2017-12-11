package transfer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_binance_Signature(t *testing.T) {
	b := MakeBinance(nil, "NhqPtmdSJYdKjVHjA7PZj4Mge3R5YNiP1e3UZjInClVN65XAbvqqM6A7H5fATj0j", "vmPUZE6mv9SD5VNHk4HlWFsOr6aKE2zvsw0MuIgwCIPy6utIco14y7Ju91duEh8A")
	assert.Equal(t,
		"157fb937ec848b5f802daa4d9f62bea08becbf4f311203bda2bd34cd9853e320",
		b.(*binance).Signature("asset=ETH&address=0x6915f16f8791d0a1cc2bf47c13a6b2a92000504b&amount=1&recvWindow=5000&name=test&timestamp=1510903211000"),
	)
}

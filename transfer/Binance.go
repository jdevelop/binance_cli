package transfer

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const Api = "https://api.binance.com/"

type binance struct {
	client  *http.Client
	pubKey  string
	privKey string
}

type Asset = string

type HistoryRecord struct {
	ID        string  `json:"id,omitempty"`
	Amount    float64 `json:"amount,omitempty"`
	Address   string  `json:"address,omitempty"`
	Asset     Asset   `json:"asset,omitempty"`
	Txid      string  `json:"txId,omitempty"`
	ApplyTime int64   `json:"applyTime,omitempty"`
	Status    int     `json:"status,omitempty"`
}

type HasHistory interface {
	History() []HistoryRecord
}

type WithdrawHistoryWrapper struct {
	WithdrawList []HistoryRecord `json:"withdrawList,omitempty"`
}

func (wh *WithdrawHistoryWrapper) History() []HistoryRecord {
	return wh.WithdrawList
}

type DepositHistoryWrapper struct {
	DepositList []HistoryRecord `json:"depositList,omitempty"`
}

func (wh *DepositHistoryWrapper) History() []HistoryRecord {
	return wh.DepositList
}

type UserStatus struct {
	Msg     string   `json:"msg,omitempty"`
	Success bool     `json:"success,omitempty"`
	Objs    []string `json:"objs,omitempty"`
}

type BinanceAccess interface {
	DepositHistory() ([]HistoryRecord, error)
	WithdrawHistory() ([]HistoryRecord, error)
	Status() (UserStatus, error)
	Withdraw(asset Asset, address string, amt float64) error
}

func (ba *binance) Signature(req string) (hash string) {
	mac := hmac.New(sha256.New, []byte(ba.privKey))
	mac.Write([]byte(req))
	hash = hex.EncodeToString(mac.Sum(nil))
	return
}

func (ba *binance) SignRequest(req *http.Request) {
	req.Header.Set("X-MBX-APIKEY", ba.pubKey)
}

func timestampStr() string {
	return fmt.Sprintf("timestamp=%d", time.Now().Unix()*1000)
}

func (ba *binance) Status() (ua UserStatus, err error) {
	ct := timestampStr()
	get, _ := http.NewRequest("GET", Api+"/wapi/v3/accountStatus.html?"+ct+"&signature="+ba.Signature(ct), nil)
	ba.SignRequest(get)
	resp, err := ba.client.Do(get)
	if err != nil {
		return
	}
	str, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(str))

	if err != nil {
		return
	}

	err = json.Unmarshal(str, &ua)
	return
}

func (ba *binance) DepositHistory() (res []HistoryRecord, err error) {
	res, err = ba.history("/wapi/v3/depositHistory.html", func() HasHistory {
		return &DepositHistoryWrapper{}
	})
	return
}

func (ba *binance) WithdrawHistory() (res []HistoryRecord, err error) {
	res, err = ba.history("/wapi/v3/withdrawHistory.html", func() HasHistory {
		return &WithdrawHistoryWrapper{}
	})
	return
}

func (ba *binance) history(path string, p func() HasHistory) (res []HistoryRecord, err error) {
	ct := timestampStr()
	get, _ := http.NewRequest("GET", Api+path+"?"+ct+"&signature="+ba.Signature(ct), nil)
	ba.SignRequest(get)
	resp, err := ba.client.Do(get)
	if err != nil {
		return
	}
	str, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	// fmt.Println(string(str))

	hw := p()

	err = json.Unmarshal(str, hw)
	if err != nil {
		return
	}

	res = hw.History()

	return
}

type WithdrawStatus struct {
	Complete bool   `json:"success,omitempty"`
	Message  string `json:"msg,omitempty"`
}

func (ba *binance) Withdraw(asset Asset, address string, amt float64) (err error) {
	content := "asset=" + (strings.ToUpper(asset)) + "&address=" + address + "&amount=" + fmt.Sprintf("%.5f", amt) + "&recvWindow=4000&" + timestampStr()
	sig := ba.Signature(content)
	post, _ := http.NewRequest("POST", Api+"/wapi/v3/withdraw.html", strings.NewReader(content+"&signature="+sig))
	ba.SignRequest(post)
	resp, err := ba.client.Do(post)
	if err != nil {
		return
	}
	str, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var ws WithdrawStatus

	err = json.Unmarshal(str, &ws)
	if err != nil {
		return
	}

	if !ws.Complete {
		err = errors.New(ws.Message)
	}

	return
}

func MakeBinance(client *http.Client, priv string, public string) BinanceAccess {
	return &binance{
		client:  client,
		privKey: priv,
		pubKey:  public,
	}
}

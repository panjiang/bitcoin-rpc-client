package rpcclient

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type jsonRPCRequest struct {
	Version string        `json:"jsonrpc"` // default '1.0'
	ID      string        `json:"id"`      // default 'jsonrpc'
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

// JSONRPCResult repc result
type JSONRPCResult struct {
	Data  interface{}   `json:"result"`
	Error *JSONRPCError `json:"error"`
	ID    string        `json:"id"`
}

// JSONRPCError will not be nil while failed
type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err *JSONRPCError) Error() error {
	return fmt.Errorf("[code: %d] %s", err.Code, err.Message)
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// BitcoinRPC struct need be initialed
type BitcoinRPC struct {
	Debug   bool
	Host    string
	Version string
	ID      string
	User    string
	Pass    string
}

// Call certain method
func (btc *BitcoinRPC) Call(method string, params []interface{}) (interface{}, error) {
	if btc.Version == "" {
		btc.Version = "1.0"
	}
	if btc.ID == "" {
		btc.ID = "jsonrpc"
	}

	rpcReq := jsonRPCRequest{Version: btc.Version, ID: btc.ID, Method: method, Params: params}
	rpcReqData, err := json.Marshal(rpcReq)
	if err != nil {
		return nil, err
	}

	if btc.Debug {
		fmt.Printf("[json-rpc-debug] Call: %s\n", rpcReqData)
	}

	req, err := http.NewRequest("POST", btc.Host, bytes.NewBuffer(rpcReqData))
	req.Header.Set("content-type", "text/json")
	req.Header.Add("Authorization", "Basic "+basicAuth(btc.User, btc.Pass))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if btc.Debug {
		fmt.Printf("[json-rpc-debug] Resp: %s\n", body)
	}

	var result JSONRPCResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Error != nil {
		return nil, result.Error.Error()
	}
	return result.Data, nil
}

// GetNewAddress returns a new Qtum address for receiving payments
func (btc *BitcoinRPC) GetNewAddress() (string, error) {
	ret, err := btc.Call("getnewaddress", []interface{}{})
	if err != nil {
		return "", err
	}

	return ret.(string), nil
}

// GetAddressesByAccount returns the list of addresses for the given account.
func (btc *BitcoinRPC) GetAddressesByAccount(account string) ([]string, error) {
	ret, err := btc.Call("getaddressesbyaccount", []interface{}{account})
	if err != nil {
		return nil, err
	}

	addresses := []string{}
	for _, v := range ret.([]interface{}) {
		s := v.(string)
		addresses = append(addresses, s)
	}
	return addresses, nil
}

// GetAccountAddress returns the current address for receiving payments to this account.
// If the account don't exist, it creates both the account and address.
// Once a payment has been received to the address, future calls to this RPC for
// the same account will return a different address.
func (btc *BitcoinRPC) GetAccountAddress(account string) (string, error) {
	ret, err := btc.Call("getaccountaddress", []interface{}{account})
	if err != nil {
		return "", err
	}

	return ret.(string), nil
}

// GetReceivedByAddress returns the total amount received by the given address
// in transactions with at least 1 confirmation.
func (btc *BitcoinRPC) GetReceivedByAddress(address string) (float64, error) {
	ret, err := btc.Call("getreceivedbyaddress", []interface{}{address})
	if err != nil {
		return 0, err
	}

	return ret.(float64), nil
}

// GetReceivedByAccount returns the total amount received by addresses in a particular
// account from transactions with the specified number of confirmations.
func (btc *BitcoinRPC) GetReceivedByAccount(account string) (float64, error) {
	ret, err := btc.Call("getreceivedbyaccount", []interface{}{account})
	if err != nil {
		return 0, err
	}

	return ret.(float64), nil
}

// SendToAddress send an amount to a given address, returns its transaction id.
func (btc *BitcoinRPC) SendToAddress(address string, amount float64) (string, error) {
	ret, err := btc.Call("sendtoaddress", []interface{}{address})
	if err != nil {
		return "", err
	}

	return ret.(string), nil
}

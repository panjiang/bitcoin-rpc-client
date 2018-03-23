package rpcclient

import (
	"fmt"
	"testing"
)

func TestBitcoinRPC(t *testing.T) {
	rpc := BitcoinRPC{
		Host:    "http://127.0.0.1:13889",
		ID:      "777bingo",
		Version: "1.0",
		User:    "test",
		Pass:    "test1234",
	}

	r, err := rpc.Call("getinfo", nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("#1 info: %+v\n", r)

	r, err = rpc.Call("getaddressesbyaccount", []interface{}{"product.1"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("#2 addresses list: %+v\n", r)

	r, err = rpc.Call("getaccountaddress", []interface{}{"product.1"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("#3 address: %+v\n", r)

	r, err = rpc.Call("sendtoaddress", []interface{}{"qcqCSPnAMJBMszPSu2puFNPtMivaYQG6Ep", 21.12345678})
	if err != nil {
		if err.Error() == "[code: -6] Insufficient funds" {
			fmt.Println("#4 Insufficient funds")
		} else {
			panic(err)
		}
	} else {
		fmt.Printf("#4 txid: %+v\n", r)
	}

	r, err = rpc.Call("getreceivedbyaddress", []interface{}{"qcqCSPnAMJBMszPSu2puFNPtMivaYQG6Ep", 1})
	if err != nil {
		panic(err)
	}
	fmt.Printf("#5 received: %+v\n", r)

	r, err = rpc.Call("getblockcount", nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("#6 blockcount: %+v\n", r)

	r, err = rpc.Call("getblockcount", []interface{}{1, "a"})
	if err != nil {
		panic(err)
	}
}

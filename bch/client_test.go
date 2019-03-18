package bch

import (
	"fmt"
	"testing"
)

func Test_Client(t *testing.T) {
	c, err := New("192.168.1.152", 8332, "test", "test")
	if err != nil {
		fmt.Println(err)
	}
	// li, err := c.ListTransaction("bchtest:qp635vnwt6dtuceft6a0zn2t5xgye9fy3v6nu52mw7")
	// fmt.Println(li, ":", err)

	li, err := c.GetBalanceOfAddr("bchtest:qp635vnwt6dtuceft6a0zn2t5xgye9fy3v6nu52mw7")
	fmt.Println(li, ":", err)
	// addr, err := c.GetNewAddress("mr-bch")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("addr:", addr)

	{
		balance, err := c.GetBalance("mr-bch1")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("balance:", balance)
	}
	{
		balance, err := c.GetBalance("mr-bch2")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("balance:", balance)
	}
	// SendBtc(fromAccount, to string, amount string) (string, error)
	{
		// balance, err := c.SendBtc("mr-bch2", "msDrFfNZP6VjsQC92t1LgJBc8Q8qhJkKVG", "0.001")
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// fmt.Println("balance:", balance)
	}
}

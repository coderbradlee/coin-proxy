package main

import (
	"github.com/lzxm160/coin-proxy/xlog"
	"fmt"
	"log"
	"github.com/lzxm160/coin-proxy/martini"
	"net/http"
)

const (
	etc_dir  = "./keys/etc"
	neo_dir  = "./keys/neo"
	true_dir = "./keys/true"
	trx_dir  = "./keys/trx"
	eos_dir  = "./keys/eos"
	ada_dir  = "./keys/ada/"
	iota_dir = "./keys/iota/"
	_dir     = "./keys"
	// ont_dir  = "./keys/ont/"
	ipfs_dir   = "./keys/ipfs/"
	bchabc_dir = "./keys/bchabc/"
	bchsv_dir  = "./keys/bchsv/"
)

func init() {
	//xlog.Create(_dir)
	//xlog.Create(eos_dir)
	//xlog.Create(trx_dir)
	//xlog.Create(etc_dir)
	//xlog.Create(true_dir)
	//xlog.Create(ada_dir)
	//xlog.Create(neo_dir)
	//xlog.Create(iota_dir)
	//// xlog.Create(ont_dir)
	//xlog.Create(ipfs_dir)
	//xlog.Create(bchabc_dir)
	//xlog.Create(bchsv_dir)
}

type EosConf struct {
	NODE_PRODUCER_NAME string
	NODE_PUB_KEY       string
	// ENV_EOS_SRC_PATH			 string
	// ENV_EOSGO_PATH				 string
	API_PORT       int
	API_URL        string
	LOCAL_API_PORT int
	LOCAL_API_URL  string
	API_METHOD     string
	// LOGGING_MODE                 string
	WALLET_NAME                  string
	WALLET_PRIV_KEY              string
	TRANSACTION_EXPIRATION_DELAY int
	Account_History              string
}
type YoyowConf struct {
	NODE_PRODUCER_NAME           string
	NODE_PUB_KEY                 string
	API_PORT                     int
	API_URL                      string
	LOCAL_API_PORT               int
	LOCAL_API_URL                string
	API_METHOD                   string
	WALLET_NAME                  string
	WALLET_PRIV_KEY              string
	TRANSACTION_EXPIRATION_DELAY int
}
type OmniConf struct {
	Host       string `json:"host,omitempty"`
	Port       int    `json:"port,omitempty"`
	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	Account    string `json:"account,omitempty"`
	Tokenid    int64  `json:"tokenid,omitempty"`
	WalletPass string `json:"walletPass,omitempty"`
	Net        string `json:"net,omitempty"`
}
type BtcConf struct {
	Host       string `json:"host,omitempty"`
	Port       int    `json:"port,omitempty"`
	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	Account    string `json:"account,omitempty"`
	Tokenid    int64  `json:"tokenid,omitempty"`
	WalletPass string `json:"walletPass,omitempty"`
	Net        string `json:"net,omitempty"`
}
type BchConf struct {
	Host       string `json:"host,omitempty"`
	Port       int    `json:"port,omitempty"`
	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	Account    string `json:"account,omitempty"`
	Tokenid    int64  `json:"tokenid,omitempty"`
	WalletPass string `json:"walletPass,omitempty"`
	Net        string `json:"net,omitempty"`
}
type ZcashConf struct {
	Host       string `json:"host,omitempty"`
	Port       int    `json:"port,omitempty"`
	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	Account    string `json:"account,omitempty"`
	WalletPass string `json:"walletPass,omitempty"`
}
type LtcConf struct {
	Host       string `json:"host,omitempty"`
	Port       int    `json:"port,omitempty"`
	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	Account    string `json:"account,omitempty"`
	WalletPass string `json:"walletPass,omitempty"`
	Net        string `json:"net,omitempty"`
}
type TrxConf struct {
	Host       string `json:"host,omitempty"`
	Port       string `json:"port,omitempty"`
	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	Account    string `json:"account,omitempty"`
	Tokenid    int64  `json:"tokenid,omitempty"`
	WalletPass string `json:"walletPass,omitempty"`
	Net        string `json:"net,omitempty"`
	Private    string `json:"private,omitempty"`
	Method     string `json:"method,omitempty"`
	Localurl   string `json:"localurl,omitempty"`
}
type XrpConf struct {
	Host     string `json:"host,omitempty"`
	Port     string `json:"port,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Private  string `json:"private,omitempty"`
}
type XlmConf struct {
	Url         string `json:"url,omitempty"`
	Seed        string `json:"seed,omitempty"`
	Public      string `json:"public,omitempty"`
	Networkpass string `json:"networkpass,omitempty"`
}
type EtcConf struct {
	Host        string `json:"host,omitempty"`
	Port        string `json:"port,omitempty"`
	AddressPass string `json:"addressPass,omitempty"`
}
type OntConf struct {
	Host            string `json:"host,omitempty"`
	Port            string `json:"port,omitempty"`
	EncryptFilePass string `json:"encryptFilePass,omitempty"`
	Wallet          string `json:"wallet,omitempty"`
}
type TrueConf struct {
	Host        string `json:"host,omitempty"`
	Port        string `json:"port,omitempty"`
	AddressPass string `json:"addressPass,omitempty"`
}
type IotaConf struct {
	Host   string `json:"host,omitempty"`
	Port   string `json:"port,omitempty"`
	Method string `json:"method,omitempty"`
	Seed   string `json:"seed,omitempty"`
}
type NeoConf struct {
	Url        string `json:"url,omitempty"`
	WalletPass string `json:"walletPass,omitempty"`
	NeoAssetID string `json:"neoAssetID,omitempty"`
	GasAssetID string `json:"gasAssetID,omitempty"`
}
type DashConf struct {
	Host       string `json:"host,omitempty"`
	Port       int    `json:"port,omitempty"`
	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	Account    string `json:"account,omitempty"`
	WalletPass string `json:"walletPass,omitempty"`
}
type XmrConf struct {
	Url      string `json:"url,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}
type XmcConf struct {
	Url      string `json:"url,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}
type IpfsConf struct {
	Url    string `json:"url,omitempty"`
	Peerid string `json:"peerid,omitempty"`
}
type Conf struct {
	ListenPort string    `json:"listenPort,omitempty"`
	Httpuser   string    `json:"httpuser,omitempty"`
	Httppass   string    `json:"httppass,omitempty"`
	Omni       OmniConf  `json:"omni,omitempty"`
	Eos        EosConf   `json:"eos,omitempty"`
	Yoyow      YoyowConf `json:"yoyow,omitempty"`
	Bchabc     BchConf   `json:"bchabc,omitempty"`
	Bchsv      BchConf   `json:"bchsv,omitempty"`
	Trx        TrxConf   `json:"trx,omitempty"`
	Xrp        XrpConf   `json:"xrp,omitempty"`
	Ltc        LtcConf   `json:"ltc,omitempty"`
	Etc        EtcConf   `json:"etc,omitempty"`
	Ada        AdaConf   `json:"ada,omitempty"`
	True       TrueConf  `json:"true,omitempty"`
	Neo        NeoConf   `json:"neo,omitempty"`
	Zcash      ZcashConf `json:"zcash,omitempty"`
	Dash       DashConf  `json:"dash,omitempty"`
	Iota       IotaConf  `json:"iota,omitempty"`
	Ont        OntConf   `json:"ont,omitempty"`
	Xlm        XlmConf   `json:"xlm,omitempty"`
	Xmr        XmrConf   `json:"xmr,omitempty"`
	Xmc        XmcConf   `json:"xmc,omitempty"`
	Ipfs       IpfsConf  `json:"ipfs,omitempty"`
	Btc        BtcConf   `json:"btc,omitempty"`
}
type AdaConf struct {
	API_PORT         int
	API_URL          string
	API_METHOD       string
	WALLET_PRIV_KEY  string
	AccountIndex     string
	WalletId         string
	SpendingPassword string
	Capath           string
}

var cfg Conf //proxy.Config

func logPanics(function func(http.ResponseWriter,
	*http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if x := recover(); x != nil {
				log.Println(fmt.Sprintf("[%v] caught panic: %v", request.RemoteAddr, x))
				fmt.Println(fmt.Sprintf("[%v] caught panic: %v", request.RemoteAddr, x))
			}
		}()
		function(writer, request)
	}
}

func startMartini() {
	m := martini.Classic()
	//m.Use(auth_BasicFunc(func(username, password string) bool {
	//	return auth_SecureCompare(username, cfg.Httpuser) && auth_SecureCompare(password, cfg.Httppass)
	//}))
	//m.Post("/", logPanics(Handler))                //for usdt
	//m.Post("/usdt", logPanics(Handler))            //for usdt
	//m.Post("/eos", logPanics(EosHandler))          //for eos
	//m.Post("/yoyow", logPanics(YoyowHandler))      //for yoyow
	//m.Post("/bchabc", logPanics(BchabcHandler))    //for bch
	//m.Post("/bchsv", logPanics(BchsvHandler))      //for bch
	//m.Post("/trx", logPanics(TrxHandler))          //for trx
	//m.Post("/xrp", logPanics(XrpHandler))          //for xrp
	//m.Post("/ltc", logPanics(LtcHandler))          //for ltc
	//m.Post("/etc", logPanics(EtcHandler))          //for etc
	//m.Post("/ada", logPanics(AdaHandler))          //for ada
	//m.Post("/true", logPanics(TrueHandler))        //for true
	//m.Post("/neo", logPanics(NeoHandler))          //for neo
	//m.Post("/zcash", logPanics(ZcashHandler))      //for zcash
	//m.Post("/dash", logPanics(DashHandler))        //for dash
	//m.Post("/iota", logPanics(IotaHandler))        //for iota
	//m.Post("/ont", logPanics(OntHandler))          //for ont
	//m.Post("/xlm", logPanics(XlmHandler))          //for xlm
	//m.Post("/xmr", logPanics(XmrHandler))          //for xlr
	//m.Post("/xmc", logPanics(XmcHandler))          //for xlc
	//m.Post("/wallet_ipfs", logPanics(IpfsHandler)) //for ipfs app wallet
	//m.Post("/wallet_btc", logPanics(BtcHandler))   //for btc app wallet
	m.Post("/", logPanics(ethminerHandler))
	m.RunOnAddr(cfg.ListenPort)
}

func main() {
	xlog.XX()
	if !LoadConfig("config.toml", &cfg) {
		return
	}
	log.Println(cfg)
	startMartini()
	quit := make(chan bool)
	<-quit
}

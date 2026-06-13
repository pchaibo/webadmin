package binan

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var Strongpassword string = "ZH123456++"

func Ceatkeystore() (path string, address string) {
	ksDir := "./keystore"
	if err := os.MkdirAll(ksDir, 0700); err != nil {
		log.Fatal(err)
	}

	ks := keystore.NewKeyStore(ksDir, keystore.StandardScryptN, keystore.StandardScryptP)

	passphrase := Strongpassword
	account, err := ks.NewAccount(passphrase)
	if err != nil {
		log.Fatal(err)
	}
	//firname := filepath.Join(ksDir, account.URL.Path)
	fmt.Println("Saved keystore file to:", account.URL.Path)
	fmt.Println("Address:", account.Address.Hex())
	address = account.Address.Hex()
	path = account.URL.Path
	//keystoreadd(account.URL.Path, passphrase)
	return

	// 如果你想从私钥生成 keystore 文件也可以：
	// privKey, _ := crypto.GenerateKey()
	// _, err = ks.ImportECDSA(privKey, passphrase)

}

func Keystoreadd(keystoreFile string, password string) (privateKey *ecdsa.PrivateKey, address string) {
	// keystore 文件路径
	//keystoreFile := "./keystore/UTC--2025-11-28T07-17-38.451347000Z--c0b3fd7884171d7a51b2cf170d7353d2f297c843"
	//password := "strong-password-123"

	// 读取文件内容
	keyjson, err := os.ReadFile(keystoreFile)
	if err != nil {
		log.Fatal(err)
	}

	// 解锁 keystore（得到 ECDSA 私钥）
	key, err := keystore.DecryptKey(keyjson, password)
	if err != nil {
		log.Fatal(err)
	}

	privateKey = key.PrivateKey
	address = key.Address.Hex()
	fmt.Println("Private key (hex):", fmt.Sprintf("%x", crypto.FromECDSA(privateKey)))
	fmt.Println("Address:", address)
	return

	// 注意：privateKey 就是 *ecdsa.PrivateKey，可直接用于签名交易
}

// 转账
func Transfer() {
	// BSC USDT 合约地址（BEP20）
	var usdtAddress = common.HexToAddress("0x55d398326f99059fF775485246999027B3197955")

	// 精度 18
	var decimals = big.NewInt(1e18)

	// ---- ① 读取 keystore 并解密 ----
	keystoreFile := "./keystore/UTC--2025-11-28T08-37-03.331081800Z--85ed4def2dda5723065b73f4b7c9699905618e87"
	//password := "strong-password-123"
	password := Strongpassword

	keyJson, err := os.ReadFile(keystoreFile)
	if err != nil {
		log.Fatal(err)
	}

	key, err := keystore.DecryptKey(keyJson, password)
	if err != nil {
		log.Fatal(err)
	}
	privateKey := key.PrivateKey
	fromAddress := key.Address

	fmt.Println("From:", fromAddress.Hex())

	// ---- ② RPC 连接 BSC ----
	client, err := ethclient.Dial("https://bsc-dataseed.binance.org/")
	if err != nil {
		log.Fatal(err)
	}

	// ---- ③ 获取 nonce ----
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	// ---- ④ 设置 gasPrice、gasLimit ----
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 收款地址
	toAddress := common.HexToAddress("0xReceiverAddressHere")

	// 转账金额（例如 10 USDT）
	amount := new(big.Int).Mul(big.NewInt(10), decimals)

	// ---- ⑤ 构建 ERC20 transfer data ----
	transferFnSignature := []byte("transfer(address,uint256)")
	hash := crypto.Keccak256Hash(transferFnSignature)
	methodID := hash.Bytes()[0:4]

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	// gasLimit 一般 60,000 足够
	gasLimit := uint64(60000)

	// ---- ⑥ 构建交易 ----
	chainID := big.NewInt(56) // BSC 主网
	tx := types.NewTransaction(
		nonce,
		usdtAddress,
		big.NewInt(0), // ERC20 转账主币 value 必须为 0
		gasLimit,
		gasPrice,
		data,
	)

	// ---- ⑦ 签名交易 ----
	signer := types.NewEIP155Signer(chainID)
	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// ---- ⑧ 发送交易 ----
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("USDT transfer sent!")
	fmt.Println("TX Hash:", signedTx.Hash().Hex())

}

// 普通生成
func Cateadd() {
	// 生成随机私钥
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	// 私钥的 hex（64 字节）表示（不包含 0x 前缀）
	privBytes := crypto.FromECDSA(privateKey)
	fmt.Printf("Private key (hex): %x\n", privBytes)

	// 从私钥派生公钥（ECDSA）
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	// 得到地址（20 字节），并以 EIP-55 checksum 格式展示
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Printf("Address: %s\n", address.Hex()) // 0x 开头
	// 可选：小写无校验格式
	fmt.Printf("Address (no checksum): %s\n", common.BytesToAddress(address.Bytes()).Hex())
}

func Keytoadd() {
	hexKey := "4c0883a69... (64 hex chars)"
	privateKey, err := crypto.HexToECDSA(hexKey)
	if err != nil {
		// 处理错误
	}
	pubKey := privateKey.Public()
	pubKeyECDSA := pubKey.(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*pubKeyECDSA)
	fmt.Println(address.Hex())

}

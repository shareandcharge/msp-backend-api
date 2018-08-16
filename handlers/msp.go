package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/gin-gonic/gin/binding"
	"fmt"
	"github.com/motionwerkGmbH/msp-backend-api/tools"
	"github.com/motionwerkGmbH/msp-backend-api/configs"
	"encoding/json"
	"log"
	"strings"
)

func MspCreate(c *gin.Context) {

	type MspInfo struct {
		Name        string `json:"name"`
		Address1    string `json:"address_1"`
		Address2    string `json:"address_2"`
		Town        string `json:"town"`
		Postcode    string `json:"postcode"`
		MailAddress string `json:"mail_address"`
		Website     string `json:"website"`
		VatNumber   string `json:"vat_number"`
	}
	var mspInfo MspInfo

	if err := c.MustBindWith(&mspInfo, binding.JSON); err == nil {
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//check if there is already an msp registered
	rows, err := tools.DB.Query("SELECT msp_id FROM msp")
	tools.ErrorCheck(err, "msp.go", true)
	defer rows.Close()

	//check if we already have an MSP registered
	if rows.Next() {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "there's already an MSP registered on this backend"})
		return
	}

	//if not, insert a new one with ID = 1, unique.
	query := "INSERT INTO msp (msp_id, wallet, seed, name, address_1, address_2, town, postcode, mail_address, website, vat_number) VALUES (%d, '%s', '%s','%s','%s','%s','%s','%s','%s','%s','%s')"
	command := fmt.Sprintf(query, 1, "", "", mspInfo.Name, mspInfo.Address1, mspInfo.Address2, mspInfo.Town, mspInfo.Postcode, mspInfo.MailAddress, mspInfo.Website, mspInfo.VatNumber)
	tools.DB.MustExec(command)

	c.JSON(http.StatusOK, gin.H{"status": "created ok"})
}

//returns the info for the MSP
func MspInfo(c *gin.Context) {

	rows, _ := tools.DB.Query("SELECT msp_id FROM msp")
	defer rows.Close()

	//check if we already have an MSP registered
	if rows.Next() == false {
		c.JSON(http.StatusNotFound, gin.H{"error": "we couldn't find any MPS registered in the database."})
		return
	}

	msp := tools.MSP{}

	tools.DB.QueryRowx("SELECT * FROM msp LIMIT 1").StructScan(&msp)
	c.JSON(http.StatusOK, msp)
}

//returns the info for the MSP
func MspGetSeed(c *gin.Context) {

	msp := tools.MSP{}
	tools.DB.QueryRowx("SELECT * FROM msp LIMIT 1").StructScan(&msp)

	if msp.Seed == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "there isn't any seed in the msp account. Maybe you need to create the wallet first ?."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"seed": msp.Seed})
}

//generates a new wallet for the msp
func MspGenerateWallet(c *gin.Context) {

	type WalletInfo struct {
		Seed string `json:"seed"`
		Addr string `json:"address"`
	}
	var walletInfo WalletInfo

	// Leave this commented code here please
	//body := tools.GetRequest("http://localhost:3000/api/wallet/create")
	//log.Printf("<- %s", string(body))
	//err := json.Unmarshal(body, &walletInfo)
	//if err != nil {
	//	log.Panic(err)
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "ops! it's our fault. This error should never happen."})
	//	return
	//}

	config := configs.Load()
	walletInfo.Addr = config.GetString("msp.wallet_address")
	walletInfo.Seed = config.GetString("msp.wallet_seed")

	//update the db for MSP
	query := "UPDATE msp SET wallet='%s', seed='%s' WHERE msp_id = 1"
	command := fmt.Sprintf(query, walletInfo.Addr, walletInfo.Seed)
	tools.DB.MustExec(command)

	//update the ~/.sharecharge/config.json
	configs.UpdateBaseAccountSeedInSCConfig(walletInfo.Seed)

	c.JSON(http.StatusOK, walletInfo)
}

//Gets the history for the MSP
func MSPHistory(c *gin.Context) {

	config := configs.Load()
	mspAddress := config.GetString("msp.wallet_address")

	type History struct {
		Block           int    `json:"block" db:"block"`
		FromAddr        string `json:"from_addr" db:"from_addr"`
		ToAddr          string `json:"to_addr" db:"to_addr"`
		Amount          uint64 `json:"amount" db:"amount"`
		Currency        string `json:"currency" db:"currency"`
		CreatedAt       uint64 `json:"created_at" db:"created_at"`
		TransactionHash string `json:"transaction_hash" db:"transaction_hash"`
	}
	var histories []History

	var transactions []tools.TxTransaction
	err := tools.MDB.Select(&transactions, "SELECT * FROM transactions WHERE (to_addr = ? OR from_addr = ?) ORDER BY blockNumber DESC", mspAddress, mspAddress)
	tools.ErrorCheck(err, "cpo.go", false)

	for _, tx := range transactions {
		if tx.Value == "0x0" {
			//we have a contract tx

			var txResponse tools.TxReceiptResponse
			err := tools.MDB.QueryRowx("SELECT * FROM transaction_receipts WHERE transactionHash = ?", tx.Hash).StructScan(&txResponse)
			tools.ErrorCheck(err, "cpo.go", false)
			calculatedGas := tools.HexToUInt(txResponse.GasUsed) * tools.HexToUInt(tx.GasPrice)
			histories = append(histories, History{Block: tx.BlockNumber, FromAddr: tx.From, ToAddr: tx.To, Amount: calculatedGas, Currency: "wei", CreatedAt: tx.Timestamp, TransactionHash: tx.Hash})

		} else {
			//we have eth transfer
			histories = append(histories, History{Block: tx.BlockNumber, FromAddr: tx.From, ToAddr: tx.To, Amount: tools.HexToUInt(tx.Value), Currency: "wei", CreatedAt: tx.Timestamp, TransactionHash: tx.Hash})
		}
	}

	c.JSON(http.StatusOK, histories)

}

// Shows the latest transaction from a driver
// https://trello.com/c/vokmESXz/182-see-latest-charging-sessions-1
func GetDriverHistory(c *gin.Context) {

	driverAddr := c.Param("addr")



	body := tools.GETRequest("http://localhost:3000/api/cdr/info") //+ ?controller=driverAddr

	var cdrs []tools.CDR
	err := json.Unmarshal(body, &cdrs)
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ops! it's our fault. This error should never happen."})
		return
	}

	var cdrsOutput []tools.CDR
	for _, cdr := range cdrs {
		cdr.Currency = "Charge & Fuel Token"

		log.Printf("processing.. %s, %s", driverAddr, cdr.Controller)

		//TODO: after filtering works, remove this part
		//filter by the driver
		if strings.ToLower(cdr.Controller) == strings.ToLower(driverAddr) {
			log.Println("adding")
			cdrsOutput = append(cdrsOutput, cdr)
		}
	}

	if len(cdrsOutput) == 0 {
		c.JSON(http.StatusOK, []string{})
		return
	}

	c.JSON(http.StatusOK, cdrsOutput)

}

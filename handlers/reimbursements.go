package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/motionwerkGmbH/msp-backend-api/tools"
	"net/http"
	"strings"
)

func SetReimbursementStatus(c *gin.Context) {
	reimbursementId := c.Param("reimbursement_id")
	status := c.Param("status")

	count := 0
	row := tools.MDB.QueryRow("SELECT COUNT(*) as count FROM reimbursements WHERE reimbursement_id = '" + reimbursementId + "'")
	row.Scan(&count)

	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "could not find any reimbursement with this id"})
		return
	}

	//update status from the database
	query := "UPDATE reimbursements SET status='%s' WHERE reimbursement_id = '%s'"
	command := fmt.Sprintf(query, status, reimbursementId)
	tools.MDB.MustExec(command)

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func ViewCDRs(c *gin.Context) {

	reimbursementId := c.Param("reimbursement_id")

	type Reimbursement struct {
		CdrRecords string `json:"cdr_records" db:"cdr_records"`
	}
	var reimbursement Reimbursement

	err := tools.MDB.QueryRowx("SELECT cdr_records FROM reimbursements WHERE reimbursement_id = ?", reimbursementId).StructScan(&reimbursement)
	tools.ErrorCheck(err, "cpo.go", false)

	if reimbursement.CdrRecords == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "no reimbursements found"})
		return
	}

	output := strings.Replace(reimbursement.CdrRecords, "\\", "", -1)

	c.JSON(200, output)

}

// lists all reimbursements filtered by status query parameter
func ListReimbursements(c *gin.Context) {

	status := c.DefaultQuery("status", "pending")

	type CpoInfo struct {
		Name          string `json:"name"`
		PublicKey     string `json:"public_key"`
		ServerAddress string `json:"server_addr"`
	}

	type Reimbursement struct {
		Id              int    `json:"id" db:"id"`
		MspName         string `json:"msp_name" db:"msp_name"`
		CpoName         string `json:"cpo_name" db:"cpo_name"`
		Amount          int    `json:"amount" db:"amount"`
		Currency        string `json:"currency" db:"currency"`
		Timestamp       int    `json:"timestamp" db:"timestamp"`
		Status          string `json:"status" db:"status"`
		ReimbursementId string `json:"reimbursement_id" db:"reimbursement_id"`
		CdrRecords      string `json:"cdr_records" db:"cdr_records"`
	}

	type AllReimbursements struct {
		CpoName        string          `json:"cpo_name"`
		ServerUrl      string          `json:"server_url"`
		Reimbursements []Reimbursement `json:"reimbursements"`
	}

	var all_reimbursements []AllReimbursements

	var allCpos []CpoInfo

	allCpos = append(allCpos, CpoInfo{"Innogy Office", "0x7b0f2b531c018d4269a95561cfb4e038a7e3c8dc", "http://52.57.155.233:9090/api/v1"})
	allCpos = append(allCpos, CpoInfo{"Cpo2", "0x7b0f2b531c018d4269a95561cfb4e038a7e3c8dc", "https://innogy-api.shareandcharge.com/api/v1"})

	for _, cpo := range allCpos {
		logrus.Info("Processing CPO: %s", cpo.Name)

		body := tools.GETRequest(cpo.ServerAddress + "/cpo/payment/reimbursements/" + status)
		var reimb []Reimbursement
		err := json.Unmarshal(body, &reimb)
		tools.ErrorCheck(err, "msp.go", true)

		all_reimbursements = append(all_reimbursements, AllReimbursements{CpoName: cpo.Name, Reimbursements: reimb, ServerUrl: cpo.ServerAddress})

	}

	c.JSON(http.StatusOK, all_reimbursements)

	//check if there is already an msp registered
	//rows, err := tools.DB.Query("SELECT msp_id FROM msp")
	//tools.ErrorCheck(err, "msp.go", true)
	//defer rows.Close()
	//
	////check if we already have an MSP registered
	//if rows.Next() {
	//	c.JSON(http.StatusNotAcceptable, gin.H{"error": "there's already an MSP registered on this backend"})
	//	return
	//}
	//
	////if not, insert a new one with ID = 1, unique.
	//query := "INSERT INTO msp (msp_id, wallet, seed, name, address_1, address_2, town, postcode, mail_address, website, vat_number) VALUES (%d, '%s', '%s','%s','%s','%s','%s','%s','%s','%s','%s')"
	//command := fmt.Sprintf(query, 1, "", "", mspInfo.Name, mspInfo.Address1, mspInfo.Address2, mspInfo.Town, mspInfo.Postcode, mspInfo.MailAddress, mspInfo.Website, mspInfo.VatNumber)
	//tools.DB.MustExec(command)
	//
	//c.JSON(http.StatusOK, gin.H{"status": "created ok"})
}

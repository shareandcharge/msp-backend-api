package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/motionwerkGmbH/msp-backend-api/tools"
	"net/http"
		"encoding/json"
	"log"
	"bytes"
	"encoding/csv"
)

//sets a reimbursement status (ex: complete)
func SetReimbursementStatus(c *gin.Context) {
	reimbursementId := c.Param("reimbursement_id")
	status := c.Param("status")

	count := 0
	row := tools.MDB.QueryRow("SELECT COUNT(*) as count FROM reimbursements WHERE reimbursement_id = '" + reimbursementId + "'")
	row.Scan(&count)

	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "could not find any reimbursement with this id "+ reimbursementId})
		return
	}

	//update status from the database
	query := "UPDATE reimbursements SET status='%s' WHERE reimbursement_id = '%s'"
	command := fmt.Sprintf(query, status, reimbursementId)
	tools.MDB.MustExec(command)

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

//shows the CDR records of a reimbursement
func ViewCDRs(c *gin.Context) {

	reimbursementId := c.Param("reimbursement_id")

	type CDR struct {
		EvseID           string `json:"evseId"`
		ScID             string `json:"scId"`
		Controller       string `json:"controller"`
		Start            string `json:"start"`
		End              string `json:"end"`
		FinalPrice       string `json:"finalPrice"`
		TokenContract    string `json:"tokenContract"`
		Tariff           string `json:"tariff"`
		ChargedUnits     string `json:"chargedUnits"`
		ChargingContract string `json:"chargingContract"`
		TransactionHash  string `json:"transactionHash"`
		Currency         string `json:"currency"`
	}

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

	var cdrs []CDR
	err = json.Unmarshal([]byte(reimbursement.CdrRecords), &cdrs)
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ops! it's our fault. This error should never happen."})
		return
	}


	b := &bytes.Buffer{} // creates IO Writer
	wr := csv.NewWriter(b) // creates a csv writer that uses the io buffer.


	wr.Write([]string{"evseId", "scId","controller","start","end","finalPrice","tokenContract","tariff","chargedUnits","chargingContract","transactionHash","currency"})
	for _, cdr := range cdrs {
		wr.Write([]string{cdr.EvseID, cdr.ScID, cdr.Controller, cdr.Start, cdr.End, cdr.FinalPrice, cdr.TokenContract, cdr.Tariff, cdr.ChargedUnits, cdr.ChargingContract, cdr.TransactionHash, cdr.Currency})
	}
	wr.Flush()

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=cdrs_"+reimbursementId+".csv")
	c.Data(http.StatusOK, "text/csv", b.Bytes())

}

// lists all reimbursements filtered by status query parameter
func ListReimbursements(c *gin.Context) {

	status := c.DefaultQuery("status", "pending")

	var reimbursements []tools.Reimbursement

	err := tools.MDB.Select(&reimbursements, "SELECT * FROM reimbursements WHERE status = ? ORDER BY timestamp DESC", status)
	tools.ErrorCheck(err, "cpo.go", false)

	if len(reimbursements) == 0 {
		log.Println("no reimbursements found with status " + status)
		c.JSON(200, []string{})
		return
	}

	var output []tools.Reimbursement
	for k, reimb := range reimbursements {
		reimb.Index = k
		output = append(output, reimb)
	}

	c.JSON(http.StatusOK, output)

}

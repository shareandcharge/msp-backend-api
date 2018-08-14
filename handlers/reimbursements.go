package handlers

import (
	"fmt"
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

	type Reimbursement struct {
		Index           int    `json:"index"`
		Id              int    `json:"id" db:"id"`
		MspName         string `json:"msp_name" db:"msp_name"`
		CpoName         string `json:"cpo_name" db:"cpo_name"`
		Amount          int    `json:"amount" db:"amount"`
		Currency        string `json:"currency" db:"currency"`
		Timestamp       int    `json:"timestamp" db:"timestamp"`
		Status          string `json:"status" db:"status"`
		ReimbursementId string `json:"reimbursement_id" db:"reimbursement_id"`
		CdrRecords      string `json:"cdr_records" db:"cdr_records"`
		ServerAddr      string `json:"server_addr" db:"server_addr"`
	}

	var reimbursements []Reimbursement

	err := tools.MDB.Select(&reimbursements, "SELECT * FROM reimbursements WHERE status = ? ORDER BY timestamp DESC", status)
	tools.ErrorCheck(err, "cpo.go", false)

	var output []Reimbursement
	for k, reimb := range reimbursements {
		reimb.Index = k
		output = append(output, reimb)
	}

	c.JSON(http.StatusOK, output)

}

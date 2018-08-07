package tools

import (
	"github.com/spf13/viper"
	"net/http"
	"io/ioutil"
	log "github.com/Sirupsen/logrus"
	"time"
	"context"
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"crypto/sha1"
	"fmt"
	"os/exec"
	"bufio"
)

//read the config file, helper function
func ReadConfig(filename string, defaults map[string]interface{}) (*viper.Viper, error) {
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	v.SetConfigName(filename)
	v.AddConfigPath("./configs")
	v.AutomaticEnv()
	err := v.ReadInConfig()
	return v, err
}

// a general get request with 200 seconds timeout (yep, 200!)
func GETRequest(url string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Panicf("%v", err)
		return nil
	}

	ctx, cancel := context.WithTimeout(req.Context(), 200*time.Second)
	defer cancel()

	req = req.WithContext(ctx)

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		log.Panicf("%v", err)
		return nil
	}

	if contents, err := ioutil.ReadAll(res.Body); err == nil {
		return contents
	}
	return nil
}

//general POST request
func POSTRequest(url string, payload []byte) ([]byte, error) {


	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Panicf("%v", err)
		return nil, err
	}

	if contents, err := ioutil.ReadAll(resp.Body); err == nil {
		log.Info("POST Request Returned >>> " + string(contents))
		return contents, nil
	}
	return nil, err
}

// general PUT request
func PUTRequest(url string, payload []byte) ([]byte, error) {


	body := bytes.NewReader(payload)

	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		log.Panicf("%v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panicf("%v", err)
		return nil, err
	}
	defer resp.Body.Close()


	if contents, err := ioutil.ReadAll(resp.Body); err == nil {
		log.Info("PUT Request Returned >>> " + string(contents))
		return contents, nil
	}
	return nil, err
}

// general DELETE request
func DELETERequest(url string) ([]byte, error) {

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Panicf("%v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panicf("%v", err)
		return nil, err
	}
	defer resp.Body.Close()


	if contents, err := ioutil.ReadAll(resp.Body); err == nil {
		log.Info("DELETE Request Returned >>> " + string(contents))
		return contents, nil
	}
	return nil, err

}



//general POST request
func POSTJsonRequest(url string, jsonMap map[string]interface{}) ([]byte, error) {

	jsonValue, err := json.Marshal(jsonMap)
	if err != nil {
		log.Panicf("%v", err)
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Panicf("%v", err)
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)

	log.Printf("%s we got ",string(b))

	if err != nil {
		log.Panicf("%v", err)
		return nil, err
	}

	return b, nil
}



// quick function to check for an error and, optionally terminate the program
func ErrorCheck(err error, where string, kill bool) {
	if err != nil {
		if kill {
			log.WithError(err).Fatalln("Script Terminated")
		} else {
			log.WithError(err).Warnf("  >>---> error at >>---> %s\n", where)
		}
	}
}

//convert hex to int
func HexToInt(number string) int64 {
	if number[0:2] == "0x" {
		number = number[2:]
	}
	i, err := strconv.ParseInt(number, 16, 0)
	if err != nil {
		panic(err)
	}
	return i
}


//convert int to hex
func UIntToHex(number uint64) string {
	return "0x" + strconv.FormatUint(number, 16)

}

//convert hex to int
func HexToUInt(hexStr string) uint64 {
	// remove 0x suffix if found in the input string
	cleaned := strings.Replace(hexStr, "0x", "", -1)

	// base 16 for hexadecimal
	result, _ := strconv.ParseUint(cleaned, 16, 64)
	return uint64(result)
}


//generate sha1 hash from interface{}
func GetSha1Hash(payload interface{}) string {

	out, err := json.Marshal(payload)
	if err != nil {
		panic(err)
		return ""
	}

	algorithm := sha1.New()
	algorithm.Write(out)
	return fmt.Sprintf("%x", algorithm.Sum(nil))
}

// wkhtmltopdf needs to be installed
func GeneratePdf(fromFile string, toFile string) error {

	cmd := exec.Command("wkhtmltopdf", fromFile, toFile)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			log.Printf("%s\n", scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

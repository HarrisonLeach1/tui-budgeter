package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/HarrisonLeach1/xero-tui/internal/api/models"
	"github.com/spf13/viper"
)

func GetProfitAndLossStatement(fromDate string, toDate string) (models.Report, error) {
	return fetchReport("ProfitAndLoss", fromDate, toDate)
}
func GetBankSummary(fromDate string, toDate string) (models.Report, error) {
	return fetchReport("BankSummary", fromDate, toDate)
}

func fetchReport(reportType string, fromDate string, toDate string) (models.Report, error) {
	accessToken, tenantId := getHeaders()

	url := fmt.Sprintf("https://api.xero.com/api.xro/2.0/Reports/%s?fromDate=%s&toDate=%s", reportType, fromDate, toDate)

	// create the request and execute it
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("xero-tenant-id", tenantId)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("HTTP error: %s", err)
		return models.Report{}, err
	}

	if res.StatusCode == 401 {
		body, _ := ioutil.ReadAll(res.Body)

		fmt.Printf("err header: %s \n body: %s", res.Header, body)
		return models.Report{}, fmt.Errorf("err header: %s \n body: %s", res.Header, body)
	}

	if res.StatusCode != 200 {
		return models.Report{}, fmt.Errorf("Error calling API: %s", res.Status)
	}

	// process the response
	defer res.Body.Close()
	var responseData models.ReportResponse
	body, _ := ioutil.ReadAll(res.Body)

	// unmarshal the json into a string map
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		fmt.Print(body)
		fmt.Println()

		fmt.Printf("auth: JSON error: %s", err)
		return models.Report{}, err
	}

	return responseData.Reports[0], nil

}

func getHeaders() (string, string) {
	viper.SetConfigFile("config.yaml")
	viper.AddConfigPath("./") // optionally look for config in the working directory

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			panic("config.yaml file not found")
		} else {
			// Config file was found but another error was produced
			panic("error reading config.yaml")
		}
	}
	accessToken := viper.GetString("accesstoken")
	tenantId := viper.GetString("tenantId")

	return accessToken, tenantId
}

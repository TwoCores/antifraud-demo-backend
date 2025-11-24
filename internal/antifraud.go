package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

var (
	client http.Client

	antifraud_model_url string
)

func init() {
	client = http.Client{Timeout: 5 * time.Second}

	viper.AutomaticEnv()
	antifraud_model_url = viper.GetString("ANTIFRAUD_MODEL_URL")
}

type ModelFeatures struct {
	CstDimID                  string  `json:"cst_dim_id"`
	MonthlyOSChanges          int     `json:"monthly_os_changes"`
	MonthlyPhoneModelChanges  int     `json:"monthly_phone_model_changes"`
	LastPhoneModelCategorical string  `json:"last_phone_model_categorical"`
	LastOSCategorical         string  `json:"last_os_categorical"`
	LoginsLast7Days           int     `json:"logins_last_7_days"`
	LoginsLast30Days          int     `json:"logins_last_30_days"`
	LoginFrequency7d          float64 `json:"login_frequency_7d"`
	LoginFrequency30d         float64 `json:"login_frequency_30d"`
	FreqChange7dVsMean        float64 `json:"freq_change_7d_vs_mean"`
	Logins7dOver30dRatio      float64 `json:"logins_7d_over_30d_ratio"`
	AvgLoginInterval30d       float64 `json:"avg_login_interval_30d"`
	StdLoginInterval30d       float64 `json:"std_login_interval_30d"`
	VarLoginInterval30d       float64 `json:"var_login_interval_30d"`
	EwmLoginInterval7d        float64 `json:"ewm_login_interval_7d"`
	BurstinessLoginInterval   float64 `json:"burstiness_login_interval"`
	FanoFactorLoginInterval   float64 `json:"fano_factor_login_interval"`
	ZscoreAvgLoginInterval7d  float64 `json:"zscore_avg_login_interval_7d"`
}

type PredictResponse struct {
	FraudProbability float64 `json:"fraud_probability"`
	BlockTransaction bool    `json:"block_transaction"`
}

func PredictFraud(feats *ModelFeatures) (*PredictResponse, error) {
	b, err := json.Marshal(feats)
	if err != nil {
		return nil, err
	}

	resp, err := client.Post(fmt.Sprintf("%s/predict", antifraud_model_url), "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pr PredictResponse
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, err
	}
	return &pr, nil
}

package dto

type CountingResponse struct {
	ResponseCode    string                 `json:"response_code"`
	ResponseMessage string                 `json:"response_message"`
	Data            []CountingResponseData `json:"data"`
}

type CountingResponseData struct {
	KodeAkun    string  `json:"kode_akun"`
	NamaAkun    string  `json:"nama_akun"`
	TotalDebit  float64 `json:"total_debit"`
	TotalCredit float64 `json:"total_credit"`
	SaldoAkhir  float64 `json:"saldo_akhir"`
}

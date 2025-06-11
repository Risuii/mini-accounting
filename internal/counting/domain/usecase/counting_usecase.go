package usecase

import (
	Constants "mini-accounting/constants"
	Library "mini-accounting/library"

	CountingDTO "mini-accounting/internal/counting/delivery/dto"
)

type CountingUsecase interface {
	Index(param *CountingDTO.CountingRequest) (*CountingDTO.CountingResponse, error)
}

type CountingUsecaseImpl struct {
	library Library.Library
}

func NewCountingUsecase(
	library Library.Library,
) CountingUsecase {
	return &CountingUsecaseImpl{
		library: library,
	}
}

func (u *CountingUsecaseImpl) Index(param *CountingDTO.CountingRequest) (*CountingDTO.CountingResponse, error) {
	// MAPPING KODE AKUN
	kodeAkuns := make(map[string]bool)
	for _, v := range param.KodeAkun {
		kodeAkuns[v.KodeAkun] = true
	}

	// MAPPING JURNAL berdasarkan kode akun
	journal := make(map[string][]CountingDTO.CountingRequestData)
	for _, j := range param.Data {
		if kodeAkuns[j.KodeAkun] {
			journal[j.KodeAkun] = append(journal[j.KodeAkun], j)
		}
	}

	// HITUNG total debit & kredit per kode akun dan ambil nama akun dari entry pertama
	var responseData []CountingDTO.CountingResponseData
	for kode, entries := range journal {
		var totalDebit, totalCredit float64
		var namaAkun string

		for i, entry := range entries {
			if i == 0 {
				namaAkun = entry.NamaAkun // ambil nama akun dari entry pertama
			}

			if entry.Debit != Constants.NilString {
				parsedDebit, err := u.library.ParseFloat(u.library.ReplaceAll(entry.Debit, ",", ""), 64)
				if err != nil {
					return nil, err
				}
				totalDebit += parsedDebit
			}

			if entry.Credit != Constants.NilString {
				parsedCredit, err := u.library.ParseFloat(u.library.ReplaceAll(entry.Credit, ",", ""), 64)
				if err != nil {
					return nil, err
				}
				totalCredit += parsedCredit
			}
		}

		selisih := totalDebit - totalCredit
		responseData = append(responseData, CountingDTO.CountingResponseData{
			KodeAkun:    kode,
			NamaAkun:    namaAkun,
			TotalDebit:  totalDebit,
			TotalCredit: totalCredit,
			SaldoAkhir:  selisih,
		})
	}

	// Susun response akhir
	response := CountingDTO.CountingResponse{
		ResponseCode:    Constants.ResponseCodeGeneralSuccess,
		ResponseMessage: Constants.MsgSuccessRequest,
		Data:            responseData,
	}

	return &response, nil
}

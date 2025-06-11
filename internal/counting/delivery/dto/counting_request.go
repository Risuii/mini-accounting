package dto

import (
	"github.com/gin-gonic/gin"

	Library "mini-accounting/library"

	CustomValidation "mini-accounting/pkg/custom_validation"
)

type CountingRequest struct {
	Data     []CountingRequestData
	KodeAkun []CountingRequestKodeAkun
}

type CountingRequestData struct {
	Tanggal   string `csv:"Tanggal"`
	KodeAkun  string `csv:"Kode Akun"`
	NamaAkun  string `csv:"Nama Akun"`
	Deskripsi string `csv:"Deskripsi"`
	Debit     string `csv:"Debit"`
	Credit    string `csv:"Kredit"`
}

type CountingRequestKodeAkun struct {
	KodeAkun string `csv:"Kode Akun"`
}

type Summary struct {
	TotalDebit  float64
	TotalCredit float64
}

func (e *CountingRequest) Validate(v CustomValidation.CustomValidation, library Library.Library, c *gin.Context) []map[string]interface{} {
	journal, err := c.FormFile("journal")
	if err != nil {
		e = &CountingRequest{}
	}

	kodeAkun, err := c.FormFile("kodeAkun")
	if err != nil {
		e = &CountingRequest{}
	}

	j, err := journal.Open()
	if err != nil {
		e = &CountingRequest{}
	}

	k, err := kodeAkun.Open()
	if err != nil {
		e = &CountingRequest{}
	}

	defer j.Close()
	defer k.Close()

	var entriesJournal []CountingRequestData
	if err := library.GoCsvUnmarshal(j, &entriesJournal); err != nil {
		e = &CountingRequest{}
	}

	var entriesKodeAkun []CountingRequestKodeAkun
	if err := library.GoCsvUnmarshal(k, &entriesKodeAkun); err != nil {
		e = &CountingRequest{}
	}

	e.Data = entriesJournal
	e.KodeAkun = entriesKodeAkun

	var validationErrors []map[string]interface{}
	for _, entry := range entriesJournal {
		if errs := v.ValidateStruct(&entry, "name"); len(errs) > 0 {
			validationErrors = append(validationErrors, errs...)
		}
	}

	return validationErrors
}

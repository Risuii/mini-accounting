package http

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	Constants "mini-accounting/constants"
	Library "mini-accounting/library"

	CountingDTO "mini-accounting/internal/counting/delivery/dto"
	CountingUsecase "mini-accounting/internal/counting/domain/usecase"

	CustomErrorPackage "mini-accounting/pkg/custom_error"
	CustomValidationPackage "mini-accounting/pkg/custom_validation"
	LoggerPackage "mini-accounting/pkg/logger"
)

type CountingHandler interface {
	ReadCSV(c *gin.Context)
	DownloadCSV(c *gin.Context)
}

type CountingHandlerImpl struct {
	library          Library.Library
	countingUsecase  CountingUsecase.CountingUsecase
	customValidation CustomValidationPackage.CustomValidation
}

func NewCountingHandler(
	library Library.Library,
	countingUsecase CountingUsecase.CountingUsecase,
	customValidation CustomValidationPackage.CustomValidation,
) CountingHandler {
	return &CountingHandlerImpl{
		library:          library,
		countingUsecase:  countingUsecase,
		customValidation: customValidation,
	}
}

func (h *CountingHandlerImpl) ReadCSV(c *gin.Context) {
	path := "CountingHandler:ReadCSV"

	param := CountingDTO.CountingRequest{}
	errValidation := param.Validate(h.customValidation, h.library, c)
	if len(errValidation) > 0 {
		err := CustomErrorPackage.New(Constants.ErrValidation, Constants.ErrValidation, path, h.library)
		err = err.(*CustomErrorPackage.CustomError).FromListMap(errValidation)
		response := &CountingDTO.CountingResponse{
			ResponseCode:    Constants.ResponseCodeGeneralError,
			ResponseMessage: errValidation[0]["message"].(string),
			Data:            []CountingDTO.CountingResponseData{},
		}
		LoggerPackage.WriteLog(logrus.Fields{
			"request":  param,
			"path":     path,
			"response": nil,
		}).Debug(err.(*CustomErrorPackage.CustomError).GetPlain())

		c.JSON(http.StatusBadRequest, response)
		return
	}

	usecase := h.countingUsecase
	response, err := usecase.Index(&param)
	if err != nil {
		err = err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
		response := &CountingDTO.CountingResponse{
			ResponseCode:    Constants.ResponseCodeGeneralError,
			ResponseMessage: Constants.ErrInternalServerError.Error(),
			Data:            []CountingDTO.CountingResponseData{},
		}
		LoggerPackage.WriteLog(logrus.Fields{
			"request":  param,
			"path":     err.(*CustomErrorPackage.CustomError).GetPath(),
			"response": response,
		}).Debug(err.(*CustomErrorPackage.CustomError).GetPlain())

		c.JSON(http.StatusInternalServerError, response)
		return
	}

	rsp, _ := h.library.JsonMarshal(response)
	h.library.SetCache("csv", string(rsp), 300)

	// RESPONSE
	LoggerPackage.WriteLog(logrus.Fields{
		"request":  param,
		"path":     path,
		"response": response,
	}).Debug(nil)

	c.JSON(http.StatusOK, response)
}

func (h *CountingHandlerImpl) DownloadCSV(c *gin.Context) {

	var response CountingDTO.CountingResponse
	data, ok := h.library.GetCache("csv")
	h.library.JsonUnmarshal([]byte(data.(string)), &response)

	if ok {
		if str, ok := data.(string); ok {
			err := h.library.JsonUnmarshal([]byte(str), &response)
			if err != nil {
				// Tangani error jika unmarshal gagal
				log.Println("error unmarshal cache:", err)
			}
		} else {
			log.Println("cache data bukan string")
		}
	} else {
		log.Println("cache tidak ditemukan")
	}

	// ==== DOWNLOAD CSV ====
	filename := "counting_report.csv"
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", h.library.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")

	writer := h.library.CsvNewWriter(c.Writer)

	// Header kolom CSV
	writer.Write([]string{"Kode Akun", "Nama Akun", "Total Debit", "Total Kredit", "Saldo Akhir"})

	// Data isi CSV
	for _, row := range response.Data {
		writer.Write([]string{
			row.KodeAkun,
			row.NamaAkun,
			h.library.Sprintf("%.0f", row.TotalDebit),
			h.library.Sprintf("%.0f", row.TotalCredit),
			h.library.Sprintf("%.0f", row.SaldoAkhir),
		})
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		c.Status(http.StatusInternalServerError)
	}
}

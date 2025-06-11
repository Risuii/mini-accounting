package middlewares

// import (
// 	"net/http"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/sirupsen/logrus"

// 	Config "mini-accounting/config"
// 	Constants "mini-accounting/constants"
// 	Library "mini-accounting/library"

// 	LoggingEntity "mini-accounting/internal/logging/domain/entity"
// 	LogginRepository "mini-accounting/internal/logging/domain/repository"

// 	CustomErrorPackage "mini-accounting/pkg/custom_error"
// 	ExecutionResultPackage "mini-accounting/pkg/execution_result"
// 	LoggerPackage "mini-accounting/pkg/logger"
// 	LoggerHookPackage "mini-accounting/pkg/logger/hook"
// 	RequestPackage "mini-accounting/pkg/request_information"
// )

// type Middleware interface {
// 	CreateAPILog() gin.HandlerFunc
// }

// type MiddlewareImpl struct {
// 	loggingRepository LogginRepository.LoggingRepository
// 	config            Config.Config
// 	library           Library.Library
// }

// func NewMiddleware(
// 	loggingRepository LogginRepository.LoggingRepository,
// 	config Config.Config,
// 	library Library.Library,
// ) Middleware {
// 	return &MiddlewareImpl{
// 		loggingRepository: loggingRepository,
// 		config:            config,
// 		library:           library,
// 	}
// }

// func (m *MiddlewareImpl) CreateAPILog() gin.HandlerFunc {
// 	functionPath := "Middleware:CreateAPILog"
// 	return func(c *gin.Context) {
// 		// SET DEFAULT VALUE
// 		traceId, _ := m.library.GenerateUUID()
// 		c.Set("TraceId", traceId)
// 		start := m.library.GetNow()
// 		buffer := m.library.NewBufferString("")
// 		hook := LoggerHookPackage.New(buffer, m.library)
// 		LoggerPackage.GetLogger().AddHook(hook)

// 		c.Next()

// 		// AVOID LOGING IF STATUS IS UNAUTHORIZED
// 		if c.Writer.Status() == http.StatusUnauthorized {
// 			return
// 		}

// 		// AVOID LOGING IF ABORTED BY ANOTHER MIDDLEWARE
// 		_, exists := c.Get("isTimeout")
// 		if !exists && c.IsAborted() {
// 			return
// 		}

// 		finish := m.library.GetNow()
// 		var response *gin.H
// 		var responseJSON map[string]interface{}
// 		responseJSONString := hook.GetBuffer().String()
// 		m.library.JsonUnmarshal([]byte(responseJSONString), &responseJSON)
// 		// GET REQUEST BY GOROUTINE CHANNEL
// 		requestChannel := make(chan ExecutionResultPackage.ExecutionResult)
// 		go func(resultChannel chan ExecutionResultPackage.ExecutionResult) {
// 			result := ExecutionResultPackage.ExecutionResult{}
// 			switch responseJSON["request"].(type) {
// 			case string:
// 				requestString := responseJSON["request"].(string)
// 				result.SetResult(requestString, nil)
// 				resultChannel <- result
// 				close(resultChannel)
// 				return

// 			}
// 			result.SetResult(nil, nil)
// 			resultChannel <- result
// 			close(resultChannel)
// 		}(requestChannel)
// 		// GET RESPONSE BY GOROUTINE CHANNEL
// 		responseChannel := make(chan ExecutionResultPackage.ExecutionResult)
// 		go func(resultChannel chan ExecutionResultPackage.ExecutionResult) {
// 			result := ExecutionResultPackage.ExecutionResult{}
// 			switch responseJSON["response"].(type) {
// 			case map[string]interface{}:
// 			case string:
// 				result.SetResult(responseJSON["response"].(string), nil)
// 				resultChannel <- result
// 				close(resultChannel)
// 				return
// 			}
// 			result.SetResult(nil, nil)
// 			resultChannel <- result
// 			close(resultChannel)
// 		}(responseChannel)
// 		// SEND THE RESULT OF GET REQUEST BY GOROUTINE CHANNEL
// 		requestResult := <-requestChannel
// 		// SEND THE RESULT OF GET RESPONSE BY GOROUTINE CHANNEL
// 		responseResult := <-responseChannel
// 		// WHEN GET REQUEST BY GOROUTINE CHANNEL RETURNS ERROR
// 		if requestResult.GetError() != nil {
// 			err := CustomErrorPackage.New(Constants.ErrInternalServerError, requestResult.GetError(), functionPath, m.library)
// 			response = &gin.H{
// 				"response": gin.H{
// 					"responseCode":    Constants.ResponseCodeGeneralError,
// 					"responseMessage": err.(*CustomErrorPackage.CustomError).GetDisplay(),
// 				},
// 			}
// 			LoggerPackage.WriteLog(logrus.Fields{
// 				"path":     functionPath,
// 				"request":  requestResult.GetData().([]byte),
// 				"response": responseResult.GetData().([]byte),
// 			}).Debug(err.(*CustomErrorPackage.CustomError).GetPlain())

// 			c.AbortWithStatusJSON(http.StatusInternalServerError, response)
// 			return
// 		}
// 		// WHEN GET RESPONSE BY GOROUTINE CHANNEL RETURNS ERROR
// 		if responseResult.GetError() != nil {
// 			err := CustomErrorPackage.New(Constants.ErrInternalServerError, responseResult.GetError(), functionPath, m.library)
// 			response = &gin.H{
// 				"response": gin.H{
// 					"responseCode":    Constants.ResponseCodeGeneralError,
// 					"responseMessage": err.(*CustomErrorPackage.CustomError).GetDisplay(),
// 				},
// 			}
// 			LoggerPackage.WriteLog(logrus.Fields{
// 				"path":     functionPath,
// 				"request":  requestResult.GetData().([]byte),
// 				"response": responseResult.GetData().([]byte),
// 			}).Debug(err.(*CustomErrorPackage.CustomError).GetPlain())

// 			c.AbortWithStatusJSON(http.StatusInternalServerError, response)
// 			return
// 		}

// 		loggingLabel := ""
// 		clientName := "Default"
// 		if _, ok := c.Get("ClientName"); ok {
// 			clientName = c.MustGet("ClientName").(string)
// 		}
// 		apiLogChannel := make(chan ExecutionResultPackage.ExecutionResult)
// 		go func(resultChannel chan ExecutionResultPackage.ExecutionResult) {

// 			result := ExecutionResultPackage.ExecutionResult{}
// 			logEntity := LoggingEntity.InterfaceLog{
// 				TraceID:         c.MustGet("TraceId").(string),
// 				ServiceName:     c.Request.URL.Path,
// 				ClientName:      clientName,
// 				RequestPayload:  requestResult.GetData().(string),
// 				ResponsePayload: responseResult.GetData().(string),
// 				RequestDate:     start,
// 				ResponseDate:    m.library.GetNow(),
// 			}

// 			err := m.loggingRepository.InsertInterfaceLog(logEntity)
// 			if err != nil {
// 				result.SetResult(nil, err.(*CustomErrorPackage.CustomError).UnshiftPath(functionPath))
// 				resultChannel <- result
// 				close(resultChannel)
// 				return
// 			}
// 			result.SetResult(nil, nil)
// 			resultChannel <- result
// 			close(resultChannel)
// 		}(apiLogChannel)

// 		// GET "executionLabel" THAT IS PASSED FROM "TimeoutLimiter"
// 		executionLabel := ""
// 		value, exist := c.Get("executionLabel")
// 		if exist {
// 			executionLabel = value.(string)
// 		}
// 		// TIMEOUT MANAGEMENT BY GOROUTINE
// 		select {
// 		case apiLogResult := <-apiLogChannel:
// 			if apiLogResult.GetError() != nil {
// 				requestInformation := RequestPackage.RequestInformation{}
// 				err := apiLogResult.GetError()
// 				response = &gin.H{
// 					"response": gin.H{
// 						"responseCode":    Constants.ResponseCodeGeneralError,
// 						"responseMessage": err.(*CustomErrorPackage.CustomError).GetDisplay(),
// 					},
// 				}
// 				LoggerPackage.WriteLog(logrus.Fields{
// 					"path":     "Middleware:CreateAPILog",
// 					"request":  requestInformation.GetRequestBodyMap(),
// 					"response": response,
// 				}).Debug(err.(*CustomErrorPackage.CustomError).GetPlain())
// 			}
// 			break
// 		case <-m.library.TimeAfter(m.config.GetConfig().App.LoggingTimeout):
// 			loggingLabel = "Timeout!"
// 			break
// 		}
// 		final := m.library.GetNow()
// 		m.DisplayInfo(&start, &finish, &final, &executionLabel, &loggingLabel, &traceId)
// 	}
// }

// func (m *MiddlewareImpl) DisplayInfo(start, finish, final *time.Time, executionLabel, loggingLabel, traceID *string) {
// 	logging := "Logging time\t\t: -\n"
// 	loggingTime := "Logged\t\t\t: -\n"
// 	totalDuration := (*finish).Sub(*start).Abs().Milliseconds()
// 	if final != nil {
// 		logging = m.library.Sprintf("Logging time\t\t: %dms %s\n", (*final).Sub(*finish).Abs().Milliseconds(), *loggingLabel)
// 		loggingTime = m.library.Sprintf("Logged\t\t\t: %s\n", (*final).Format("2006-01-02 15:04:05 Z0700 MST"))
// 		totalDuration = (*final).Sub(*start).Abs().Milliseconds()
// 	}
// 	m.library.Println("------------------------------------------------------------")
// 	m.library.Printf("Start\t\t\t: %s\n", (*start).Format("2006-01-02 15:04:05 Z0700 MST"))
// 	m.library.Printf("Finish\t\t\t: %s\n", (*finish).Format("2006-01-02 15:04:05 Z0700 MST"))
// 	m.library.Printf(loggingTime)
// 	m.library.Printf("Execution time\t\t: %dms %s\n", (*finish).Sub(*start).Abs().Milliseconds(), *executionLabel)
// 	m.library.Printf(logging)
// 	m.library.Printf("Total response time\t: %dms\n", totalDuration)
// 	m.library.Printf("Trace ID: %s", *traceID)
// }

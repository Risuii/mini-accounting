package utils

import (
	"github.com/sirupsen/logrus"

	Constants "mini-accounting/constants"
	Library "mini-accounting/library"
	CustomErrorPackage "mini-accounting/pkg/custom_error"
	LoggerPackage "mini-accounting/pkg/logger"
)

func TernaryOperator[T interface{}](comparator bool, trueCondition T, falseCondition T) T {
	if comparator {
		return trueCondition
	}

	return falseCondition
}

func TernaryOperatorPromise[T interface{}](comparator bool, trueCallback func() T, falseCallback func() T) T {
	if comparator {
		return trueCallback()
	}

	return falseCallback()
}

func CatchPanic(path string, library Library.Library) {
	if err := recover(); err != nil {
		err := CustomErrorPackage.New(Constants.ErrPanic, err.(error), path, library)
		LoggerPackage.WriteLog(logrus.Fields{
			"path":  err.(*CustomErrorPackage.CustomError).GetPath(),
			"title": err.(*CustomErrorPackage.CustomError).GetDisplay().Error(),
		}).Panic(err.(*CustomErrorPackage.CustomError).GetPlain())
	}
}

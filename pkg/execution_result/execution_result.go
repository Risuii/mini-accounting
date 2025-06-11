package execution_result

type ExecutionResult struct {
	data interface{}
	err  error
}

func (e *ExecutionResult) SetResult(data interface{}, err error) {
	e.data = data
	e.err = err
}

func (e *ExecutionResult) GetError() error {
	return e.err
}

func (e *ExecutionResult) GetData() interface{} {
	return e.data
}

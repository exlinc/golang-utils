package sqlkit

type Response struct {
	Success      bool
	ErrorCode    string
	ErrorMessage string
}

func (resp *Response) CheckError(log FieldLogger, err error, code string) {
	if err != nil {
		resp.Success = false
		resp.ErrorCode = code
		resp.ErrorMessage = err.Error()
		if log != nil {
			log.Errorf("SQL_ERROR: %s", err.Error())
		}
	} else {
		resp.Success = true
	}
}

func (resp *Response) Error() string {
	return resp.ErrorMessage
}

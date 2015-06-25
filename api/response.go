package api

const TYPE_RESPONSE = "Response"

type Response struct {
    result bool
    errors []error
}

/* Constructs a new Response message */
func NewResponse(errors []error) *Response {
    instance := new(Response)
    instance.result = (len(errors) == 0)
    instance.errors = errors
    return instance
}

/* Retrieves the result */
func (this *Response) Result() bool {
    return this.result
}

/* Retrieves the errors */
func (this *Response) Errors() []error {
    return this.errors
}

/* Validates the message */
func (this *Response) IsValid() (bool, error) {
    return true, nil
}

/* Retrieves the message type */
func (this *Response) Type() string {
    return TYPE_RESPONSE
}

/* Converts the message to map */
func (this *Response) Mapify() interface{} {
    resultString := ""

    if this.result {
        resultString = "OK"
    } else {
        resultString = "Error"
    }

    errorsString := make([]string, len(this.errors))
    index := 0

    for _, err := range this.errors {
        errorsString[index] = err.Error()
        index += 1
    }

    return map[string]map[string]interface{} {
        TYPE_RESPONSE: {
            "Result": resultString,
            "Errors": errorsString,
        },
    }
}

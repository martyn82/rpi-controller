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

/* Retrieves the message type */
func (this *Response) Type() string {
    return TYPE_RESPONSE
}

/* Converts the message to JSON */
func (this *Response) JSON() string {
    resultString := ""

    if this.result {
        resultString = "OK"
    } else {
        resultString = "Error"
    }

    errorsString := ""

    for _, err := range this.errors {
        if errorsString != "" {
            errorsString += ","
        }

        errorsString += "\"" + err.Error() + "\""
    }

    return "{\"" + TYPE_RESPONSE + "\":{\"Result\":\"" + resultString + "\",\"Errors\":[" + errorsString + "]}}"
}

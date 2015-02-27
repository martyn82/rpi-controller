package device

/* DenonAvr type */
type DenonAvr struct {
    DeviceModel
}

/* Map of messages */
var denonAvrMessageMap = map[string]string{
    "PW:ON": "PWON\r",
    "PW:OFF": "PWSTANDBY\r",
}

/* Construct DenonAvr */
func CreateDenonAvr(name string, model string, protocol string, address string) *DenonAvr {
    d := new(DenonAvr)
    d.name = name
    d.model = model
    d.protocol = protocol
    d.address = address
    d.messageMap = denonAvrMessageMap
    return d
}

package device

import (
    "encoding/base64"
    "errors"
    "net"
    "os"
    "strconv"
    "strings"

    "github.com/martyn82/rpi-controller/messages"
)

/* SamsungTv type */
type SamsungTv struct {
    DeviceModel
    tvAppName string
    isAuthenticated bool
}

var samsungTvPropertyMap = map[string]string{
    messages.PROP_POWER: "KEY_POWER",
    messages.PROP_VOLUME: "KEY_VOL",
}

var samsungTvValueMap = map[string]string{
    messages.VAL_ON: "ON",
    messages.VAL_OFF: "OFF",
}

/* Constructs SamsungTv */
func CreateSamsungTv(name string, model string, protocol string, address string) *SamsungTv {
    d := new(SamsungTv)
    d.name = name
    d.model = model
    d.protocol = protocol
    d.address = address
    d.isAuthenticated = false

    d.mapMessage = func (message string) string {
        msg, err := messages.ParseMessage(message)

        if err != nil {
            return message
        }

        value := samsungTvValueMap[msg.Value]

        if v, err := strconv.Atoi(msg.Value); err == nil {
            if v > 0 {
                value = "UP"
            } else {
                value = "DOWN"
            }
        }

        return samsungTvPropertyMap[msg.Property] + value
    }

    d.processResponse = func (response []byte) string {
        return strconv.Quote(string(response))
    }

    return d
}

/* Connects to device */
func (d *SamsungTv) Connect() error {
    connectErr := d.DeviceModel.Connect()

    if connectErr != nil {
        return connectErr
    }

    return d.authenticate()
}

/* Retrieves network information */
func getNetworkInfo() (hostName string, ipAddress string, macAddress string, err error) {
    fqn, _ := os.Hostname()
    hostName = strings.Split(fqn, ".")[0]
    ifaces, err := net.Interfaces()

    if err != nil {
        return "", "", "", err
    }

    for _, iface := range ifaces {
        addrs, err := iface.Addrs()

        if err != nil {
            continue
        }

        for _, addr := range addrs {
            if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
                if ipnet.IP.To4() != nil {
                    return hostName, ipnet.IP.To4().String(), iface.HardwareAddr.String(), nil
                }
            }
        }
    }

    return "", "", "", errors.New("Unable to gather network information.")
}

/* Authenticates with the device */
func (d *SamsungTv) authenticate() error {
    if d.isAuthenticated {
        return nil
    }
    
    hostName, remoteIp, mac, err := getNetworkInfo()

    if err != nil {
        return err
    }

    remoteMac := strings.Replace(mac, ":", "-", -1)
    remoteName := hostName + "..iapp.samsung"
    tvAppName := remoteName + "." + d.model + ".iapp.samsung"
    d.tvAppName = tvAppName

    remoteIpEnc := base64.StdEncoding.EncodeToString([]byte(remoteIp))
    remoteMacEnc := base64.StdEncoding.EncodeToString([]byte(remoteMac))
    remoteNameEnc := base64.StdEncoding.EncodeToString([]byte(remoteName))

    tvAppNameLen, _ := strconv.Unquote(strconv.QuoteRuneToASCII(rune(len(tvAppName))))
    remoteIpLen, _ := strconv.Unquote(strconv.QuoteRuneToASCII(rune(len(remoteIpEnc))))
    remoteMacLen, _ := strconv.Unquote(strconv.QuoteRuneToASCII(rune(len(remoteMacEnc))))
    remoteNameLen, _ := strconv.Unquote(strconv.QuoteRuneToASCII(rune(len(remoteNameEnc))))

    authPayload := "\x64\x00" +
        remoteIpLen + "\x00" + remoteIpEnc +
        remoteMacLen + "\x00" + remoteMacEnc +
        remoteNameLen + "\x00" + remoteNameEnc
    authPayloadLen, _ := strconv.Unquote(strconv.QuoteRuneToASCII(rune(len(authPayload))))
    authMsg := "\x00" + tvAppNameLen + "\x00" + tvAppName +
        authPayloadLen + "\x00" + authPayload

    if writeErr := d.DeviceModel.SendMessage(authMsg); writeErr != nil {
        return writeErr
    }

    secondPayload := "\xC8\x00"
    secondPayloadLen, _ := strconv.Unquote(strconv.QuoteRuneToASCII(rune(len(secondPayload))))
    secondMsg := "\x00" + tvAppNameLen + "\x00" + tvAppName +
        secondPayloadLen + "\x00" + secondPayload

    if writeErr := d.DeviceModel.SendMessage(secondMsg); writeErr != nil {
        return writeErr
    }

    d.isAuthenticated = true
    return nil
}

/* Sends message to device */
func (d *SamsungTv) SendMessage(message string) error {
    if !d.isAuthenticated {
        d.authenticate()
    }

    message = d.mapMessage(message)

    tvAppName := d.tvAppName
    tvAppNameLen, _ := strconv.Unquote(strconv.QuoteRuneToASCII(rune(len(tvAppName))))

    keyEnc := base64.StdEncoding.EncodeToString([]byte(message))
    keyLen, _ := strconv.Unquote(strconv.QuoteRuneToASCII(rune(len(keyEnc))))
    keyPayload := "\x00\x00\x00" + keyLen + "\x00" + keyEnc
    keyPayloadLen, _ := strconv.Unquote(strconv.QuoteRuneToASCII(rune(len(keyPayload))))
    keyMsg := "\x00" + tvAppNameLen + "\x00" + tvAppName +
        keyPayloadLen + "\x00" + keyPayload

    if writeErr := d.DeviceModel.SendMessage(keyMsg); writeErr != nil {
        return writeErr
    }

    return nil
}

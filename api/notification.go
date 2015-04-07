package api

import "errors"

const (
    TYPE_NOTIFICATION = "Event"

    ERR_INVALID_NOTIFICATION = "Invalid event notification; missing device and/or property name."
)

type INotification interface {
    IMessage
    PropertyName() string
    PropertyValue() string
}

type Notification struct {
    deviceName string
    propertyName string
    propertyValue string
}

/* Create a Notification from map */
func notificationFromMap(message map[string]string) (*Notification, error) {
    var deviceName string
    var propertyName string
    var propertyValue string

    for k, v := range message {
        if k == KEY_DEVICE {
            deviceName = v
        } else {
            propertyName = k
            propertyValue = v
        }
    }

    result := NewNotification(deviceName, propertyName, propertyValue)

    if _, err := result.IsValid(); err != nil {
        return nil, err
    }

    return result, nil
}

/* Create a Notification command */
func NewNotification(deviceName string, propertyName string, propertyValue string) *Notification {
    instance := new(Notification)
    instance.deviceName = deviceName
    instance.propertyName = propertyName
    instance.propertyValue = propertyValue
    return instance
}

/* Retrieve the device name */
func (this *Notification) DeviceName() string {
    return this.deviceName
}

/* Retrieve the property name */
func (this *Notification) PropertyName() string {
    return this.propertyName
}

/* Retrieve the property value */
func (this *Notification) PropertyValue() string {
    return this.propertyValue
}

/* Validates the notification */
func (this *Notification) IsValid() (bool, error) {
    if this.deviceName == "" || this.propertyName == "" {
        return false, errors.New(ERR_INVALID_NOTIFICATION)
    }

    return true, nil
}

/* Retrieves the type of the message */
func (this *Notification) Type() string {
    return TYPE_NOTIFICATION
}

/* Convert notification to string */
func (this *Notification) JSON() string {
    return "{\"" + TYPE_NOTIFICATION + "\":{\"" + KEY_DEVICE + "\":\"" + this.deviceName + "\",\"" + this.propertyName + "\":\"" + this.propertyValue + "\"}}"
}

package command

type Notification struct {
    deviceName string
    propertyName string
    propertyValue string
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

/* Convert notification to string */
func (this *Notification) String() string {
    return this.deviceName + ":" + this.propertyName + ":" + this.propertyValue
}

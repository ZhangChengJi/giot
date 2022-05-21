package event

type event struct {
	methodType string //add、update、delete
	sponsor    string //product、device
	productId  string //
	deviceId   string
	action     string //product、device 、slave、model、alarm
}

package pfmobile

// Special program to send a lots of sms using a mobile phone!
// uses logfilename for logging
// uses phonenumbersfilename to specify file with phonenumbers
// 2024-01-21 working!!!!
// 2024-03-10 switched to newer serial driver, implemented support for S24U and model selection
// got it working with Samsung S24Ultra! speed 14s/sms using timeout = Millisecond*700

// Special program to send a lots of sms using a mobile phone!
// uses logfilename for logging
// uses phonenumbersfilename to specify file with phonenumbers
// 2024-01-21 working!!!!
// 2024-03-10 switched to newer serial driver, implemented support for S24U and model selection
// got it working with Samsung S24Ultra! speed 14s/sms using timeout = Millisecond*700

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
	"unicode/utf16"
	"unicode/utf8"

	"fyne.io/fyne/v2"
	"go.bug.st/serial"
)

func Modemcommand(port serial.Port, modemcommand string, expectedresponse string, timeout time.Duration, description string, err error) error {
	// (port,modemcommand,expected response,possible timeout,description,previous result err)
	var r string
	if err != nil {
		return err
	}
	err =port.SetReadTimeout(timeout)
	if err != nil {
		log.Printf("Modemcommand port.SetReadTimeout(%d) failed: %s\r\n",timeout, err)
		return err
	}
	err = mywrite(port,modemcommand)
	if err!=nil {
		return err
	}
	r,err = myread(port,expectedresponse)
	if r>"" {
		fmt.Printf("Modemcommand '%s' ? %s (%s) got: '%s'\r\n",showdebugmsg(modemcommand),expectedresponse,description, showdebugmsg(r))
	}
	if err != nil {
		return errors.New("Modemcommand "+description+" error: '"+showdebugmsg(r)+err.Error()+"'")
	}		
	if !strings.Contains(r,expectedresponse) {
		return errors.New("Modemcommand "+description+" unexpected response: '"+showdebugmsg(r)+"' expected: '"+expectedresponse+"'")
	}
	return err
}
func Openmodemport(comport string) (serial.Port, error) {
	var port serial.Port
	var err error
	if comport=="" {
		comport= fyne.CurrentApp().Preferences().StringWithFallback("mobileport", "COM3")
	}
	// err=TestPort(comport)
	// if  err!= nil {
	// 	return nil,errors.New("Openmodemport #1 failed: " + err.Error())
	// }
	mode := &serial.Mode{
		BaudRate: 115200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	port, err = serial.Open(comport, mode)
	if err != nil {
		return nil,errors.New("Openmodemport #2 failed: " + err.Error())
	}
	port.Break(time.Second)
	port.Drain()
	port.ResetInputBuffer()
	return port, nil
}
func Modemreset(comport string) error {
	// modemcommand format:
	// (port,modemcommand,expected response,possible timeout,description,previous result err)
	var port serial.Port
	if comport=="" {
		comport= fyne.CurrentApp().Preferences().StringWithFallback("mobileport", "COM3")
	}
	var err error
	err = TestPort(comport)
	if  err != nil {
		return errors.New("Modemreset: Testcomport #1 failed: " + err.Error())
	}
	port,err=Openmodemport("")
	if err!=nil {
		return errors.New("Modemreset Openmodemport failed: " + err.Error())
	}
	err=Modemcommand(port,"\032\n","",time.Second*2,"wakeup",err)
	err=nil
	err=Modemcommand(port,"\032\n","",time.Second*2,"wakeup",err)
	err=Modemcommand(port,"AT+CREG?\n","",time.Second*3,"check registration",err)
	err=Modemcommand(port,"AT+CGDCONT?\n","",time.Second*3,"check CGDCONT",err)
	err=Modemcommand(port,"AT+CGATT?\n","",time.Second*3,"check GATT",err)
	err=Modemcommand(port,"AT+CFUN=0\n","",time.Second*3,"set full func off",err)
	err=Modemcommand(port,"ATZ\n","",time.Second,"reset modem",err)
	err=Modemcommand(port,"ATE1\r\n","",time.Second,"set echo on",err)
	err=Modemcommand(port,"ATV1\r\n","",time.Second,"set verbose on",err)
	err=Modemcommand(port,"AT\r\n","OK",time.Second,"test AT command",err)
	err=Modemcommand(port,"AT+CFUN=1\r\n","",time.Second*2,"set full func on",err)
	err=Modemcommand(port,"AT\r\n","OK",time.Second,"test AT command",err)
	// Set SMS center number = Universal 	= 0046735480000
	// Set SMS center number = Telenor 		= 0046708000708
	// err=modemcommand(port,"AT+CSCA=\"0046708000708\"\r\n","OK",time.Second,"set SMS center number",err)
	// err=modemcommand(port,"AT\r\n","OK",time.Second,"test AT command",err)
	// err=modemcommand(port,"AT+CMGF=0\r\n","OK",time.Second,"set PDU mode",err)
	// err=modemcommand(port,string("AT+DEVCONINFO\r\n"),"OK",time.Second,"get device info",err)
	// err=modemcommand(port, "AT+CGMI\n","OK",time.Second,"get manufacturer",err)		//GET MANUFACTURER
	// err=modemcommand(port, "AT+CGMM\n","OK",time.Second,"get model",err)		//GET MODEL
	port.Close()
	if err!=nil {
		return errors.New("Modemreset failed: " + err.Error())
	}
	// r should be "AT+CFUN=0,0\r\n\r\nOK\r\nAT+CMGF=0\r\nATE1\r\n\r\nOK\r\nAT+CFUN=1,0\r\n\r\nOK\r\n"?!
	return nil
}
func SendSMS(port serial.Port, phoneNumber string, message string) error {
	var pduarray []string
	var cmd1 []string
	var cmd2 []string
	var err error
	if port == nil {
		port, err = Openmodemport("")
		if err != nil {
			return errors.New("SendSMS #1 openmodemport failed: " + err.Error())
		}
	}
	pduarray = CreateLongPDU(phoneNumber, message)
	for i := 0; i < len(pduarray); i++ {
		cmd1 = append(cmd1, "AT+CMGS="+fmt.Sprintf("%d", (len(pduarray[i])-2)/2)+"\r\n")
		cmd2 = append(cmd2, pduarray[i]+string(rune(26)))
	}
	for i := 0; i < len(cmd1); i++ {
		fmt.Println("mywrite:", showdebugmsg(cmd1[i]))
		err =Modemcommand(port,cmd1[i],">",time.Second * 2,"cmd1[i] ? >",nil)
		if err !=nil {
			port.Close()
			return errors.New("SendSMS #2 myread failed: " + err.Error())
		}
		err =Modemcommand(port,cmd2[i],">",time.Second * 3,"cmd2[i] ? OK",nil)
		fmt.Println("mywrite:", showdebugmsg(cmd2[i]))
		if err !=nil {
			port.Close()
			return errors.New("SendSMS #3 myread failed: " + err.Error())
		}
	}
	return nil
}
func WriteWithTimeout(port serial.Port, s string, timeout time.Duration) (int, error) {
	type result struct {
		n   int
		err error
	}
	done := make(chan result, 1)

	go func() {
		n, err := port.Write([]byte(s))
		done <- result{n, err}
	}()

	select {
	case res := <-done:
		return res.n, res.err
	case <-time.After(timeout):
		return 0, fmt.Errorf("write timeout")
	}
}
func mywrite(port serial.Port,s string ) error {
	n,err:=port.Write([]byte(s))
	if err!=nil {
		return err
	}
	if n== 0 {
		return errors.New("no bytes written to port")
	}
	err = port.Drain()
	if err!=nil {
		return err
	}
	return nil
}
func myread(port serial.Port,response string) (string,error) {
	var r string
	var err error
	var n int
	timeout:=time.Second * 20
	buff := make([]byte, 100)
	startTime := time.Now()
	for {
		n, err = port.Read(buff)
		if err != nil {
			break 
		}
		if n > 0 {
			r= fmt.Sprintf("%v", string(buff[:n]))
		} 
		if  r>"" {
			if strings.Contains(r,response) || response==""{
				break
			}
		}
		if response=="" && n==0 {
			break
		}
		if time.Since(startTime) > timeout {
            err= errors.New("2 second read timeout")
			break
        }
	}
	return r,err
}
func CreateLongPDU(phoneNumber string, message string) []string {
	const maxCharsPerSegment = 67 // Maximum characters per segment
	var segments []string
	var pdus []string
	tmsg := message
	for tmsg > "" {
		if utf8.RuneCountInString(tmsg) < maxCharsPerSegment {
			segments = append(segments, tmsg)
			tmsg = ""
		} else {
			segments = append(segments, string([]rune(tmsg)[:maxCharsPerSegment])) // Extract segment of the message
			tmsg = string([]rune(tmsg)[maxCharsPerSegment:])
		}
	}
	for i := 0; i < len(segments); i++ {
		var pdu string
		// UDH format:
		// UDH Length (1 byte) | Information Element Identifier (1 byte) | Information Element Data Length (1 byte) |
		// 0x00 (1 byte) | Message Reference (1 byte) | Total Parts (1 byte) | Sequence Number (1 byte)
		udh := fmt.Sprintf("05000300%02X%02X", len(segments), i+1) // Construct UDH for segmented message with total length
		pdu = CreatePDU(phoneNumber, segments[i], udh)             // Construct PDU for the segment with UDH
		pdus = append(pdus, pdu)
	}
	return pdus
}
func CreatePDU(number string, message string, udh string) string {
	// Ensure the phone number is in the correct format (e.g., with TOA)
	var pdu, pduHeader, pduMessage, pduMessageLen string
	phoneNumber := strings.TrimPrefix(number, "+")

	// Add a padding nibble if the phone number length is odd
	if len(phoneNumber)%2 != 0 {
		phoneNumber += "F" // Padding nibble
	}

	// Convert the phone number to semi-octets
	semiOctets := ""
	for i := 0; i < len(phoneNumber); i += 2 {
		semiOctets = semiOctets + string(phoneNumber[i+1]) + string(phoneNumber[i])
	}

	// Calculate the PDU length based on the semi-octets of the phone number
	pduLength := len(semiOctets) - 1

	// Convert the message to UCS-2 encoding (16-bit Unicode transformation format)
	encodedMessage := utf16.Encode([]rune(message))
	// UCS-2 encoded message as hex
	ucs2EncodedMessage := ""
	for _, char := range encodedMessage {
		ucs2EncodedMessage += fmt.Sprintf("%04X", char)
	}
	// Assemble the PDU message string
	pduMessage = ucs2EncodedMessage
	//	pduMessageLen = string(pduMessageLen[2:4]) + string(pduMessageLen[0:2])
	if udh == "" {
		// Assemble the PDU header string
		pduHeader = fmt.Sprintf("001100%02x91%s00080B", pduLength, semiOctets)
		pduMessageLen = fmt.Sprintf("%02X", len(pduMessage)/2)
		pdu = pduHeader + pduMessageLen + pduMessage
	} else {
		// Assemble the PDU header string
		pduHeader = fmt.Sprintf("005100%02x81%s00080B", pduLength, semiOctets)
		pduMessageLen = fmt.Sprintf("%02X", len(pduMessage)/2+6)
		pdu = pduHeader + pduMessageLen + udh + pduMessage
	}
	pdu = strings.Replace(pdu, "a", "A", -1)
	pdu = strings.Replace(pdu, "b", "B", -1)
	pdu = strings.Replace(pdu, "c", "C", -1)
	pdu = strings.Replace(pdu, "d", "D", -1)
	pdu = strings.Replace(pdu, "e", "E", -1)
	pdu = strings.Replace(pdu, "f", "F", -1)
	pdu = strings.TrimSpace(pdu)
	return pdu
}
func GetPortsList() ([]string, error) {
	return serial.GetPortsList()
}
func TestPort(comport string) error {
	var port serial.Port
	var err error
	var n int
	mode := &serial.Mode{
		BaudRate: 115200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	port, err = serial.Open(comport, mode)
	port.Drain()
	port.ResetInputBuffer()
	if err != nil {
		return err // err.Error("#1 serial.Open(comport)")
	}
	// err= port.SetReadTimeout(time.Second * 2)
	// if err != nil {
	// 	return err //err.Error("#2 port.SetReadTimeout(time.Second * 2)")
	// }
	// err = port.Break(time.Second)
	// if err != nil {
	// 	return err //("#3 port.Break(time.Second)", err)
	// }
	// check port write possibility
	n,err =WriteWithTimeout(port,"AT\r\n",time.Second*2)
	if err!=nil {
		port.Close()
		return err
	}
	if n==0 {
		return errors.New("no bytes written to port")
	}
	port.SetReadTimeout(time.Second)
	buff := make([]byte, 100)
	n,err=port.Read(buff)
	fmt.Print("...Got ",n," characters, error: ", err)
	if n>0 {
		fmt.Printf("...Data: '%s'", showdebugmsg(string(buff[:n])))
	}
	port.Close()
	if !strings.Contains(string(buff[:n]), "OK") {
		return errors.New("no OK received from port")
	}
	return nil
}
func GetMobilePort() (string, error) {
	ports, err := GetPortsList()
	if err != nil {
		return "", err
	}
	for i := 0; i < len(ports); i++ {
		err = TestPort(ports[i])
		if err == nil {
			return ports[i], nil
		}
	}
	return "", errors.New("no valid mobile port found")
}
func showdebugmsg(s string) string {
	r2 := s
	r2 = strings.Replace(r2, string(rune(13)), "\\r", -1)
	r2 = strings.Replace(r2, string(rune(10)), "\\n", -1)
	r2 = strings.Replace(r2, string(rune(0)), "\\0", -1)
	r2 = strings.Replace(r2, string(rune(9)), "\\t", -1)
	r2 = strings.Replace(r2, string(rune(26)), "\\z", -1)
	return r2
}

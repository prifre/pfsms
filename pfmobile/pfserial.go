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
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"
	"unicode/utf8"

	"go.bug.st/serial"
)

func Modemcommand(port serial.Port, modemcommand string, expectedresponse string, timeout int, description string, err error) error {
	// (port,modemcommand,expected response,possible timeout,description,previous result err)
	var r string
	if err != nil {
		return err
	}
	err =port.SetReadTimeout(time.Duration(timeout/100))
	if err != nil {
		log.Printf("Modemcommand Error #1 port.SetReadTimeout(%d) failed: %s\r\n",timeout, err)
		return err
	}
	err = mywrite(port,modemcommand)
	if err!=nil {
		return err
	}
	r,err = myread(port,expectedresponse,timeout)
	if err != nil {
		return errors.New("Modemcommand Error #2 "+description+" readerror '"+showdebugmsg(r)+err.Error()+"'")
	}		
	if !strings.Contains(r,expectedresponse) {
		return errors.New("Modemcommand Error #3 "+description+" unexpected response: '"+showdebugmsg(r)+"' expected: '"+expectedresponse+"'")
	}
	// if r>"" {
	// 	log.Printf("Modemcommand '%s' ? %s (%s) got: '%s'\r\n",showdebugmsg(modemcommand),expectedresponse,description, showdebugmsg(r))
	// }
	return err
}
func Openmodemport(comport string) (serial.Port, error) {
	var port serial.Port
	var err error
	if comport=="" {
		comport,err=GetMobilePort()
		if err!=nil {
			return nil,errors.New("Openmodemport #0 failed: " + err.Error())
		}
	}
	// err=TestPort(comport)
	// if  err!= nil {
	// 	return nil,errors.New("Openmodemport #1 failed: " + err.Error())
	// }
	mode := &serial.Mode{
		BaudRate: 115200,
		DataBits: 8,
		Parity:   serial.NoParity,
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
func Modemreset(comport string) (serial.Port,error) {
	// modemcommand format:
	// (port,modemcommand,expected response,possible timeout,description,previous result err)
	var port serial.Port
	var err error
	if comport=="" {
		comport,err=GetMobilePort()
		if err!=nil {
			return nil,errors.New("Modemreset #1 GetMobilePort failed: " + err.Error())
		}
	}
	port,err=Openmodemport(comport)
	if err!=nil {
		return nil,errors.New("Modemreset #2 Openmodemport failed: " + err.Error())
	}
	err=Modemcommand(port,"\032\r","",1,"ESC",err)
	err=Modemcommand(port,"AT\r","OK",1,"AT",err)
	err=Modemcommand(port,"AT+CFUN=1\r","OK",10,"reset modem",nil)
	err=Modemcommand(port,"ATZ\r","OK",4,"reset modem2",err)
	err=Modemcommand(port,"ATE0; V1\r","OK",4,"echo & verbose",err)
	err=Modemcommand(port,"AT+CSCA=\"0046708000708\"\r","OK",2,"setup SMS CENTER",err)
	err=Modemcommand(port,"AT+CMGF=0\r","OK",2,"set PDU mode",err)
	err=Modemcommand(port,"AT\r","OK",1,"AT",err)
	//	err=Modemcommand(port,"AT&F\r","OK",1,"reset modem2",nil)
	//	err=Modemcommand(port,"AT+CSMP=17,167,0,16\r","OK",1,"Flash SMS ON",err)
	// err=Modemcommand(port,"AT+CSMP=17,167,0,0\r","OK",2,"Flash SMS OFF",err)
	// err=Modemcommand(port,"\032\r","",2,"wakeup",err)
	// err=Modemcommand(port,"AT+CREG?\r","",3,"check registration",err)
	// err=Modemcommand(port,"AT+CGDCONT?\r","",3,"check CGDCONT",err)
	// err=Modemcommand(port,"AT+CGATT?\r","",3,"check GATT",err)
	// err=Modemcommand(port,"AT+CFUN=0\r","",3,"set full func off",err)
	// err=Modemcommand(port,"ATZ\r","",1,"reset modem",err)
	// err=Modemcommand(port,"ATE0\r","",1,"set echo on",err)
	// err=Modemcommand(port,"ATV1\r","",1,"set verbose on",err)
	// err=Modemcommand(port,"AT\r","OK",1,"test AT command",err)
	// err=Modemcommand(port,"AT+CFUN=1\r","",2,"set full func on",err)
	// err=Modemcommand(port,"AT\r","OK",1,"test AT command",err)
	// Set SMS center number = Universal 	= 0046735480000
	// Set SMS center number = Telenor 		= 0046708000708
	// err=Modemcommand(port,"AT\r","OK",time.Second,"test AT command",err)
	// err=Modemcommand(port,string("AT+DEVCONINFO\r"),"OK",time.Second,"get device info",err)
	// err=Modemcommand(port, "AT+CGMI\r","OK",time.Second,"get manufacturer",err)		//GET MANUFACTURER
	// err=Modemcommand(port, "AT+CGMM\r","OK",time.Second,"get model",err)		//GET MODEL
	if err!=nil {
		port.Close()
		port=nil
		return nil,errors.New("Modemreset #2 failed: " + err.Error())
	}
	// r should be "AT+CFUN=0,0\r\rOK\rAT+CMGF=0\rATE1\r\rOK\rAT+CFUN=1,0\r\rOK\r"?!
	return port,nil
}
func SendDirectSMS(port serial.Port, phoneNumber string, message string) error {
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
	err =Modemcommand(port,"AT\r","OK",5,"AT",nil)
	if err !=nil {
		port.Close()
		return errors.New("SendSMS #2 AT to wakeup failed: " + err.Error())
	}
	pduarray = CreateLongPDU(phoneNumber, message)
	for i := 0; i < len(pduarray); i++ {
		cmd1 = append(cmd1, "AT+CMGS="+fmt.Sprintf("%d", (len(pduarray[i])-2)/2)+"\r")
		cmd2 = append(cmd2, pduarray[i]+string(rune(26)))
	}
	for i := 0; i < len(cmd1); i++ {
		err =Modemcommand(port,cmd1[i],">",2,"length",nil)
		if err !=nil {
			port.Close()
			return errors.New("SendSMS #3 myread failed: " + err.Error())
		}
		err =Modemcommand(port,cmd2[i],"OK",30,"PDU",nil)
		if err !=nil {
			port.Close()
			return errors.New("SendSMS #4 myread failed: " + err.Error())
		}
	}
	return nil
}
func SendStoredSMS(port serial.Port, phoneNumber string, message string) error {
	var pduarray []string
	var cmd1 []string
	var cmd2 []string
	var r string
	var err error
	if port == nil {
		port, err = Openmodemport("")
		if err != nil {
			return errors.New("SendSMS #1 openmodemport failed: " + err.Error())
		}
	}
	err =Modemcommand(port,"AT+CMGD=1,4\r" ,"OK",2,"delete all stored sms",nil)
	pduarray = CreateLongPDU(phoneNumber, message)
	for i := 0; i < len(pduarray); i++ {
		cmd1 = append(cmd1, "AT+CMGW="+fmt.Sprintf("%d", (len(pduarray[i])-2)/2)+"\r")
		cmd2 = append(cmd2, pduarray[i]+string(rune(26)))
	}
	for i := 0; i < len(cmd1); i++ {
		err =Modemcommand(port,cmd1[i],">",3 ,"length",nil)
		if err !=nil {
			port.Close()
			return errors.New("SendSMS #2 myread failed: " + err.Error())
		}
		err =mywrite(port,cmd2[i])
		if err !=nil {
			port.Close()
			return errors.New("SendSMS #3 mywrite failed: " + err.Error())
		}
		log.Println("Modemcommand '"+cmd2[i]+"' sent.\r\n")
		err =mywrite(port,"\032")
		}
	r, err = myread(port, "C", 50)
	if err != nil {
		return errors.New("SendSMS #4 myread error: " + showdebugmsg(r)+" = "+err.Error())
	}
	if !strings.Contains(r,"+CMGW:") {
		port.Close()
		return errors.New("SendSMS #5 myread failed: no +CMGW: response, got: " + showdebugmsg(r))
	}
	log.Println("Modemcommand response: '"+showdebugmsg(r)+"'.\r\n")
	var n int
	cleanStr:= regexp.MustCompile(`[^0-9]+`).ReplaceAllString(r,"")
	n, err = strconv.Atoi(cleanStr)
	if err != nil {
		return errors.New("SendSMS #6 strconv err failed: bad result from: " + showdebugmsg(r)+" = "+err.Error())
	}
	sndstring:=fmt.Sprintf("AT+CMSS=%d\r\n",n)
	fmt.Println(sndstring)
	err =Modemcommand(port,sndstring,"OK",10,"send sms",nil)
	if err != nil {
		return errors.New("SendSMS #7 No OK: " + err.Error())
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
func myread(port serial.Port,response string,timeout int) (string,error) {
	var r string
	var err error
	var n int
	var myreadtimeout time.Duration= time.Second
	var myreadtimestart time.Time= time.Now()
	var startTime time.Time = time.Now()
	var buff []byte = make([]byte, 100)
	for {
		n, err = port.Read(buff)
		if err != nil {
			break 
		}
		if n > 0 {
			r= fmt.Sprintf("%v", string(buff[:n]))
		} 
		if strings.Contains(r,response) {
			break
		}
		if response=="" {
			break
		}
		if time.Since(myreadtimestart) > myreadtimeout {
			myreadtimestart = time.Now()
			timeout--
            log.Printf("myread timeout %dms, waiting for %d (got '%s').\r\n", time.Since(startTime).Milliseconds(),timeout,showdebugmsg(r))
		}
		if timeout<=0 {
            return r,fmt.Errorf("myread timeout exceeded, no expected response: %s within %d seconds (got '%s')", response, timeout, showdebugmsg(r))
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
	mode := &serial.Mode{
		BaudRate: 115200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	port, err = serial.Open(comport, mode)
	if err != nil {
		return err // err.Error("#1 serial.Open(comport)")
	}
	err = Modemcommand(port, "AT\r", "", 2, "wakeup", err)
	// if n==0 {
	// 	port.Close()
	// 	return errors.New("TestPort error got no data while from port "+comport)
	// }
	port.Close()
	if err != nil {
		return errors.New("TestPort error on port "+comport+": " + err.Error())
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

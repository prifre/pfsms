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
	"fmt"
	"log"
	"strings"
	"time"
	"unicode/utf16"
	"unicode/utf8"

	"fyne.io/fyne/v2"
	"go.bug.st/serial"
)

func Modemreset(comport string) bool {
	var port serial.Port
	var r string
	var err error
	mode := &serial.Mode{
		BaudRate: 115200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	port, err = serial.Open(comport, mode)
	if err != nil {
		log.Fatal("#1 serial.Open(comport)", err)
	}

	port.Break(time.Second)
	r = ""

	mywrite(port,string("AT+DEVCONINFO\r\n"))
	r += myread(port,"OK")
	mywrite(port, "AT+CGMM\r\n")
	r += myread(port,"")
	// model := myread(port)
	// model += myread(port)
	// model = strings.TrimSpace(model)
	// fmt.Println("MODEL: ", model)
	mywrite(port,"\032\r\n")
	r += myread(port,"")
	mywrite(port,"AT+CFUN=0\r\n")
	r += myread(port,"OK")
	mywrite(port,"ATZ\r\n") // Reset modem
	r += myread(port,"OK")
	mywrite(port,"ATE0\r\n") // set echo on...
	r += myread(port,"OK")
	mywrite(port,"AT+CFUN=1\r\n")
	r += myread(port,"OK")
	mywrite(port,"AT+CSCA=\"+46735480000 \"\r\n")
	r += myread(port,"OK")
	mywrite(port,"AT\r\n")
	r += myread(port,"OK")
	mywrite(port, "AT+CSCA?\r\n")
	r += myread(port,"OK")
	mywrite(port, "AT+CREG?\r\n")
	r += myread(port,"OK")
	mywrite(port, "AT\r\n")
	r += myread(port,"OK")
	mywrite(port,"AT+CMGF=0\r\n") // Set PDU mode
	r += myread(port,"OK")

	port.Close()
	port=nil

	if strings.Contains(strings.ToUpper(r), "ERROR") || !strings.Contains(r, "OK") || len(r) == 0 {
		return false
	}
	// r should be "AT+CFUN=0,0\r\n\r\nOK\r\nAT+CMGF=0\r\nATE1\r\n\r\nOK\r\nAT+CFUN=1,0\r\n\r\nOK\r\n"?!
	return true
}

func SendSMS(comport string, phoneNumber string, message string) bool {
	var pduarray []string
	var cmd1 []string
	var cmd2 []string
	var r string
	var err error
	var port serial.Port 
	mode := &serial.Mode{
		BaudRate: 115200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	port, err = serial.Open(comport, mode)
	if err != nil {
		log.Fatal("#1 serial.Open(comport)", err)
	}
	pduarray = CreateLongPDU(phoneNumber, message)
	for i := 0; i < len(pduarray); i++ {
		cmd1 = append(cmd1, "AT+CMGS="+fmt.Sprintf("%d", (len(pduarray[i])-2)/2)+"\r\n")
		cmd2 = append(cmd2, pduarray[i]+string(rune(26)))
	}
	for i := 0; i < len(cmd1); i++ {
		mywrite(port,cmd1[i])
		r = myread(port,">")
		if !strings.Contains(r, ">") {
			log.Printf("ERROR #1: no '>' in part %d: %s", i, r)

			port.Close()
			port=nil
		
			return false
		}
		mywrite(port,cmd2[i])
		r = myread(port,"OK")
		if !strings.Contains(r, "OK") {
			log.Printf("ERROR #2: no 'OK' in part %d: %s", i, r)
			port.Close()
			port=nil
			return false
		}
	}

	port.Close()
	return true
}
func mywrite(port serial.Port,s string ) {
	if fyne.CurrentApp().Preferences().Bool("debug") {
		log.Printf("WROTE: %s\r\n",showdebugmsg(s))
	}
	port.Write([]byte(s))
	port.Drain()
}
func myread(port serial.Port,response string) string {
	var r string
	timeout:=time.Millisecond * 700
	port.SetReadTimeout(timeout)
	buff := make([]byte, 100)
	startTime := time.Now()
	for {
		n, err := port.Read(buff)
		if err != nil {
			log.Fatal(err)
			break
		}
		if n > 0 {
			r += fmt.Sprintf("%v", string(buff[:n]))
		} 
		if  r>"" {
			if strings.Contains(r,response) || response==""{
				break
			}
		}
		if time.Since(startTime) > 5*time.Second {
            break
        }
	}
	if fyne.CurrentApp().Preferences().Bool("debug") {
		log.Printf("READ: %s (?='%s')\r\n",showdebugmsg(r), response)
	}
	return r
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
func showdebugmsg(s string) string {
	r2 := s
	r2 = strings.Replace(r2, string(rune(13)), "\\r", -1)
	r2 = strings.Replace(r2, string(rune(10)), "\\n", -1)
	r2 = strings.Replace(r2, string(rune(0)), "\\0", -1)
	r2 = strings.Replace(r2, string(rune(9)), "\\t", -1)
	r2 = strings.Replace(r2, string(rune(26)), "\\z", -1)
	return r2
}

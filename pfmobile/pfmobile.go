package pfmobile

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

type SMStype struct {
	mydebug 	bool
	Comport 	string 
	timeout 	time.Duration
	starttime 	time.Time
	port		serial.Port
	Addhash		bool
}

func (s *SMStype) SendMessage(phonenumbers []string, message string) [][]string {
	// Replace with the correct serial port of the modem
	s.Comport = fyne.CurrentApp().Preferences().StringWithFallback("mobilePort", "COM2")
	s.Addhash=fyne.CurrentApp().Preferences().Bool("addHash")
	
	var sendtext, phoneNumber string
	var failures, success int
	var result [][]string
	var err error
	s.mydebug= true
	// s.Setuplog()
	s.starttime = time.Now()
	s.timeout = time.Millisecond * 700
	mode := &serial.Mode{
		BaudRate: 115200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	s.port, err = serial.Open(s.Comport, mode)
	if err != nil {
		log.Printf("#1 serial.Open(comport) %s\r\n", err.Error())
		return nil
	}
	s.port.SetReadTimeout(s.timeout)
	message = strings.TrimSpace(message)
	log.Printf("Got %d phonenumbers to send ok.\r\n", len(phonenumbers))
	for i, record := range phonenumbers {
		rec:=strings.Split(record,"\t")
		phoneNumber = rec[0]
		sendtext=message
		if strings.Contains(sendtext,"<<Fname>>") || strings.Contains(sendtext,"<<Lname>>") {
			sendtext = strings.Replace(sendtext, "<<Fname>>", rec[1], -1)
			sendtext = strings.Replace(sendtext, "<<Lname>>", rec[2], -1)
		}
		if s.Addhash {
			sendtext = fmt.Sprintf(sendtext+"\r\n#=%d", i+1)
		}
		sentok := false
		for !sentok {
			sentok = s.SendSMS(phoneNumber, sendtext)
			if !sentok {
				log.Println("--------------------SENDSMS FAILED")
				modemresetfail:=0
				for !s.Modemreset() && modemresetfail<10 {
					log.Println("--------------------MODEMRESET FAIL: ",modemresetfail)
					modemresetfail++
				}
				if modemresetfail>8 {
					return nil
				}
				log.Println("--------------------MODEMRESET OK")
				failures++
			}
		}
		success++
		log.Printf("Message %d/%d to phone %s sent! (failures: %d)\r\n", i+1, len(phonenumbers), phoneNumber, failures)
		tstamp:=time.Now().Format("200601012150405")
		result =append(result,[]string{tstamp,phoneNumber,sendtext})
		if !s.mydebug {
			log.Printf("%s Message %d/%d to phone %s sent! (failures: %d)\r\n", time.Now().Format("2006-01-02 15:04:05"), i+1, len(phonenumbers), phoneNumber, failures)
		}
	}
	log.Printf("RESULT OF SMS SENDING: Failures: %d Success: %d\r\n", failures, success)
	s1:=s.starttime.Format("2006-01-02 15:04:05")
	s2:=time.Now().Format("2006-01-02 15:04:05")
	log.Printf("Started: %s  Finished: %s  Duration: %s\r\n",s1 , s2, time.Since(s.starttime))
	log.Printf("Speed: %ds/sms\r\n", int(time.Since(s.starttime).Seconds())/len(phonenumbers))
	s.port.Close()
	return result
}
func (s SMStype) Modemreset() bool {
	var r string
	// var model string
	s.port.Break(time.Second)
	// port.Write([]byte(string("AT+DEVCONINFO\r\n")))
	// slowwrite(s.port, "AT+CGMM\r\n")
	// model = myread(s.port)
	// model += myread(s.port)
	// model = strings.TrimSpace(model)
	// fmt.Println("MODEL: ", model)
	r = ""
	s.port.Write([]byte("\032\r\n"))
	r += myread(s.port)
	s.port.Write([]byte("AT+CFUN=0\r\n"))
	r += myread(s.port)
	s.port.Write([]byte("ATZ\r\n")) // set echo on...
	r += myread(s.port)
	s.port.Write([]byte("ATE0\r\n")) // set echo on...
	r += myread(s.port)
	s.port.Write([]byte("AT+CFUN=1\r\n"))
	r += myread(s.port)
	s.port.Write([]byte("AT+CSCA=\"+46735480000\""))
	r += myread(s.port)
	slowwrite(s.port, "AT+CSCA?\r\n")
	r += myread(s.port)
	s.port.Write([]byte("AT+CMGF=0\r\n")) // Set PDU mode
	r += myread(s.port)
	// if mydebug {
		fmt.Println("MODEMRESET: ", showdebugmsg(r))
	// }
	if strings.Contains(strings.ToUpper(r), "ERROR") || !strings.Contains(r, "OK") || len(r) == 0 {
		return false
	}
	// r should be "AT+CFUN=0,0\r\n\r\nOK\r\nAT+CMGF=0\r\nATE1\r\n\r\nOK\r\nAT+CFUN=1,0\r\n\r\nOK\r\n"?!
	return true
}
func (s SMStype) SendSMS(phoneNumber string, message string) bool {
	var pduarray []string
	var cmd1 []string
	var cmd2 []string
	var r string
	pduarray = CreateLongPDU(phoneNumber, message)
	for i := 0; i < len(pduarray); i++ {
		cmd1 = append(cmd1, "AT+CMGS="+fmt.Sprintf("%d", (len(pduarray[i])-2)/2)+"\r\n")
		cmd2 = append(cmd2, pduarray[i]+string(rune(26)))
	}
	for i := 0; i < len(cmd1); i++ {
		s.port.Write([]byte(cmd1[i]))
		r = myread(s.port)
		if !strings.Contains(r, ">") {
			fmt.Printf("ERROR #1: no '>' in part %d: %s", i, r)
			return false
		}
		s.port.Write([]byte(cmd2[i]))
		r = myread(s.port)
		r += myread(s.port)
		r += myread(s.port)
		if !strings.Contains(r, "OK") {
			fmt.Printf("ERROR #2: no 'OK' in part %d: %s", i, r)
			return false
		}
	}
	return true
}
func myread(port serial.Port) string {
	var r string
	buff := make([]byte, 100)
	for {
		n, err := port.Read(buff)
		if err != nil {
			log.Println("Error #1 myread()", err.Error())
			break
		}
		if n == 0 {
			break
		}
		r += fmt.Sprintf("%v", string(buff[:n]))
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
func slowwrite(port serial.Port, s string) {
	// port.Drain()
	// port.ResetOutputBuffer()
	for i := 0; i < len(s); i++ {
		port.Write([]byte(string(s[i]))) // sms central...
	}
	port.Drain()
}
func (s SMStype) GetPortsList() ([]string, error) {
		return serial.GetPortsList()
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

package ariasms

// Special program to send a lots of sms using a mobile phone!
// uses logfilename for logging
// uses phonenumbersfilename to specify file with phonenumbers
// 2024-01-21 working!!!!
// 2024-03-10 switched to newer serial driver, implemented support for S24U and model selection
// got it working with Samsung S24Ultra! speed 14s/sms using timeout = Millisecond*700

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
	"unicode/utf16"
	"unicode/utf8"

	"go.bug.st/serial"
)

const logfilename = "ariasms.log"
const phonenumbersfilename = "ariasms.tel"
const messagefilename = "ariasms.txt"
const timeout = time.Millisecond * 700
const comport = "COM2"

var wrt io.Writer
var mydebug bool = true

func main() {
	// Replace with the correct serial port of the modem
	var message, sendtext, phoneNumber, Fname, Lname string
	var failures, success int
	var phonenumbers []string
	var err error
	var port serial.Port
	Setuplog()
	starttime := time.Now()
	hname, _ := os.Hostname()
	log.Println("Starting " + hname)
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
	port.SetReadTimeout(timeout)
	var m []byte
	m, err = os.ReadFile(messagefilename)
	if err != nil {
		log.Fatal("#2 ariasms.txt error", err)
	}
	message = string(m)
	// message += "Hej <<Fname>>,\r\n"
	// message += "Den 14/3, 14-18 kommer Per Ståbi att hålla en kalligrafikurs på Sollentuna Ram.\r\n"
	// message += "De som inte betalat kursen ännu kan gärna Swisha 1700:- till 0736290839.\r\n"
	// message += "Vid frågor, ring!\r\n"
	// message += "Med glada hälsningar\r\n"
	// message += "Peter & Per"
	// message = "test"
	phonenumbers = getphonenumbers(phonenumbersfilename)
	log.Printf("Read %d from %s ok.\r\n", len(phonenumbers), phonenumbersfilename)
	for i, record := range phonenumbers {
		if !strings.Contains(record, "\t") {
			phoneNumber = record
		} else {
			rec := strings.Split(record, "\t")
			phoneNumber = rec[0]
			if len(rec) > 0 {
				Fname = rec[1]
			}
			if len(rec) > 1 {
				Lname = rec[2]
			}
		}
		sendtext = fmt.Sprintf(message+"\r\n#=%d", i+1)
		if Lname > "" {
			sendtext = strings.Replace(sendtext, "<<Fname>>", Fname, -1)
		}
		if Lname > "" {
			sendtext = strings.Replace(sendtext, "<<Lname>>", Lname, -1)
		}
		sentok := false
		for !sentok {
			sentok = sendSMS(port, phoneNumber, sendtext)
			if !sentok {
				log.Println("--------------------SENDSMS FAILED")
				for !modemreset(port) {
					log.Println("--------------------MODEMRESET FAILED")
				}
				log.Println("--------------------MODEMRESET OK")
				failures++
			}
		}
		success++
		log.Printf("Message %d/%d to phone %s sent! (failures: %d)\r\n", i+1, len(phonenumbers), phoneNumber, failures)
		documentsms(time.Now(), phoneNumber, showdebugmsg(sendtext))
		if !mydebug {
			fmt.Printf("%s %s Message %d/%d to phone %s sent! (failures: %d)\r\n", time.Now().Format("2006-01-02"), time.Now().Format("15:04"), i+1, len(phonenumbers), phoneNumber, failures)
		}
	}
	log.Printf("\r\nRESULT OF SMS SENDING\r\nFailures: %d\r\nSuccess: %d\r\n", failures, success)
	log.Printf("Started: %s\r\nFinished: %s\r\nDuration: %s\r\n", starttime, time.Now(), time.Since(starttime))
	fmt.Printf("Speed: %ds/sms", int(time.Since(starttime).Seconds())/len(phonenumbers))
	port.Close()
}
func sendSMS(port serial.Port, phoneNumber string, message string) bool {
	var pduarray []string
	var cmd1 []string
	var cmd2 []string
	var r string
	pduarray = createLongPDU(phoneNumber, message)
	for i := 0; i < len(pduarray); i++ {
		cmd1 = append(cmd1, "AT+CMGS="+fmt.Sprintf("%d", (len(pduarray[i])-2)/2)+"\r\n")
		cmd2 = append(cmd2, pduarray[i]+string(rune(26)))
	}
	for i := 0; i < len(cmd1); i++ {
		port.Write([]byte(cmd1[i]))
		r = myread(port)
		if !strings.Contains(r, ">") {
			fmt.Printf("ERROR #1: no '>' in part %d: %s", i, r)
			return false
		}
		port.Write([]byte(cmd2[i]))
		r = myread(port)
		r += myread(port)
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
			log.Fatal(err)
			break
		}
		if n == 0 {
			break
		}
		r += fmt.Sprintf("%v", string(buff[:n]))
	}
	return r
}
func createLongPDU(phoneNumber string, message string) []string {
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
		pdu = createPDU(phoneNumber, segments[i], udh)             // Construct PDU for the segment with UDH
		pdus = append(pdus, pdu)
	}
	return pdus
}
func createPDU(number string, message string, udh string) string {
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
func Setuplog() {
	f, err := os.OpenFile(logfilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	//	defer f.Close()
	if mydebug { // write to both screen and file!!!
		wrt = io.MultiWriter(os.Stdout, f)
	} else {
		wrt = f
	}
	log.SetOutput(wrt)
}
func getphonenumbers(fn string) []string {
	// Read the text file
	data, err := os.ReadFile(fn)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	// Convert the data to a string
	text := string(data)
	text = strings.TrimSpace(text)
	text = strings.Replace(text, "\n", "\r", -1)
	text = strings.Replace(text, "\r\r", "\r", -1)
	text = strings.Replace(text, "\r0", "\r+46", -1)
	if text[:1] == "0" {
		text = "+46" + text[1:]
	}

	// Split the text into an array of phone numbers
	if !strings.Contains(text, "\t") && strings.Contains(text, " ") {
		text = strings.Replace(text, " ", "\t", -1)
		for strings.Contains(text, "\t\t") {
			text = strings.Replace(text, "\t\t", "\t", -1)
		}
	}
	phoneNumbers := strings.Split(text, "\r")
	return phoneNumbers
}
func documentsms(t time.Time, p string, msg string) {
	fn := strings.Replace(phonenumbersfilename, ".txt", ".log", -1)
	f, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	p = strings.Replace(p, "+46", "0", -1)
	s := fmt.Sprintf("%s\t%s\t%s\t%s\r\n", t.Format("2006-01-02"), t.Format("15:04"), p, msg)
	f.WriteString(s)
	f.Close()
}
func modemreset(port serial.Port) bool {
	var r string
	port.Break(time.Second)
	// port.Write([]byte(string("AT+DEVCONINFO\r\n")))
	// slowwrite(port, "AT+CGMM\r\n")
	// model := myread(port)
	// model += myread(port)
	// model = strings.TrimSpace(model)
	// fmt.Println("MODEL: ", model)
	r = ""
	port.Write([]byte("\032\r\n"))
	r += myread(port)
	port.Write([]byte("AT+CFUN=0\r\n"))
	r += myread(port)
	port.Write([]byte("ATZ\r\n")) // set echo on...
	r += myread(port)
	port.Write([]byte("ATE0\r\n")) // set echo on...
	r += myread(port)
	port.Write([]byte("AT+CFUN=1\r\n"))
	r += myread(port)
	port.Write([]byte("AT+CSCA=\"+46735480000\""))
	r += myread(port)
	slowwrite(port, "AT+CSCA?\r\n")
	r += myread(port)
	port.Write([]byte("AT+CMGF=0\r\n")) // Set PDU mode
	r += myread(port)
	if mydebug {
		fmt.Println("MODEMRESET: ", showdebugmsg(r))
	}
	if strings.Contains(strings.ToUpper(r), "ERROR") || !strings.Contains(r, "OK") || len(r) == 0 {
		return false
	}
	// r should be "AT+CFUN=0,0\r\n\r\nOK\r\nAT+CMGF=0\r\nATE1\r\n\r\nOK\r\nAT+CFUN=1,0\r\n\r\nOK\r\n"?!
	return true
}
func slowwrite(port serial.Port, s string) {
	// port.Drain()
	// port.ResetOutputBuffer()
	for i := 0; i < len(s); i++ {
		port.Write([]byte(string(s[i]))) // sms central...
	}
	port.Drain()
}

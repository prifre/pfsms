package pfmobile

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	"unicode/utf16"

	"fyne.io/fyne/v2"
	"go.bug.st/serial"
)

// Global räknare för UDH-referensnummer vid flerdels-SMS
var globalMsgRef uint32 = 0

func Modemcommand(port serial.Port, modemcommand string, expectedresponse string, timeout int, description string, err error) error {
	if err != nil {
		return err
	}
	// Sätt read timeout i millisekunder utifrån skickad variabel
	err = port.SetReadTimeout(time.Duration(timeout) * 100 * time.Millisecond)
	if err != nil {
		log.Printf("Modemcommand Error #1 port.SetReadTimeout(%d) failed: %s\r\n", timeout, err)
		return err
	}
	err = mywrite(port, modemcommand)
	if err != nil {
		return err
	}
	r, err := myread(port, expectedresponse, timeout)
	if err != nil {
		return errors.New("Modemcommand Error #2 " + description + " readerror '" + showdebugmsg(r) + " " + err.Error() + "'")
	}
	if expectedresponse != "" && !strings.Contains(r, expectedresponse) {
		return errors.New("Modemcommand Error #3 " + description + " unexpected response: '" + showdebugmsg(r) + "' expected: '" + expectedresponse + "'")
	}
	return nil
}

func Openmodemport(comport string) (serial.Port, error) {
	var port serial.Port
	var err error
	if comport == "" {
		comport = fyne.CurrentApp().Preferences().StringWithFallback("mobileport", "")
		if comport == "" {
			comport, err = GetMobilePort()
			if err != nil {
				return nil, errors.New("Openmodemport #0 GetMobilePort() failed: " + err.Error())
			}
		}
	}
	mode := &serial.Mode{
		BaudRate: 115200,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}
	port, err = serial.Open(comport, mode)
	if err != nil {
		return nil, errors.New("Openmodemport #2 failed: " + err.Error())
	}
	port.Drain()
	port.ResetInputBuffer()
	return port, nil
}

func Modemreset(comport string) (serial.Port, error) {
	var port serial.Port
	var err error
	if comport == "" {
		comport, err = GetMobilePort()
		if err != nil {
			return nil, errors.New("Modemreset #1 GetMobilePort failed: " + err.Error())
		}
	}
	port, err = Openmodemport(comport)
	if err != nil {
		return nil, errors.New("Modemreset #2 Openmodemport failed: " + err.Error())
	}
	err = Modemcommand(port, "\032\r", "", 1, "ESC", err)
	err = Modemcommand(port, "AT\r", "OK", 1, "AT", err)
	err = Modemcommand(port, "AT+CFUN=1\r", "OK", 10, "reset modem", nil)
	err = Modemcommand(port, "ATZ\r", "OK", 4, "reset modem2", err)
	err = Modemcommand(port, "ATE0; V1\r", "OK", 4, "echo & verbose", err)
	err = Modemcommand(port, "AT+CMGF=0\r", "OK", 2, "set PDU mode", err)
	err = Modemcommand(port, "AT\r", "OK", 1, "AT", err)

	if err != nil {
		port.Close()
		return nil, errors.New("Modemreset #2 failed: " + err.Error())
	}
	return port, nil
}

func SendDirectSMS(port serial.Port, phoneNumber string, message string) error {
	var pduarray []string
	var cmd1 []string
	var cmd2 []string
	var err error
	var closeport bool = false

	if port == nil {
		closeport = true
		port, err = Openmodemport("")
		if err != nil {
			return errors.New("SendSMS #1 openmodemport failed: " + err.Error())
		}
	}

	err = Modemcommand(port, "AT\r", "OK", 5, "AT", nil)
	if err != nil {
		if closeport {
			port.Close()
		}
		return errors.New("SendSMS #2 AT to wakeup failed: " + err.Error())
	}

	pduarray = CreateLongPDU(phoneNumber, message)
	for i := 0; i < len(pduarray); i++ {
		cmd1 = append(cmd1, fmt.Sprintf("AT+CMGS=%d\r", (len(pduarray[i])-2)/2))
		cmd2 = append(cmd2, pduarray[i]+string(rune(26)))
	}

	for i := 0; i < len(cmd1); i++ {
		err = Modemcommand(port, cmd1[i], ">", 2, "length", nil)
		if err != nil {
			if closeport {
				port.Close()
			}
			return errors.New("SendSMS #3 myread failed: " + err.Error())
		}
		err = Modemcommand(port, cmd2[i], "OK", 30, "PDU", nil)
		if err != nil {
			if closeport {
				port.Close()
			}
			return errors.New("SendSMS #4 myread failed: " + err.Error())
		}
	}

	if closeport {
		port.Close()
	}
	time.Sleep(time.Millisecond * 500)
	return nil
}

func SendStoredSMS(port serial.Port, phoneNumber string, message string) error {
	var err error
	var closeport bool = false

	if port == nil {
		closeport = true
		port, err = Openmodemport("")
		if err != nil {
			return errors.New("SendStoredSMS #1 openmodemport failed: " + err.Error())
		}
	}

	// Se till att porten alltid stängs om den öppnades i denna funktion
	defer func() {
		if closeport && port != nil {
			port.Close()
		}
	}()

	err = Modemcommand(port, "AT+CMGD=1,4\r", "OK", 2, "delete all stored sms", nil)
	if err != nil {
		return errors.New("SendStoredSMS delete failed: " + err.Error())
	}

	pduarray := CreateLongPDU(phoneNumber, message)
	storedIndexes := make([]int, 0, len(pduarray))

	re := regexp.MustCompile(`\+CMGW:\s*(\d+)`)

	// Lagra alla PDU-delar i modemet först
	for i := 0; i < len(pduarray); i++ {
		cmd1 := fmt.Sprintf("AT+CMGW=%d\r", (len(pduarray[i])-2)/2)
		cmd2 := pduarray[i] + "\032"

		err = Modemcommand(port, cmd1, ">", 3, "CMGW length", nil)
		if err != nil {
			return errors.New("SendStoredSMS CMGW prompt failed: " + err.Error())
		}

		err = mywrite(port, cmd2)
		if err != nil {
			return errors.New("SendStoredSMS mywrite failed: " + err.Error())
		}

		// Läs svaret specifikt för +CMGW
		r, err := myread(port, "+CMGW:", 50)
		if err != nil {
			return errors.New("SendStoredSMS CMGW response read error: " + err.Error())
		}

		matches := re.FindStringSubmatch(r)
		if len(matches) < 2 {
			return errors.New("SendStoredSMS could not parse +CMGW index from response: " + showdebugmsg(r))
		}

		idx, err := strconv.Atoi(matches[1])
		if err != nil {
			return errors.New("SendStoredSMS invalid index: " + err.Error())
		}
		storedIndexes = append(storedIndexes, idx)
	}

	// Skicka alla lagrade PDU-segment
	for _, idx := range storedIndexes {
		sndstring := fmt.Sprintf("AT+CMSS=%d\r", idx)
		err = Modemcommand(port, sndstring, "OK", 15, "send stored sms", nil)
		if err != nil {
			return fmt.Errorf("SendStoredSMS AT+CMSS=%d failed: %w", idx, err)
		}
	}
	return nil
}

func mywrite(port serial.Port, s string) error {
	n, err := port.Write([]byte(s))
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("no bytes written to port")
	}
	err = port.Drain()
	if err != nil {
		return err
	}
	return nil
}

func myread(port serial.Port, response string, timeout int) (string, error) {
	var r string
	var err error
	var n int
	var myreadtimeout time.Duration = time.Second
	var myreadtimestart time.Time = time.Now()
	var startTime time.Time = time.Now()
	var buff []byte = make([]byte, 100)

	if response == "" {
		time.Sleep(100 * time.Millisecond)
		n, _ = port.Read(buff)
		if n > 0 {
			r = string(buff[:n])
		}
		return r, nil
	}

	for {
		n, err = port.Read(buff)
		if n > 0 {
			r += string(buff[:n])
			if strings.Contains(r, response) {
				return r, nil
			}
		}

		if time.Since(myreadtimestart) > myreadtimeout {
			myreadtimestart = time.Now()
			timeout--
		}

		if timeout <= 0 {
			var errStr string
			if err != nil {
				errStr = " (last error: " + err.Error() + ")"
			}
			return r, fmt.Errorf("myread started %s timeout exceeded, no expected response: %s within target time (got '%s')%s",
				startTime.Format("15:04:05"), response, showdebugmsg(r), errStr)
		}
	}
}

func CreateLongPDU(phoneNumber string, message string) []string {
	const maxCharsPerSegment = 67
	var pdus []string

	runes := []rune(message)
	totalRunes := len(runes)

	if totalRunes <= 70 {
		return []string{CreatePDU(phoneNumber, message, "")}
	}

	var segments []string
	for i := 0; i < totalRunes; i += maxCharsPerSegment {
		end := i + maxCharsPerSegment
		if end > totalRunes {
			end = totalRunes
		}
		segments = append(segments, string(runes[i:end]))
	}

	// Generera unikt referensnummer (1-255) per långt meddelande
	refNum := atomic.AddUint32(&globalMsgRef, 1) % 256
	if refNum == 0 {
		refNum = 1
	}

	for i, seg := range segments {
		udh := fmt.Sprintf("050003%02X%02X%02X", refNum, len(segments), i+1)
		pdus = append(pdus, CreatePDU(phoneNumber, seg, udh))
	}
	return pdus
}

func CreatePDU(number string, message string, udh string) string {
	var pdu, pduHeader, pduMessage, pduMessageLen string

	// Tvätta telefonnumret
	phoneNumber := strings.TrimPrefix(number, "+")
	phoneNumber = strings.TrimPrefix(phoneNumber, "00")

	// Om svenskt nummer som börjar på '07', omvandla till internationellt '467'
	if strings.HasPrefix(phoneNumber, "0") {
		phoneNumber = "46" + phoneNumber[1:]
	}

	// Lägg till padding-siffra (F) om längden är udda
	if len(phoneNumber)%2 != 0 {
		phoneNumber += "F"
	}

	// Omvandla till semi-octets
	semiOctets := ""
	for i := 0; i < len(phoneNumber); i += 2 {
		semiOctets = semiOctets + string(phoneNumber[i+1]) + string(phoneNumber[i])
	}

	pduLength := len(strings.TrimSuffix(phoneNumber, "F"))

	encodedMessage := utf16.Encode([]rune(message))
	ucs2EncodedMessage := ""
	for _, char := range encodedMessage {
		ucs2EncodedMessage += fmt.Sprintf("%04X", char)
	}
	pduMessage = ucs2EncodedMessage

	if udh == "" {
		pduHeader = fmt.Sprintf("001100%02X91%s00080B", pduLength, semiOctets)
		pduMessageLen = fmt.Sprintf("%02X", len(pduMessage)/2)
		pdu = pduHeader + pduMessageLen + pduMessage
	} else {
		pduHeader = fmt.Sprintf("005100%02X91%s00080B", pduLength, semiOctets)
		pduMessageLen = fmt.Sprintf("%02X", (len(pduMessage)/2)+6) // 6 bytes UDH
		pdu = pduHeader + pduMessageLen + udh + pduMessage
	}

	return strings.ToUpper(pdu)
}

func GetPortsList() ([]string, error) {
	return serial.GetPortsList()
}

func TestPort(comport string) error {
	mode := &serial.Mode{
		BaudRate: 115200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	port, err := serial.Open(comport, mode)
	if err != nil {
		return err
	}
	err = Modemcommand(port, "AT\r", "OK", 2, "wakeup", nil)
	port.Close()
	if err != nil {
		return errors.New("TestPort error on port " + comport + ": " + err.Error())
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
	r2 = strings.ReplaceAll(r2, "\r", "\\r")
	r2 = strings.ReplaceAll(r2, "\n", "\\n")
	r2 = strings.ReplaceAll(r2, "\x00", "\\0")
	r2 = strings.ReplaceAll(r2, "\t", "\\t")
	r2 = strings.ReplaceAll(r2, "\x1a", "\\z")
	return r2
}

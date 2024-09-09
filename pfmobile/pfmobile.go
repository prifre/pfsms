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

	"fyne.io/fyne/v2"
	"go.bug.st/serial"
)

type SMStype struct {
	mydebug   bool
	Comport   string
	timeout   time.Duration
	starttime time.Time
	Addhash   bool
}

func (s *SMStype) SendMessage(phonenumbers []string, message string) [][]string {
	// Replace with the correct serial port of the modem
	if s.Comport=="" {
		s.Comport = fyne.CurrentApp().Preferences().StringWithFallback("mobileport", "COM2")
		s.Addhash = fyne.CurrentApp().Preferences().Bool("addhash")
	}
	var sendtext, phoneNumber string
	var failures, success int
	var result [][]string
	s.mydebug = true
	// s.Setuplog()
	s.starttime = time.Now()
	s.timeout = time.Millisecond * 700
	message = strings.TrimSpace(message)
	log.Printf("Got %d phonenumbers to send ok.\r\n", len(phonenumbers))
	for i, record := range phonenumbers {
		rec := strings.Split(record, "\t")
		phoneNumber = rec[0]
		sendtext = message
		if strings.Contains(sendtext, "<<Fname>>") || strings.Contains(sendtext, "<<Lname>>") {
			sendtext = strings.Replace(sendtext, "<<Fname>>", rec[1], -1)
			sendtext = strings.Replace(sendtext, "<<Lname>>", rec[2], -1)
		}
		if s.Addhash {
			sendtext = fmt.Sprintf(sendtext+"\r\n#=%d", i+1)
		}
		sentok := false
		for !sentok {
			sentok = SendSMS(s.Comport,phoneNumber, sendtext)
			if !sentok {
				log.Println("--------------------SENDSMS FAILED")
				modemresetfail := 0
				for !Modemreset(s.Comport) && modemresetfail < 10 {
					log.Println("--------------------MODEMRESET FAIL: ", modemresetfail)
					modemresetfail++
				}
				if modemresetfail > 8 {
					return nil
				}
				log.Println("--------------------MODEMRESET OK")
				failures++
			}
		}
		success++
		log.Printf("Message %d/%d to phone %s sent! (failures: %d)\r\n", i+1, len(phonenumbers), phoneNumber, failures)
		tstamp := time.Now().Format("20060102150405")
		result = append(result, []string{tstamp, phoneNumber, sendtext})
		if !s.mydebug {
			log.Printf("%s Message %d/%d to phone %s sent! (failures: %d)\r\n", time.Now().Format("2006-01-02 15:04:05"), i+1, len(phonenumbers), phoneNumber, failures)
		}
	}
	log.Printf("RESULT OF SMS SENDING: Failures: %d Success: %d\r\n", failures, success)
	s1 := s.starttime.Format("2006-01-02 15:04:05")
	s2 := time.Now().Format("2006-01-02 15:04:05")
	log.Printf("Started: %s  Finished: %s  Duration: %s\r\n", s1, s2, time.Since(s.starttime))
	log.Printf("Speed: %ds/sms\r\n", int(time.Since(s.starttime).Seconds())/len(phonenumbers))
	return result
}

func (s SMStype) GetPortsList() ([]string, error) {
	return serial.GetPortsList()
}
package pfmobile

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

/*
  TEST MODEM SEQUENCE VIA Putty/Terminal:
  AT+CSQ          // Kontrollera signalstyrka (>10 rekommenderas)
  AT+CMGF=1       // Sätt text-läge (1) eller PDU-läge (0)
  AT+CSCA?        // Visa meddelandecentralens nummer
  AT+CSCS="UCS2"  // Ställ in teckenuppsättning för UTF-16 / svenska tecken (ÅÄÖ)
  AT+CMGS="num"   // Skicka SMS
*/

func TestModemreset(t *testing.T) {
	port, err := GetMobilePort()
	if err != nil {
		t.Skipf("Hoppar över modemtest, hittade ingen mobilport: %v", err)
	}

	start := time.Now()
	_, err = Modemreset(port)
	if err != nil {
		t.Fatalf("Modemreset misslyckades på %s: %v", port, err)
	}

	t.Logf("Modemreset OK på %s (tog %s)", port, time.Since(start))
}

func TestPDU_unicode(t *testing.T) {
	message := "hello\r\nhelloÅÄÖ"
	number := "+46736290839"
	want := "0011000B916437260938F900080B1E00680065006C006C006F000D000A00680065006C006C006F00C500C400D6"

	pduarray := CreateLongPDU(number, message)
	if len(pduarray) == 0 {
		t.Fatalf("CreateLongPDU returnerade en tom slice")
	}

	pdu := pduarray[0]

	if len(pdu) < 30 || len(want) < 30 {
		t.Fatalf("PDU-strängen är för kort. Got len=%d, Want len=%d", len(pdu), len(want))
	}

	// Kontrollera telefonnumrets längdfält i PDU
	wantLen, _ := strconv.ParseInt(want[6:8], 16, 32)
	pduPhoneLen, _ := strconv.ParseInt(pdu[6:8], 16, 32)
	if wantLen != pduPhoneLen {
		t.Errorf("Nummerlängd matchar inte: WANT %d, GOT %d", wantLen, pduPhoneLen)
	}

	// Kontrollera själva telefonnummer-segmentet
	if want[10:22] != pdu[10:22] {
		t.Errorf("Telefonnummer i PDU matchar inte: WANT %s, GOT %s", want[10:22], pdu[10:22])
	}

	// Kontrollera datalängd (meddelandelängd i bytes)
	wantMLen, _ := strconv.ParseInt(want[28:30], 16, 32)
	pduMLen, _ := strconv.ParseInt(pdu[28:30], 16, 32)
	if wantMLen != pduMLen {
		t.Errorf("Meddelandelängd (bytes) matchar inte: WANT %d, GOT %d", wantMLen, pduMLen)
	}

	// Exakt jämförelse av hela PDU-strängen
	if want != pdu {
		t.Errorf("PDU matchar inte helt:\nWANT: %s\nGOT:  %s", want, pdu)
	} else {
		t.Log("Alla Unicode PDU-tester godkända!")
	}
}

func TestCreateLongPDU(t *testing.T) {
	message := "Välkommen till akvarellkursen 18/1, 18.00, Sollentuna Ram."
	message = message + message + message + "X"
	number := "+46736290839"

	want := []string{
		"0051000B816437260938F900080B8C050003000301005600E4006C006B006F006D006D0065006E002000740069006C006C00200061006B0076006100720065006C006C006B0075007200730065006E002000310038002F0031002C002000310038002E00300030002C00200053006F006C006C0065006E00740075006E0061002000520061006D002E005600E4006C006B006F006D006D0065006E",
		"0051000B816437260938F900080B8C050003000302002000740069006C006C00200061006B0076006100720065006C006C006B0075007200730065006E002000310038002F0031002C002000310038002E00300030002C00200053006F006C006C0065006E00740075006E0061002000520061006D002E005600E4006C006B006F006D006D0065006E002000740069006C006C00200061006B0076",
		"0051000B816437260938F900080B58050003000303006100720065006C006C006B0075007200730065006E002000310038002F0031002C002000310038002E00300030002C00200053006F006C006C0065006E00740075006E0061002000520061006D002E0058",
	}

	pduarray := CreateLongPDU(number, message)

	if len(pduarray) != len(want) {
		t.Fatalf("Förväntade %d PDU-delar, fick %d", len(want), len(pduarray))
	}

	for g := 0; g < len(want); g++ {
		if pduarray[g] != want[g] {
			t.Errorf("Del %d matchar inte:\nWANT: %s\nGOT:  %s", g+1, want[g], pduarray[g])
		}
	}
}

func TestSendDirectMessage(t *testing.T) {
	port, err := GetMobilePort()
	if err != nil {
		t.Skipf("Hoppar över skickatest: ingen aktiv mobilport funnen (%v)", err)
	}

	starttime := time.Now()
	msg := fmt.Sprintf("This is a SendDirectMessage test, sent %s! ÅÄÖ åäö", time.Now().Format("2006-01-02 15:04:05"))

	p, err := Modemreset(port)
	if err != nil {
		t.Fatalf("Modemreset misslyckades: %v", err)
	}

	err = SendDirectSMS(p, "0046736290839", msg)
	if err != nil {
		t.Fatalf("SendDirectSMS misslyckades: %v", err)
	}

	t.Logf("SendSMS slutfördes framgångsrikt på %s (tog %s)", port, time.Since(starttime))
}

func TestGetPortsList(t *testing.T) {
	ports, err := GetPortsList()
	if err != nil {
		t.Fatalf("GetPortsList fel: %v", err)
	}
	t.Logf("Hittade serieportar: %v", ports)
}

func TestAllTestPorts(t *testing.T) {
	ports, err := GetPortsList()
	if err != nil || len(ports) == 0 {
		t.Skip("Inga tillgängliga serieportar att testa")
	}

	for _, port := range ports {
		err := TestPort(port)
		if err == nil {
			t.Logf("Port %s OK!", port)
		} else {
			t.Logf("Port %s misslyckades eller svarade inte: %v", port, err)
		}
	}
}

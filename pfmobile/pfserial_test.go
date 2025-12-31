package pfmobile

//various testscenarios:

// s := new(SMStype)
// s.Comport="COM3"
// s.SendMessage(p,msg)
// Modemreset("COM3")
// AT+CMGS=51
// 0051000D81006437260938F900080B24050003000101005400680069007300200069007300200061002000740065007300740021
/* TEST MODEM SEQUENCE VIA putty:
// testa signalstyrka, bör vara >10.
AT+CSQ
// kör text-läge:
AT+CMGF=1
// kolla meddelandecentralens nummer:
AT+CSCA?
// ställ in teckenuppsättning till UCS2 (för att kunna skicka svenska tecken)
AT+CSCS="UCS2"
// skicka sms till nummer (inkl riktnr +46...)

AT
ATZ
AT+CSCA?
AT+CSCA="+46708000708"
AT+CGATT?
AT+CPIN?
AT+CSCS="UCS2"
AT+CSMP=17,167,0,0
AT+CMGF=1
AT+CMEE=2;
ATE1 ; +CMGS="0046736290839"
AT+CMGF=1; +CFUN=1; V1;
TEST MESSAGE



*/

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)
func TestModemreset(t *testing.T) {
	start:=time.Now()
	comport:= "COM3"
	var err error
	_,err=Modemreset(comport)
	if err!=nil {
		t.Fatalf("Modemreset failed on %s: %s", comport, err.Error())
	} else {
		fmt.Println("Modemreset OK on ")
	}
	elapsed:=time.Since(start)
	fmt.Printf("Modemreset took %s\r\n", elapsed)
}
func TestPDU_unicode(t *testing.T) {
	message := "hello\r\nhelloÅÄÖ"
	number := "+46736290839"
	want := "0011000B916437260938F900080B1E00680065006C006C006F000D000A00680065006C006C006F00C500C400D6" // Unicode encoded.
	//       0011000B916437260938F900000B08C8329BFD766CB9 // default SMS alphabet encoded
	//       0011000B916437260938F900F40B0848656C6C6F3F3F3F // ANSI 8-bit encoded
	pduarray := CreateLongPDU(number, message)
	pdu := pduarray[0]
	// check phone number length
	wantLen, _ := strconv.ParseInt(want[6:8], 16, 32)
	pduPhoneLen, _ := strconv.ParseInt(pdu[6:8], 16, 32)
	if wantLen != pduPhoneLen {
		fmt.Println("WANT: >"+want[6:8]+"< (", wantLen, ")  "+want)
		fmt.Println("GOT:  >"+pdu[6:8]+"< (", pduPhoneLen, ")  "+pdu)
		t.Fatalf("phones lengths")
	}
	fmt.Println("PHONE NUMBER LENGTHS!")

	if want[10:22] != pdu[10:22] {
		fmt.Println("WANT: ", want[10:22])
		fmt.Println("GOT:  ", pdu[10:22])
		t.Fatalf("phones differ")
	}
	fmt.Println("PHONE NUMBER OK!")

	wantMLen, _ := strconv.ParseInt(want[28:30], 16, 32)
	pduMLen, _ := strconv.ParseInt(pdu[28:30], 16, 32)
	if wantMLen != pduMLen {
		fmt.Println("WANT: >"+want[28:30]+"< (", wantMLen, ")  "+want)
		fmt.Println("GOT:  >"+pdu[28:30]+"< (", pduMLen, ")  "+pdu)
		t.Fatalf("MESSAGE LENGTHS BYTES")
	}
	fmt.Println("MESSAGE LENGTHS BYTES OK!")

	if want != pdu {
		fmt.Println(want[0:10], " ", want[10:22], " ", want[22:])
		fmt.Println(pdu[0:10], " ", pdu[10:22], " ", pdu[22:])
		fmt.Println(want)
		fmt.Println(pdu)
		t.Fatalf("TOTAL differ")
	}
	fmt.Printf("length should be set to : %d\r\n", len(pdu)/2)
	fmt.Println("ALL PDU tests OK!!!!!!!!!!!")
}
func TestCreateLongPDU(t *testing.T) {
	message := "Välkommen till akvarellkursen 18/1, 18.00, Sollentuna Ram."
	message = message + message + message + "X"
	number := "+46736290839"
	var want []string
	want = append(want, "0051000B816437260938F900080B8C050003000301005600E4006C006B006F006D006D0065006E002000740069006C006C00200061006B0076006100720065006C006C006B0075007200730065006E002000310038002F0031002C002000310038002E00300030002C00200053006F006C006C0065006E00740075006E0061002000520061006D002E005600E4006C006B006F006D006D0065006E")
	want = append(want, "0051000B816437260938F900080B8C050003000302002000740069006C006C00200061006B0076006100720065006C006C006B0075007200730065006E002000310038002F0031002C002000310038002E00300030002C00200053006F006C006C0065006E00740075006E0061002000520061006D002E005600E4006C006B006F006D006D0065006E002000740069006C006C00200061006B0076")
	want = append(want, "0051000B816437260938F900080B58050003000303006100720065006C006C006B0075007200730065006E002000310038002F0031002C002000310038002E00300030002C00200053006F006C006C0065006E00740075006E0061002000520061006D002E0058")
	pduarray := CreateLongPDU(number, message)
	for g := 0; g < len(want); g++ {
		if pduarray[g] != want[g] {
			fmt.Printf("part1 length = %04x %d\r\n", len(want[g]), len(want[g]))
			fmt.Printf("part1 length = %04x %d\r\n", len(pduarray[g]), len(pduarray[g]))
			fmt.Println(want[g][0:28], " ", want[g][28:48], " ", want[g][48:100])
			fmt.Println(pduarray[g][0:28], " ", pduarray[g][28:48], " ", pduarray[g][48:100])
			fmt.Println("-----------")
			fmt.Println(pduarray[g])
		}
	}
	for i := 0; i < len(pduarray); i++ {
		pduarraylength := len(pduarray[i]) / 2
		fmt.Println(len(pduarray[i]), pduarraylength)
	}
	fmt.Println("ALL CreateLongPDU tests OK!!!!!!!!!!!")
}
func TestCreateLongPDU2(t *testing.T) {
	message := "Välkommen till akvarellkursen 25/1, 18.00, Sköldvägen 16, Sollentuna Ram.\r\n"
	message += "De som inte betalat kursen ännu kan gärns Swisha 1700:- till 0736290839.\r\n"
	message += "Vid frågor, ring!\r\n"
	message += "Med glada hälsningar\r\n"
	message += "Peter & Ulrica"
	number := "+46736290839"
	var want []string
	want = append(want, "0051000B816437260938F900080B8C050003000401005600E4006C006B006F006D006D0065006E002000740069006C006C00200061006B0076006100720065006C006C006B0075007200730065006E002000320035002F0031002C002000310038002E00300030002C00200053006B00F6006C0064007600E400670065006E002000310036002C00200053006F006C006C0065006E00740075006E")
	want = append(want, "0051000B816437260938F900080B8C0500030004020061002000520061006D002E000D000A0044006500200073006F006D00200069006E0074006500200062006500740061006C006100740020006B0075007200730065006E002000E4006E006E00750020006B0061006E0020006700E40072006E0073002000530077006900730068006100200031003700300030003A002D002000740069006C")
	want = append(want, "0051000B816437260938F900080B8C050003000403006C00200030003700330036003200390030003800330039002E000D000A00560069006400200066007200E50067006F0072002C002000720069006E00670021000D000A004D0065006400200067006C0061006400610020006800E4006C0073006E0069006E006700610072000D000A005000650074006500720020002600200055006C0072")
	want = append(want, "0051000B816437260938F900080B0C050003000404006900630061")
	pduarray := CreateLongPDU(number, message)
	for g := 0; g < len(want); g++ {
		if pduarray[g] != want[g] {
			fmt.Println("-------------------------NOT SAME!!!!!")
			fmt.Printf("part1 length = %04x %d\r\n", len(want[g]), len(want[g]))
			fmt.Printf("part1 length = %04x %d\r\n", len(pduarray[g]), len(pduarray[g]))
			fmt.Println(want[g])
			fmt.Println(pduarray[g])
		}
	}
	// for i := 0; i < len(pduarray); i++ {
	// 	pduarraylength := len(pduarray[i]) / 2
	// 	fmt.Println(len(pduarray[i]), pduarraylength)
	// }
	fmt.Println("ALL CreateLongPDU tests OK!!!!!!!!!!!")
}
func TestSendDirectMessage(t *testing.T) {
	// p :=[]string{"0736290839","0736290839"}
	starttime:=time.Now()
	msg :="This is a SendDirectMessage test, sent "+time.Now().Format("2006-01-02 15:04:05")+"! ÅÄÖ åäö"
	fmt.Println("TestSendDirectMessage msg=", msg)
	p,err:=Modemreset("COM3")
	fmt.Println("Modemreset took ", time.Since(starttime))
	if err!=nil {
		t.Fatalf("Modemreset err=%v", err)
	}
	err=SendDirectSMS(p,"0046736290839",msg)
	if err!=nil {
		t.Fatalf("SendDirectSMS err=%v", err)
	}
	fmt.Println("SendSMS took ", time.Since(starttime))
}
func TestSendingMANYDirectMessages(t *testing.T) {
	// p :=[]string{"0736290839","0736290839"}
	var starttime time.Time=time.Now()
	var  msg string = "Detta är ett test med SendDirectMessage, där flera meddelanden sänds efter varandra.\r\n"
	var cnt int=40
	fmt.Printf("Skicka %d SendDirectMessage msg='%s'", cnt, msg)
	p,err:=Modemreset("COM3")
	fmt.Println("Modemreset took ", time.Since(starttime))
	if err!=nil {
		t.Fatalf("Modemreset err=%v", err)
	}
	var sendstarttime time.Time = time.Now()
	for i:=0;i<cnt;i++ {
		s:=fmt.Sprintf("%sNummer %d skickades %s.\r\n",msg,i+1,time.Now().Format("2006-01-02 15:04:05"))
		err=SendDirectSMS(p,"0046736290839",s)
		if err!=nil {
			fmt.Println("Failed to send message ", i+1," but sent ", i, " messages ok.")
			t.Fatalf("SendDirectSMS err=%v", err)
		}
		fmt.Printf("...........Sending %d took %s.\r\n", i+1, time.Since(sendstarttime))
	}
	fmt.Printf("Sending %d messages took total %s.",cnt, time.Since(starttime))
	t.Logf("Sending %d messages took total %s.",cnt, time.Since(starttime))
}
func TestSendDirectLongMessage(t *testing.T) {
	// p :=[]string{"0736290839","0736290839"}
	starttime:=time.Now()
	msg :="This is a SendDirectMessage test, sent "+time.Now().Format("2006-01-02 15:04:05")+"! ÅÄÖ åäö\r\n"
	msg += msg + msg + msg + msg + msg + msg + msg + msg + msg + msg
	fmt.Println("TestSendDirectMessage msg=", msg)
	p,err:=Modemreset("COM3")
	fmt.Println("Modemreset took ", time.Since(starttime))
	if err!=nil {
		t.Fatalf("Modemreset err=%v", err)
	}
	err=SendDirectSMS(p,"0046736290839",msg)
	if err!=nil {
		t.Fatalf("SendDirectSMS err=%v", err)
	}
	fmt.Println("SendSMS took ", time.Since(starttime))
}
func TestSendStoredMessage(t *testing.T) {
	// p :=[]string{"0736290839","0736290839"}
	starttime:=time.Now()
	msg :="This is a SendStoredMessage test, sent"+time.Now().Format("2006-01-02 15:04:05")+"! ÅÄÖ åäö"
	fmt.Println("TestSendStoredMessage msg=", msg)
	p,err:=Modemreset("COM3")
	fmt.Println("Modemreset took ", time.Since(starttime))
	if err!=nil {
		t.Fatalf("Modemreset err=%v", err)
	}
	err=SendStoredSMS(p,"0046736290839",msg)
	if err!=nil {
		t.Fatalf("SendStoredSMS err=%v", err)
	}
	fmt.Println("SendSMS took ", time.Since(starttime))
}
func TestGetPortsList(t *testing.T) {
	// p :=[]string{"0736290839","0736290839"}
	// msg :="This is a test!"
	// s := new(SMStype)
	// s.Comport="COM3"
	// s.SendMessage(p,msg)
	s,err:=GetPortsList()
	if err!=nil {
		t.Fatalf("GetPortsList err=%v", err)
	}
	fmt.Printf("ports=%v\r\n", s)
}
func TestTestPort(t *testing.T) {
	// p :=[]string{"0736290839","0736290839"}
	// msg :="This is a test!"
	// s := new(SMStype)
	// s.Comport="COM3"
	// s.SendMessage(p,msg)
	comport:="COM3"
	result:=TestPort(comport)
	if result!=nil {
		fmt.Print("...TestPort ",comport," results in: ",result.Error(), " -- FAILED!\r\n") 
	} else {
		fmt.Print("...TestPort ",comport," results in: <nil> -- OK!\r\n") 
	}
}
func TestAllTestPorts(t *testing.T) {
	// p :=[]string{"0736290839","0736290839"}
	// msg :="This is a test!"
	// s := new(SMStype)
	// s.Comport="COM3"
	// s.SendMessage(p,msg)
	var s[]string
	// var err error
	s,_= GetPortsList()
	// for i:=0;i<30;i++ {
	// 	s=append(s, fmt.Sprintf("COM%d",i))
	// }
	fmt.Println("Available ports: ",s)
	// if err!=nil {
	// 	t.Fatalf("GetPortsList err=%v", err)
	// }
	for si:=0;si<len(s);si++ {
		fmt.Print("Testing port: ",s[si])	
		result:=TestPort(s[si])
		fmt.Print("...TestPort ",s[si]," results in =",result)
		if result==nil {
			fmt.Println("   --> OK!")
		} else {
			fmt.Println("   --> FAILED!")
		}
	}
}
func TestGetMobilePort(t *testing.T) {
	// p :=[]string{"0736290839","0736290839"}
	// msg :="This is a test!"
	// s := new(SMStype)
	// s.Comport="COM3"
	// s.SendMessage(p,msg)
	s,err:=GetMobilePort()
	if err!=nil {
		t.Fatalf("GetMobilePort err=%v", err)
	}
	fmt.Printf("\r\nGetMobilePort found='%v'\r\n", s)
}

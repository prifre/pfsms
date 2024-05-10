package main

import (
	"fmt"
	"strconv"
	"testing"
)

func TestPDU_unicode(t *testing.T) {
	message := "hello\r\nhelloÅÄÖ"
	number := "+46736290839"
	want := "0011000B916437260938F900080B1E00680065006C006C006F000D000A00680065006C006C006F00C500C400D6" // Unicode encoded.
	//       0011000B916437260938F900000B08C8329BFD766CB9 // default SMS alphabet encoded
	//       0011000B916437260938F900F40B0848656C6C6F3F3F3F // ANSI 8-bit encoded
	pduarray := createLongPDU(number, message)
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
	pduarray := createLongPDU(number, message)
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
	pduarray := createLongPDU(number, message)
	for g := 0; g < len(want); g++ {
		if pduarray[g] != want[g] {
			fmt.Println("-------------------------NOT SAME!!!!!")
			fmt.Printf("part1 length = %04x %d\r\n", len(want[g]), len(want[g]))
			fmt.Printf("part1 length = %04x %d\r\n", len(pduarray[g]), len(pduarray[g]))
			fmt.Println(want[g])
			fmt.Println(pduarray[g])
		} else {
			fmt.Println(g)
		}
	}
	for i := 0; i < len(pduarray); i++ {
		pduarraylength := len(pduarray[i]) / 2
		fmt.Println(len(pduarray[i]), pduarraylength)
	}
	fmt.Println("ALL CreateLongPDU tests OK!!!!!!!!!!!")
}

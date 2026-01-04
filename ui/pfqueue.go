package ui

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

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/prifre/pfsms/pfdatabase"
	"github.com/prifre/pfsms/pfmobile"

	"go.bug.st/serial"
)
type theQueue struct {
	tableQueue *widget.Table
	tableHistory *widget.Table
	dataQueue		[][]string
	dataHistory		[][]string
	btnSubmit      *widget.Button
	btnSubmit2      *widget.Button
	btnDeleteQueue      *widget.Button
	btnDeleteOnefromQueue	*widget.Button
	btnExportHistory	*widget.Button
	Phone    string
	Fname    string
	Lname    string
	Message  string
	logtext  *widget.Label
	mydebug   bool
	Comport   string
	starttime time.Time
	window    fyne.Window
}
func NewQueue( w fyne.Window) *theQueue {
	return &theQueue{ window: w}
}
func (s *theQueue) buildTableQueue() *container.Scroll {
	s.dataQueue=new(pfdatabase.DBtype).ShowQueue()
	s.tableQueue = widget.NewTable(nil,nil,nil)
	s.tableQueue.Length = func() (int, int) {	
		if len(s.dataQueue)<=0 {
			s.dataQueue = [][]string{{"No Queue",""}}
		}
		return len(s.dataQueue),len(s.dataQueue[0])	
	}
	s.tableQueue.CreateCell = func() fyne.CanvasObject {
		return widget.NewLabelWithStyle("0123456789012345",fyne.TextAlignCenter,fyne.TextStyle{Monospace: true})
	}
	s.tableQueue.UpdateCell = func(i widget.TableCellID, o fyne.CanvasObject) {
		o.(*widget.Label).SetText(strings.TrimSpace(s.dataQueue[i.Row][i.Col]))
		o.(*widget.Label).Refresh()
	}
	// s.tableQueue.SetColumnWidth(0,s.window.Content().Size().Width*0.25)	//tstamp
	// s.tableQueue.SetColumnWidth(1,s.window.Content().Size().Width*0.10)	//groupname
	s.tableQueue.SetColumnWidth(0,s.window.Content().Size().Width*0.20)	//phone
	// s.tableQueue.SetColumnWidth(3,s.window.Content().Size().Width*0.45)	//message
	s.tableQueue.Refresh()
	return container.NewScroll(s.tableQueue)
}
func (s *theQueue) buildTableHistory() *container.Scroll {
	s.dataHistory=new(pfdatabase.DBtype).ShowHistory()
	s.tableHistory = widget.NewTable(nil,nil,nil)
	s.tableHistory.Length = func() (int, int) {	
		if len(s.dataHistory)<=0 {
			s.dataHistory = [][]string{{"No History",""}}
		}
		return len(s.dataHistory),len(s.dataHistory[0])	
	}
	s.tableHistory.CreateCell = func() fyne.CanvasObject {
		return widget.NewLabelWithStyle("",fyne.TextAlignLeading,fyne.TextStyle{Monospace: true})
	}
	s.tableHistory.UpdateCell = func(i widget.TableCellID, o fyne.CanvasObject) {
		if len(s.dataHistory)>0 {
			o.(*widget.Label).SetText(strings.TrimSpace(s.dataHistory[i.Row][i.Col]))
			o.(*widget.Label).Refresh()
		}
	}
	s.tableHistory.SetColumnWidth(0,s.window.Content().Size().Width*0.17)	//tstamp
	s.tableHistory.SetColumnWidth(1,s.window.Content().Size().Width*0.1)	//groupname
	s.tableHistory.SetColumnWidth(2,s.window.Content().Size().Width*0.12)	//phone
	s.tableHistory.SetColumnWidth(3,s.window.Content().Size().Width*0.5)	//message
	return container.NewScroll(s.tableHistory)
}
func (s *theQueue) SendQueuedMessages() error {
	// Replace with the correct serial port of the modem
	var comport, tstamp,groupname,sendtext, phoneNumber string
	var success int
	var err error
	var port serial.Port
	comport = fyne.CurrentApp().Preferences().StringWithFallback("mobileport", "COM2")
	addhash := fyne.CurrentApp().Preferences().BoolWithFallback("addhash",false)
	// get phonenumbers & messages to send from queue
	s.starttime = time.Now()
	port,err= pfmobile.Modemreset(comport)
	if err!=nil {
		return errors.New("SendQueuedMessages #1 Modemreset failed: " + err.Error())
	}
	db:=new(pfdatabase.DBtype)
	var id string
	var tot = db.CountinQueue()
	for db.CountinQueue()>0 {
	// pfmobile.Modemcommand(port,"AT+CMMS=2\r","OK",time.Second*2,"Quicksend start",nil)
		rec,err:= db.GetNextinQueue()
		if err!=nil {
			return errors.New("SendQueuedMessages #2 GetNextinQueue failed: " + err.Error())
		}
		// id, tstamp, groupname, phone, message
		id=rec[0]
		tstamp = rec[1]
		groupname=rec[2]
		phoneNumber=rec[3]
		sendtext=rec[4]
		if addhash {
			sendtext = fmt.Sprintf(sendtext+"\r\n#=%d", success+1)
		}
		msg:=sendtext
		if strings.Contains(sendtext,"\\r") {
			msg=strings.ReplaceAll(msg,"\\r","\r")
		}
		if strings.Contains(msg,"\\n") {
			strings.ReplaceAll(msg,"\\n","\n")
		}
		err = pfmobile.SendDirectSMS(port,phoneNumber, msg)
		if err!=nil {
			return errors.New("SendQueuedMessages #3 SendDirectSMS failed to " + phoneNumber + " message: " + fmt.Sprintf("%d", success+1) + ": " + err.Error())
		}
		success++
		err =db.DeletefromQueue(id)
		if err!=nil {
			return errors.New("SendQueuedMessages #4 DeletefromQueue failed for id " + id + ": " + err.Error())
		}
		tstamp = time.Now().Format("20060102150405")
		err = db.SaveHistory([]string{tstamp, groupname, phoneNumber, sendtext})
		if err!=nil {
			return errors.New("SendQueuedMessages #5 SaveHistory failed for phone " + phoneNumber + ": " + err.Error())
		}
		log.Printf("Message %d/%d to phone %s sent! \r\n", success, tot, phoneNumber)
		s.RefreshTables()
	}
	return err
}
func (s *theQueue) SimulateSendMessages() error {
	// Replace with the correct serial port of the modem
	var tstamp, groupname, sendtext, phoneNumber string
	var success int
	var err error
	addhash := fyne.CurrentApp().Preferences().BoolWithFallback("addhash",false)
	// get phonenumbers & messages to send from queue
	s.starttime = time.Now()
	db:=new(pfdatabase.DBtype)
	var id string
	var tot int = db.CountinQueue()
	for db.CountinQueue()>0 {
	// pfmobile.Modemcommand(port,"AT+CMMS=2\r","OK",time.Second*2,"Quicksend start",nil)
		rec,err:= db.GetNextinQueue()
		if err!=nil {
			return errors.New("SendSMS GetNextinQueue failed: " + err.Error())
		}
		// id, tstamp, groupname, phone, message
		id=rec[0]
		tstamp=rec[1]
		groupname =rec[2]
		phoneNumber=rec[3]
		sendtext=rec[4]
		if addhash {
			sendtext = fmt.Sprintf(sendtext+"\r\n#=%d", success+1)
		}
		// err = pfmobile.SendDirectSMS(port,phoneNumber, sendtext)
		// if err!=nil {
		// 	return errors.New("SendSMS SendDirectSMS failed to " + phoneNumber + " message " + fmt.Sprintf("%d", i+1) + ": " + err.Error())
		// }
		success++
		err =db.DeletefromQueue(id)
		if err!=nil {
			return errors.New("SendSMS DeletefromQueue failed for id " + id + ": " + err.Error())
		}
		tstamp = time.Now().Format("20060102150405")
		err = db.SaveHistory([]string{tstamp, groupname, phoneNumber, sendtext})
		if err!=nil {
			return errors.New("SendSMS SaveHistory failed for phone " + phoneNumber + ": " + err.Error())
		}
		s.RefreshTables()
		log.Printf("Simlated Message %d/%d to phone %s sent!\r\n", success, tot, phoneNumber)
		time.Sleep(time.Second*5) // Sleeping here is fine because it's in a goroutine
	}
	return err
}
func (s *theQueue) buildUI() *container.Scroll {
	var err error
    s.btnSubmit = widget.NewButton("Send messages in queue", func() {
		go func() {
			err := s.SendQueuedMessages()
			if err != nil {
				log.Println("Error sending messages from queue: ", err.Error())
			}
		}()
    })
    s.btnSubmit2 = widget.NewButton("Simulate Sending & Updating queue", func() {
		go func() {
			err := s.SimulateSendMessages()
			if err != nil {
				log.Println("Error Simulated sending messages: ", err.Error())
			}
		}()
    })
    s.btnDeleteQueue = widget.NewButton("Clear Whole Queue", func() {
		db:=new(pfdatabase.DBtype)
		err = db.DeleteAllQueue()
        if err != nil {
            log.Println("Error sending messages from queue: ", err.Error())
        }
    })
    s.btnDeleteOnefromQueue = widget.NewButton("Clear One from Queue", func() {
		db:=new(pfdatabase.DBtype)
		rec,_ := db.GetNextinQueue()
		db.DeletefromQueue(rec[0])
		s.RefreshTables()
    })
	s.btnExportHistory = widget.NewButton("Export History", func() {
		new(pfdatabase.DBtype).ExportHistory(fyne.CurrentApp().Preferences().String("historyfile"))
		s.RefreshTables()
})

    // Lägg dem i en Split-container (Queue till vänster, History till höger)
    // split := container.NewHSplit(
	// 	&widget.Card{Title: "Queue", Content: s.buildTableQueue()}, 
	// 	&widget.Card{Title: "History", Content:s.buildTableHistory()})
    // split.Offset = 0.6 // Ger 60% av platsen till kön som standard
	// Skapa vänster sida (Queue)
	queueView := container.NewBorder(NewBoldLabel("Queue"), nil, nil, nil, s.buildTableQueue())

	// Skapa höger sida (History)
	historyView := container.NewBorder(NewBoldLabel("History"), nil, nil, nil, s.buildTableHistory())

	// Lägg in dem i din split
	split := container.NewHSplit(queueView, historyView)
	split.Offset = 0.2
    // return container.NewBorder(
	// 	NewBoldLabel("Queue                           History"),         // Topp
    //     s.btnSubmit, // Botten
    //     nil,         // Vänster
    //     nil,         // Höger
    //     split, // Split-containern fyller nu hela mittenytan
    // )
	return container.NewScroll(container.NewBorder(
		nil,         // Topp
        container.NewHBox(s.btnSubmit,s.btnSubmit2,s.btnDeleteQueue,s.btnDeleteOnefromQueue),	// Botten
        nil,         // Vänster
		nil,	// Höger
		split,
	))
}
func (s *theQueue) RefreshTables() {
	s.dataQueue = new(pfdatabase.DBtype).ShowQueue()
	s.dataHistory = new(pfdatabase.DBtype).ShowHistory()
    // // Säg till den EXISTERANDE tabellen att rita om sig
	s.tableQueue.ScrollToBottom()
	s.tableHistory.ScrollToBottom()
    s.tableQueue.Refresh()
	s.tableHistory.Refresh()
	s.tableQueue.ScrollToBottom()
	s.tableHistory.ScrollToBottom()
 	// s.window.Canvas().Refresh(s.tableQueue)
	// s.window.Canvas().Refresh(s.tableHistory)
}
func (s *theQueue) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Queue", Icon: theme.StorageIcon(), Content: s.buildUI()}
}
package ui

import (
	"fmt"
	"log"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"pfsms/general"
	"pfsms/pfdatabase"
	"pfsms/pfmobile"
)

type theQueue struct {
	tableQueue            *widget.Table
	tableHistory          *widget.Table
	dataQueue             [][]string
	dataHistory           [][]string
	btnSubmit             *widget.Button
	btnSubmit2            *widget.Button
	btnDeleteQueue        *widget.Button
	btnDeleteOnefromQueue *widget.Button
	btnExportHistory      *widget.Button
	CountinQueue          *widget.Label
	logtext               *widget.Entry
	Phone                 string
	Fname                 string
	Lname                 string
	Message               string
	mydebug               bool
	Comport               string
	starttime             time.Time
	window                fyne.Window
}

func NewQueue(w fyne.Window) *theQueue {
	return &theQueue{window: w}
}

func (s *theQueue) buildQueue() fyne.CanvasObject {
	s.logtext = widget.NewMultiLineEntry()
	s.logtext.SetMinRowsVisible(5)
	s.updateLog()

	s.btnSubmit = widget.NewButton("Send messages in queue", func() {
		s.btnSubmit.Disable()
		go func() {
			defer fyne.Do(func() {
				s.btnSubmit.Enable()
				s.RefreshTables()
			})

			err := s.SendQueuedMessages()
			if err != nil {
				log.Println("Error sending messages from queue: ", err.Error())
			}
		}()
	})

	s.btnDeleteQueue = widget.NewButton("Clear Whole Queue", func() {
		db := new(pfdatabase.DBtype)
		err := db.DeleteAllQueue()
		if err != nil {
			log.Println("Error deleting queue: ", err.Error())
		}
		s.RefreshTables()
	})

	s.btnDeleteOnefromQueue = widget.NewButton("Clear One from Queue", func() {
		db := new(pfdatabase.DBtype)
		rec, err := db.GetNextinQueue()
		if err == nil && len(rec) > 0 {
			db.DeletefromQueue(rec[0])
		} else {
			log.Println("Error deleting from queue:", err)
		}
		s.RefreshTables()
	})

	s.btnExportHistory = widget.NewButton("Export History", func() {
		new(pfdatabase.DBtype).ExportHistory(fyne.CurrentApp().Preferences().String("historyfile"))
	})

	queueView := container.NewBorder(general.NewBoldLabel("Queue"), nil, nil, nil, s.buildTableQueue())
	historyView := container.NewBorder(general.NewBoldLabel("History"), nil, nil, nil, s.buildTableHistory())

	split := container.NewHSplit(queueView, historyView)
	split.Offset = 0.25 // 35% bredd till Queue, 65% till History

	bottomControls := container.NewVBox(
		s.logtext,
		container.NewHBox(
			s.btnSubmit,
			s.btnDeleteQueue,
			s.btnDeleteOnefromQueue,
			s.btnExportHistory,
			s.CountinQueue,
		),
	)

	// Använd en ren Border-layout (Mitten fylls av split-vyn, botten av knappar/logg)
	return container.NewBorder(nil, bottomControls, nil, nil, split)
}

func (s *theQueue) buildTableQueue() fyne.CanvasObject {
	s.CountinQueue = widget.NewLabel(fmt.Sprintf("In Queue: %d", new(pfdatabase.DBtype).CountinQueue()))
	s.dataQueue = new(pfdatabase.DBtype).ShowQueue()

	s.tableQueue = widget.NewTable(
		func() (int, int) {
			if len(s.dataQueue) == 0 {
				return 1, 3
			}
			return len(s.dataQueue), 3
		},
		func() fyne.CanvasObject {
			return widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			if len(s.dataQueue) == 0 {
				if i.Col == 0 {
					label.SetText("No Queue")
				} else {
					label.SetText("")
				}
				return
			}
			if i.Row < len(s.dataQueue) && i.Col < len(s.dataQueue[i.Row]) {
				label.SetText(strings.TrimSpace(s.dataQueue[i.Row][i.Col]))
			}
		},
	)

	// Fasta kolumnbredder för Queue (Telefon, Förnamn, Efternamn)
	s.tableQueue.SetColumnWidth(0, 80)
	s.tableQueue.SetColumnWidth(1, 50)
	s.tableQueue.SetColumnWidth(2, 50)

	return s.tableQueue
}

func (s *theQueue) buildTableHistory() fyne.CanvasObject {
	s.dataHistory = new(pfdatabase.DBtype).ShowHistory()

	s.tableHistory = widget.NewTable(
		func() (int, int) {
			if len(s.dataHistory) == 0 {
				return 1, 4
			}
			return len(s.dataHistory), 4
		},
		func() fyne.CanvasObject {
			return widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			if len(s.dataHistory) == 0 {
				if i.Col == 0 {
					label.SetText("No History")
				} else {
					label.SetText("")
				}
				return
			}
			if i.Row < len(s.dataHistory) && i.Col < len(s.dataHistory[i.Row]) {
				label.SetText(strings.TrimSpace(s.dataHistory[i.Row][i.Col]))
			}
		},
	)

	// Fasta kolumnbredder för History (Tidsstämpel, Grupp, Telefon, Meddelande)
	s.tableHistory.SetColumnWidth(0, 175) // Tidsstämpel
	s.tableHistory.SetColumnWidth(1, 70)  // Grupp
	s.tableHistory.SetColumnWidth(2, 120) // Telefon
	s.tableHistory.SetColumnWidth(3, 400) // Meddelande

	return s.tableHistory
}

func (s *theQueue) SendQueuedMessages() error {
	comport := fyne.CurrentApp().Preferences().StringWithFallback("mobileport", "COM2")
	addhash := fyne.CurrentApp().Preferences().BoolWithFallback("addhash", false)

	s.starttime = time.Now()
	port, err := pfmobile.Modemreset(comport)
	if err != nil {
		return fmt.Errorf("SendQueuedMessages #1 Modemreset failed: %w", err)
	}
	defer port.Close()

	db := new(pfdatabase.DBtype)
	var success int
	tot := db.CountinQueue()
	if tot > 0 {
		log.Printf("*************** Started sending messages from Queue, where %d are to be sent!\r\n", tot)
	} else {
		log.Printf("*************** No messages in Queue!\r\n\r\n")
	}
	for db.CountinQueue() > 0 {
		rec, err := db.GetNextinQueue()
		if err != nil {
			return fmt.Errorf("SendQueuedMessages #2 GetNextinQueue failed: %w", err)
		}

		id := rec[0]
		groupname := rec[2]
		phoneNumber := rec[3]
		sendtext := rec[4]

		if addhash {
			sendtext = fmt.Sprintf("%s\r\n#=%d", sendtext, success+1)
		}

		msg := strings.ReplaceAll(sendtext, `\r`, "\r")
		msg = strings.ReplaceAll(msg, `\n`, "\n")

		err = pfmobile.SendDirectSMS(port, phoneNumber, msg)
		if err != nil {
			return fmt.Errorf("SendQueuedMessages #3 SendDirectSMS failed to %s: %w", phoneNumber, err)
		}

		success++
		if err := db.DeletefromQueue(id); err != nil {
			return fmt.Errorf("SendQueuedMessages #4 DeletefromQueue failed for id %s: %w", id, err)
		}

		tstamp := time.Now().Format("20060102150405")
		if err := db.SaveHistory([]string{tstamp, groupname, phoneNumber, sendtext}); err != nil {
			return fmt.Errorf("SendQueuedMessages #5 SaveHistory failed for phone %s: %w", phoneNumber, err)
		}

		log.Printf("Message %d/%d to phone %s sent!\r\n", success, tot, phoneNumber)

		fyne.Do(func() {
			s.RefreshTables()
		})
	}
	fyne.Do(func() {
		s.RefreshTables()
	})
	return nil
}

func (s *theQueue) RefreshTables() {
	s.dataQueue = new(pfdatabase.DBtype).ShowQueue()
	s.dataHistory = new(pfdatabase.DBtype).ShowHistory()

	if s.tableQueue != nil {
		s.tableQueue.ScrollToBottom()
		s.tableQueue.Refresh()
	}
	if s.tableHistory != nil {
		s.tableHistory.ScrollToBottom()
		s.tableHistory.Refresh()
	}

	if s.CountinQueue != nil {
		s.CountinQueue.SetText(fmt.Sprintf("In Queue: %d", new(pfdatabase.DBtype).CountinQueue()))
	}

	s.updateLog()
}

func (s *theQueue) TabItem() *container.TabItem {
	return &container.TabItem{Text: "Queue", Icon: theme.StorageIcon(), Content: s.buildQueue()}
}

func (s *theQueue) updateLog() {
	m, err := general.ReadLastLineWithSeek(fyne.CurrentApp().Preferences().String("pfsmslog"), 5)
	if err != nil {
		log.Println("Read log failed:", err)
		return
	}
	if s.logtext != nil {
		s.logtext.SetText(m)
	}
}

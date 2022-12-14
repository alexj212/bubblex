package bubblex

import (
    "fmt"
    "github.com/charmbracelet/lipgloss"
    "github.com/evertras/bubble-table/table"
    "github.com/potakhov/loge"
    "os"
    "time"
)

type logeHandler struct {
}

func LogeInit() func() {

    c := &logeHandler{}

    logeShutdown := loge.Init(
        loge.Path("."),
        loge.EnableOutputConsole(false),
        loge.EnableOutputFile(false),
        //loge.ConsoleOutput(os.Stdout),
        loge.EnableDebug(),
        loge.EnableError(),
        loge.EnableInfo(),
        loge.EnableWarning(),

        loge.Transports(func(list loge.TransactionList) []loge.Transport {
            transport := loge.WrapTransport(list, c)
            return []loge.Transport{transport}
        }),
    )
    return logeShutdown
}

func AppendEvent(event *LogEvent) {

    row := table.NewRow(table.RowData{ //current_time.Format(time.RubyDate))
        columnKeyTime:    fmt.Sprintf("%s", event.Timestamp.Local().Format("15:04:05.000 EST")),
        columnKeyLevel:   table.NewStyledCell(Name(event.Level), lipgloss.NewStyle().Foreground(lipgloss.Color(Color(event.Level)))),
        columnKeyMessage: event.Message,
        // This isn't a visible column, but we can add the data here anyway for later retrieval
        columnKeyData: event,
    })

    rows = append(rows, row)
    go NotifyAll(rows)
    //m.Update(m.rows)
}

func (m *logeHandler) WriteOutTransaction(tr *loge.Transaction) {

    /*
       type BufferElement struct {
       	Timestamp   time.Time                  `json:"time"`
       	Timestring  [dateTimeStringLength]byte `json:"-"`
       	Message     string                     `json:"msg"`
       	Level       uint32                     `json:"-"`
       	Levelstring string                     `json:"level,omitempty"`
       	Data        map[string]interface{}     `json:"data,omitempty"`

    */

    for i, mesg := range tr.Items {

        //if strings.Contains(mesg.Message, "\n") {
        //    mesg.Message = fmt.Sprintf("@@ %d[%s]", len(mesg.Message), mesg.Message)
        //}
        //if strings.Contains(mesg.Message, "SetDeviceId") {
        //    mesg.Message = fmt.Sprintf("@@ SetDeviceId [%s]", mesg.Message[0:len(mesg.Message)-5])
        //}
        //if strings.Contains(mesg.Message, "ClientControl3") {
        //    mesg.Message = fmt.Sprintf("@@ ClientControl3 [%s]", mesg.Message[0:len(mesg.Message)-5])
        //}

        //if len(mesg.Message) > 20 {
        //    mesg.Message = mesg.Message[0:20]
        //}

        event := &LogEvent{
            EventID:   fmt.Sprintf("%v:%d", tr.ID, i),
            Level:     mesg.Level,
            Message:   mesg.Message,
            Source:    "loge",
            Timestamp: mesg.Timestamp,
            Data:      mesg.Data,
        }

        AppendEvent(event)
    }

}
func (m *logeHandler) FlushTransactions() {

}

func defaultShutdown(sig os.Signal) {
    loge.Printf("caught sig: %v\n\n", sig)
    os.Exit(0)
}

type LogEvent struct {
    EventID   string
    Source    string
    Timestamp time.Time
    Level     uint32
    Message   string
    Data      interface{}
}

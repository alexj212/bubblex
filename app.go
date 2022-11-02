package bubblex

import (
    "github.com/alexj212/gox/commandr"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/gliderlabs/ssh"
)

func NewApp(c *TuiClient) *Ui {
    return NewUi(c)
}

func NewTuiClient(outputHandler func(data []byte), inputHandler func(input string)) *TuiClient {
    console := &TuiClient{}
    console.m = NewUi(console)
    console.m.LogViewer.OnCommandEntered = inputHandler
    console.outputHandler = outputHandler
    console.p = tea.NewProgram(console.m, tea.WithAltScreen())
    RegisterClient(console.p)
    return console
}

type TuiClient struct {
    p             *tea.Program
    s             ssh.Session
    m             *Ui
    user          *sshUser
    outputHandler func(data []byte)
    inputHandler  func(data []byte)
}

//Close interface func implementation to close client down
func (s *TuiClient) Close() {
    if s.s != nil {
        _ = s.s.Close()
    }
}

//ExecLevel interface func implementation to return client exec level
func (s *TuiClient) ExecLevel() commandr.ExecLevel {
    if s.user != nil {
        return s.user.level
    }

    return commandr.SuperAdmin
}

//UserName interface func implementation to return client user name
func (s *TuiClient) UserName() string {
    if s.s != nil {
        return s.s.User()
    }

    return "console"
}

//History interface func implementation to return client command history
func (s *TuiClient) History() []string {
    return s.user.history
}

//Write interface func implementation to write to clients stream
func (s *TuiClient) Write(p []byte) (n int, err error) {
    s.AddReplContent(string(p))
    if s.outputHandler != nil {
        s.outputHandler(p)
    }

    return len(p), nil
}

//WriteString interface func implementation to write string to clients stream
func (s *TuiClient) WriteString(p string) {
    s.AddReplContent(p)
    if s.outputHandler != nil {
        s.outputHandler([]byte(p))
    }
}

func (s *TuiClient) AddReplContent(line string) {
    s.m.r.content = s.m.r.content + line
}
func (s *TuiClient) ClsReplContent() {
    s.m.r.content = ""
}

func (s *TuiClient) Start() error {
    return s.p.Start()
}

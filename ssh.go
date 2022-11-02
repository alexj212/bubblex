package bubblex

import (
    "fmt"
    "github.com/alexj212/gox"
    "github.com/alexj212/gox/commandr"
    "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/wish"
    bm "github.com/charmbracelet/wish/bubbletea"
    "github.com/charmbracelet/wish/logging"
    "github.com/fatih/color"
    "github.com/gliderlabs/ssh"
    "github.com/go-errors/errors"
    "github.com/muesli/termenv"
    "github.com/potakhov/loge"
    gossh "golang.org/x/crypto/ssh"
    "io"
    "log"
    "os/user"
    "path/filepath"
    "time"
)

// ClientDecorator - func def for client initializer
type ClientDecorator func(*sshService)

type TeaAppLauncher func(*TuiClient) *Ui

// SshService interface of ssh service
type SshService interface {
    // Close shut down ssh service
    Close() error
    // Spawn start new go routine serving ssh
    Spawn()
    // RegisterUser register a user on the system
    RegisterUser(user string, level commandr.ExecLevel, keys []*gox.SshKey, history []string)
    // LookupUser lookup a user by name
    LookupUser(username string) (user *sshUser, ok bool)
    // AddCommand add commands to be executed
    AddCommand(cmds ...*commandr.Command)
}

// SshClient interface of ssh client
type SshClient interface {
    //Close interface func implementation to close client down
    Close()
    //UserName interface func implementation to return client user name
    UserName() string
    //ExecLevel interface func implementation to return client exec level
    ExecLevel() commandr.ExecLevel
    //History interface func implementation to return client command history
    History() []string
}
type sshClient struct {
    s    ssh.Session
    user *sshUser
}
type sshUser struct {
    name    string
    level   commandr.ExecLevel
    keys    map[string]bool
    history []string
}

type sshService struct {
    s     *ssh.Server
    users map[string]*sshUser
}

//
//func LaunchSsh(host string, port int, m TeaAppLauncher) {
//    s, err := wish.NewServer(
//        wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
//        wish.WithHostKeyPath(".ssh/term_info_ed25519"),
//        wish.WithMiddleware(
//            BubbleteaProgramMiddleware(m()),
//            logging.Middleware(),
//        ),
//    )
//    if err != nil {
//        log.Fatalln(err)
//    }
//
//    done := make(chan os.Signal, 1)
//    signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
//    loge.Info("Starting SSH server on %s:%d", host, port)
//    go func() {
//        if err = s.ListenAndServe(); err != nil {
//            log.Fatalln(err)
//        }
//    }()
//
//    <-done
//    log.Println("Stopping SSH server")
//    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//    defer func() { cancel() }()
//    if err := s.Shutdown(ctx); err != nil {
//        log.Fatalln(err)
//    }
//}

// NewSshService create a new instance of ssh service
func NewSshService(port int, hostKey gossh.Signer, m TeaAppLauncher, decorators ...ClientDecorator) (SshService, error) {

    svc := &sshService{
        users: make(map[string]*sshUser),
    }

    s, err := wish.NewServer(
        wish.WithAddress(fmt.Sprintf("%s:%d", "0.0.0.0", port)),
        ssh.PublicKeyAuth(svc.publicKeyValidator),
        WithHostKey(hostKey),
        wish.WithMiddleware(
            BubbleteaProgramMiddleware(svc, m),
            logging.Middleware(),
        ),
    )
    if err != nil {
        log.Fatalln(err)
    }

    svc.s = s

    // server.MaxTimeout = 30 * time.Second  // absolute connection timeout, none if empty
    s.IdleTimeout = 60 * time.Second // connection timeout when no activity, none if empty
    s.AddHostKey(hostKey)

    for _, decorator := range decorators {
        decorator(svc)
    }

    log.Printf("starting ssh server on port %s...\n", s.Addr)
    log.Printf("connections will only last %s\n", s.MaxTimeout)
    log.Printf("and timeout after %s of no activity\n", s.IdleTimeout)

    return svc, nil
}

// BubbleteaProgramMiddleware You can write your own custom bubbletea middleware that wraps tea.Program.
// Make sure you set the program input and output to ssh.Session.
func BubbleteaProgramMiddleware(svc SshService, m TeaAppLauncher) wish.Middleware {

    return func(sh ssh.Handler) ssh.Handler {

        return func(s ssh.Session) {
            user, ok := svc.LookupUser(s.User())
            if !ok {
                io.WriteString(s, color.RedString("\nUnknown user: %v\n\n\n", s.User()))
                s.Close()
                return
            }
            loge.Info("user logged in: %v", s.User())

            authorizedKey := gossh.MarshalAuthorizedKey(s.PublicKey())
            loge.Info("public key used by %s    key: %s\n", s.User(), string(authorizedKey))

            var c *TuiClient
            teaHandler := func(s ssh.Session) *tea.Program {
                pty, _, active := s.Pty()
                if !active {
                    loge.Info("no active terminal, skipping")
                    _ = s.Exit(1)
                    return nil
                }

                hpk := s.PublicKey() != nil
                c = &TuiClient{
                    s:    s,
                    user: user,
                }

                app := m(c)
                c.m = app

                loge.Info("%s my connect %s %v %v %s %v %v\n", s.User(), s.RemoteAddr().String(), hpk, s.Command(), pty.Term, pty.Window.Width, pty.Window.Height)
                c.p = tea.NewProgram(app, tea.WithInput(s), tea.WithOutput(s), tea.WithAltScreen())
                RegisterClient(c.p)
                c.p.Send(c)
                return c.p
            }
            val := bm.MiddlewareWithProgramHandler(teaHandler, termenv.ANSI256)
            h := val(sh)
            h(s)
            loge.Info("%s my UnregisterClient %s\n", s.User(), s.RemoteAddr().String())

            if c != nil && c.p != nil {
                UnregisterClient(c.p)
            }

        }
    }
}

// LookupUser lookup a user by name
func (svc *sshService) LookupUser(username string) (user *sshUser, ok bool) {
    user, ok = svc.users[username]
    return user, ok
}

// RegisterUser register a user on the system
func (svc *sshService) RegisterUser(user string, level commandr.ExecLevel, keys []*gox.SshKey, history []string) {

    u := &sshUser{
        name:    user,
        level:   level,
        keys:    make(map[string]bool),
        history: history,
    }

    if u.history == nil {
        u.history = make([]string, 0)
    }

    for _, v := range keys {
        u.keys[string(v.Key)] = true
    }
    svc.users[user] = u
}

// AddCommand add commands to be executed
func (svc *sshService) AddCommand(cmds ...*commandr.Command) {
    commandr.DefaultCommands.AddCommand(cmds...)
}

//Close interface func implementation to close client down
func (s *sshClient) Close() {
    _ = s.s.Close()
}

// Close shut down ssh service
func (svc *sshService) Close() error {
    return svc.s.Close()
}

// Spawn start new go routine serving ssh
func (svc *sshService) Spawn() {
    go svc.s.ListenAndServe()
}

//ExecLevel interface func implementation to return client exec level
func (s *sshClient) ExecLevel() commandr.ExecLevel {
    return s.user.level
}

//UserName interface func implementation to return client user name
func (s *sshClient) UserName() string {
    return s.s.User()
}

//History interface func implementation to return client command history
func (s *sshClient) History() []string {
    return s.user.history
}

//Write interface func implementation to write to clients stream
func (s *sshClient) Write(p []byte) (n int, err error) {
    return s.s.Write(p)
}

//WriteString interface func implementation to write string to clients stream
func (s *sshClient) WriteString(p string) {
    _, _ = s.Write([]byte(p))
}

func (svc *sshService) publicKeyValidator(ctx ssh.Context, key ssh.PublicKey) bool {
    user, ok := svc.LookupUser(ctx.User())
    if !ok {
        fmt.Printf("Login Attempt %v - user not found: %v\n", ctx.RemoteAddr(), ctx.User())
        // user not found
        return false
    }

    pubKey := key.Marshal()
    _, ok = user.keys[string(pubKey)]
    return ok // allow all keys, or use ssh.KeysEqual() to compare against known keys
}

// WithHostKey returns a functional option that adds HostSigners to the server
// from a PEM file as bytes.
func WithHostKey(hostKey gossh.Signer) ssh.Option {
    return func(srv *ssh.Server) error {
        srv.AddHostKey(hostKey)
        return nil
    }
}

func CreatesSshService(port int) (SshService, []*gox.SshKey, error) {

    usr, _ := user.Current()
    sshCertDir := filepath.Join(usr.HomeDir, ".ssh")
    authorizedKeyFile := filepath.Join(sshCertDir, "authorized_keys")

    keys, err := gox.LoadAuthorizedKeys(authorizedKeyFile)
    if err != nil {
        loge.Info("Unable to load authorized keys: %v", err)
        return nil, nil, errors.Errorf("Unable to load authorized keys: %v", err)
    }

    appKey, err := gox.GetAppKey()
    if err != nil {
        loge.Info("Unable to load app key: %v", err)
        return nil, nil, errors.Errorf("Unable to load app key: %v", err)
    }

    hostKey, err := gossh.NewSignerFromKey(appKey)
    if err != nil {
        loge.Info("Unable create hostKey: %v", err)
        return nil, nil, errors.Errorf("Unable create hostKey: %v", err)
    }

    svc, err := NewSshService(port, hostKey, func(c *TuiClient) *Ui {
        return NewApp(c)
    })
    if err != nil {
        loge.Info("Unable to launch ssh server: %v", err)
        return nil, nil, errors.Errorf("Unable to launch ssh server: %v", err)
    }

    return svc, keys, nil

}

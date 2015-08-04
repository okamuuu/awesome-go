package tcptest

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"syscall"
	"testing"
	"time"

	tcptest "github.com/lestrrat/go-tcptest"
)

func Example() {
	var cmd *exec.Cmd
	memd := func(t *tcptest.TCPTest) {
		cmd = exec.Command("memcached", "-p", fmt.Sprintf("%d", t.Port()))
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
		}

		go cmd.Run()
		for loop := true; loop; {
			select {
			case <-t.Done():
				cmd.Process.Kill()
				loop = false
			}
		}
	}

	server, err := tcptest.Start2(memd, 30*time.Second)
	if err != nil {
		log.Fatalf("Failed to start memcached: %s", err)
	}

	log.Printf("memcached started on port %d", server.Port())
	defer func() {
		if cmd != nil && cmd.Process != nil {
			cmd.Process.Signal(syscall.SIGTERM)
		}
	}()

	server.Stop()

	server.Wait()
}

func TestBasic(t *testing.T) {
	cb := func(port int) {
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			t.Fatalf("Failed to listen on port %d: %s", port, err)
		}

		for i := 0; i < 2; i++ {
			conn, err := l.Accept()
			if err != nil {
				t.Fatalf("Failed to accept connection o %d: %s", port, err)
			}
			time.Sleep(500 * time.Millisecond)
			conn.Close()
		}
		l.Close()
	}

	t.Logf("Starting callback")
	server, err := tcptest.Start(cb, time.Minute)
	if err != nil {
		log.Fatalf("Failed to start listening on random port: %s", err)
	}

	t.Logf("Attempting to connect to port %d", server.Port())
	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", server.Port()))
	if err != nil {
		log.Fatalf("Failed to connect to port %d: %s", server.Port(), err)
	}
	defer conn.Close()

	t.Logf("Successfully connected to port %d", server.Port())

	server.Wait()
}

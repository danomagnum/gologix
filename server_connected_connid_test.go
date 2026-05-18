package gologix

import (
	"bytes"
	"net"
	"os"
	"strings"
	"testing"
	"time"
)

// TestConnectedReplyUsesTOConnID is a wire-level regression guard against the
// most subtle bug we hit on this branch: every connected-mode reply
// (sendConnectedReply / sendConnectedError) was passing the
// Originator->Target connection ID to the ConnectionAddress CPF item when
// the spec mandates the Target->Originator ID. Symptom: a real
// ControlLogix MSG Read would receive the TCP bytes but its firmware
// demultiplexer silently dropped the reply, leaving the instruction at
// .EN=True forever. Writes were unaffected because they go through
// sendUnitDataReply (correct ConnID source) and carry no payload.
//
// The test boots the gologix server in-process, performs a connected read
// end-to-end through the standard client, and asserts the response makes
// it back with the correct value — proving the connection layer routes
// the reply to the originator's inbound listener.
func TestConnectedReplyUsesTOConnID(t *testing.T) {
	if probe, err := net.Listen("tcp", "0.0.0.0:44818"); err != nil {
		t.Skipf("port 44818 unavailable: %v", err)
	} else {
		probe.Close()
	}

	router := PathRouter{}
	provider := MapTagProvider{}
	path, err := ParsePath("1,0")
	if err != nil {
		t.Fatalf("parse path: %v", err)
	}
	router.Handle(path.Bytes(), &provider)

	const tag = "regressiontag"
	if err := provider.TagWrite(tag, int32(0x12345678)); err != nil {
		t.Fatalf("seed tag: %v", err)
	}

	srv := NewServer(&router)
	go func() { _ = srv.Serve() }()
	defer func() {
		if srv.TCPListener != nil {
			srv.TCPListener.Close()
		}
		if srv.UDPListener != nil {
			srv.UDPListener.Close()
		}
	}()

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		c, err := net.DialTimeout("tcp", "127.0.0.1:44818", 100*time.Millisecond)
		if err == nil {
			c.Close()
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	client := NewClient("127.0.0.1")
	if err := client.Connect(); err != nil {
		t.Fatalf("client connect: %v", err)
	}
	defer client.Disconnect()

	var got int32
	if err := client.Read(tag, &got); err != nil {
		t.Fatalf("read: %v", err)
	}
	if got != 0x12345678 {
		t.Fatalf("value mismatch: got 0x%08X, want 0x12345678", got)
	}
}

// TestServerConnectedRepliesNeverUseOTConnID is a source-level guard that
// prevents the OT/TO regression from sneaking back in. The wire-level
// behaviour is covered by TestConnectedReplyUsesTOConnID above; this static
// check catches the bug at code-review time. The intent is plain: every
// outbound reply on a connected session must address the connection's
// Target->Originator ID, so the substring `connection.OT` should never
// appear inside server_connected.go.
//
// Legitimate OT references on the REQUEST side (for example, looking up
// inbound connections by their O->T ID via GetByOT) live in other files
// (server_router.go, server.go) and are not in scope here.
func TestServerConnectedRepliesNeverUseOTConnID(t *testing.T) {
	src, err := os.ReadFile("server_connected.go")
	if err != nil {
		t.Fatalf("read server_connected.go: %v", err)
	}
	if i := bytes.Index(src, []byte("connection.OT")); i != -1 {
		// Surface the surrounding line so reviewers see exactly which
		// call site regressed.
		lineStart := bytes.LastIndexByte(src[:i], '\n') + 1
		lineEnd := i + bytes.IndexByte(src[i:], '\n')
		if lineEnd <= i {
			lineEnd = len(src)
		}
		t.Fatalf("server_connected.go re-introduced `connection.OT` in a reply call site — must be `connection.TO`:\n  %s", strings.TrimSpace(string(src[lineStart:lineEnd])))
	}
}

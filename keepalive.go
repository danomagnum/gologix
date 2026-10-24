package gologix

import (
	"errors"
	"log/slog"
	"slices"
	"time"
)

// KeepAlive starts the keepalive goroutine if it is not already running and
// clears any suppression set by a forced KeepAliveCancel. It returns
// immediately; the keepalive runs in the background until the client
// disconnects or KeepAliveCancel is called.
func (client *Client) KeepAlive() {
	client.kaMu.Lock()
	client.kaSuppressed = false
	client.kaMu.Unlock()
	client.startKeepAlive()
}

// KeepAliveCancel stops the keepalive if it is running.
//
// With force == false the call is refused when KeepAliveAutoStart is set, since
// the keepalive would just be restarted on the next Connect(). With force == true
// the keepalive is stopped and will not auto-start again until KeepAlive() is
// called explicitly.
func (client *Client) KeepAliveCancel(force bool) error {
	if client.KeepAliveAutoStart && !force {
		return errors.New("unable to cancel keepalive due to KeepAliveAutoStart == true. use force to cancel anyway")
	}
	if force {
		client.kaMu.Lock()
		client.kaSuppressed = true
		client.kaMu.Unlock()
	}
	client.stopKeepAlive()
	return nil
}

// startKeepAlive launches the keepalive loop unless one is already running or
// a forced cancel has suppressed it. Safe to call repeatedly.
func (client *Client) startKeepAlive() {
	if client.SocketTimeout == 0 {
		return
	}
	client.kaMu.Lock()
	defer client.kaMu.Unlock()
	if client.kaSuppressed || client.kaStop != nil {
		return
	}
	client.Logger.Debug("starting keep alive")
	stop := make(chan struct{})
	client.kaStop = stop
	go client.keepAliveLoop(stop)
}

// stopKeepAlive signals the running keepalive loop (if any) to exit. It does
// not wait for the goroutine to finish so it is safe to call from the loop
// itself via Disconnect().
func (client *Client) stopKeepAlive() {
	client.kaMu.Lock()
	stop := client.kaStop
	client.kaStop = nil
	client.kaMu.Unlock()
	if stop != nil {
		close(stop)
	}
}

// keepAliveLoop periodically reads the configured controller attributes so the
// connection stays open, and optionally refreshes the tag list when a
// controller change is detected. stop is the channel handed out by
// startKeepAlive for this particular run.
func (client *Client) keepAliveLoop(stop <-chan struct{}) {
	defer func() {
		// only clear the handle if it is still ours; a stop/restart may have
		// already replaced it with a new channel.
		client.kaMu.Lock()
		if client.kaStop == stop {
			client.kaStop = nil
		}
		client.kaMu.Unlock()
	}()

	originalProps, err := client.GetAttrList(CipObject_ControllerInfo, 1, client.KeepAliveProps...)
	if err != nil {
		client.Logger.Error(
			"initial keep alive property get failed",
			slog.Any("client.KeepAliveProps", client.KeepAliveProps),
			slog.Any("err", err),
		)
		return
	}

	if client.KeepAlivePollTags {
		err = client.ListAllTags(0)
		if err != nil {
			client.Logger.Error("keepalive list tags failed", slog.Any("err", err))
			return
		}
	}

	t := time.NewTicker(client.KeepAliveFrequency)
	defer t.Stop()
	for {
		select {
		case <-stop:
			return
		case <-t.C:
			if !client.Connected() {
				client.Logger.Warn("keepalive failed. not connected")
				return
			}

			newProps, err := client.GetAttrList(CipObject_ControllerInfo, 1, client.KeepAliveProps...)
			if err != nil {
				client.Logger.Error("keepalive failed", slog.Any("err", err))
				err = client.Disconnect()
				if err != nil {
					client.Logger.Warn("keepalive disconnect failed", slog.Any("err", err))
				}
				return
			}
			if !slices.Equal(newProps.Rest(), originalProps.Rest()) {
				if client.KeepAlivePollTags {
					client.Logger.Debug(
						"controller change detected. re-analyzing types",
						slog.Any("originalProps", originalProps.Rest()),
						slog.Any("newProps", newProps.Rest()),
					)
					err := client.ListAllTags(0)
					if err != nil {
						client.Logger.Warn("keepalive list tags failed.", "error", err)
						return
					}
				} else {
					client.Logger.Debug(
						"controller change detected.",
						slog.Any("originalProps", originalProps.Rest()),
						slog.Any("newProps", newProps.Rest()),
					)
				}
				originalProps = newProps
			}
		}
	}
}

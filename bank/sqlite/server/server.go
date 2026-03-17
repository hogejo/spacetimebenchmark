package main

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

type Action int

const (
	ActionContinue Action = iota
	ActionDisconnect
)

type CommandError interface {
	error
	Action() Action
}

type commandError struct {
	action Action
	msg    string
}

func (e commandError) Error() string  { return e.msg }
func (e commandError) Action() Action { return e.action }

func ContinueError(msg string) CommandError {
	return commandError{action: ActionContinue, msg: msg}
}

func DisconnectError(msg string) CommandError {
	return commandError{action: ActionDisconnect, msg: msg}
}

func runServer(ctx context.Context, config Config, db *sql.DB) {
	ln, err := net.Listen("tcp", config.addr)
	if err != nil {
		log.Fatalf("listen failed: %v", err)
	}
	defer ln.Close()
	log.Printf("listening on %s", config.addr)
	go func() {
		<-ctx.Done()
		log.Printf("shutting down...")
		ln.Close()
	}()
	for {
		conn, err := ln.Accept()
		if err != nil {
			if ctx.Err() != nil {
				log.Printf("server stopped")
				return
			}
			log.Printf("accept failed: %v", err)
			continue
		}
		go handleConn(conn, db)
	}
}

func handleConn(conn net.Conn, db *sql.DB) {
	defer conn.Close()
	remote := conn.RemoteAddr().String()
	log.Printf("client connected: %s", remote)
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("client disconnected: %s", remote)
			} else {
				log.Printf("read error from %s: %v", remote, err)
			}
			return
		}
		msg := strings.TrimRight(line, "\r\n")
		// log.Printf("request from %s: %q", remote, msg)
		parts := strings.Fields(msg)
		if len(parts) == 0 {
			fmt.Fprintf(conn, "ERROR\n")
			log.Printf("invalid request from %s, disconnecting", remote)
			return
		}
		var f handlerFunc
		switch parts[0] {
		case "transfer":
			f = handleTransfer
		case "get":
			f = handleGet
		case "get_total":
			f = handleGetTotal
		case "count_accounts":
			f = handleCountAccounts
		default:
			err := invalidRequest(conn, remote)
			if err.Action() == ActionDisconnect {
				return
			}
			continue
		}
		if err := f(parts, conn, remote, db); err != nil {
			if err.Action() == ActionDisconnect {
				log.Print(err.Error())
				return
			}
			continue
		}
	}
}

type handlerFunc func(parts []string, conn net.Conn, remote string, db *sql.DB) CommandError

func invalidRequest(conn net.Conn, remote string) CommandError {
	if _, werr := fmt.Fprintf(conn, "ERROR\n"); werr != nil {
		return DisconnectError(fmt.Sprintf("write error to %s: %v", remote, werr))
	}
	return DisconnectError(fmt.Sprintf("invalid request from %s, disconnecting", remote))
}

func failedCommand(command string, err error, conn net.Conn, remote string) CommandError {
	if _, werr := fmt.Fprintf(conn, "%s\n", err); werr != nil {
		return DisconnectError(fmt.Sprintf("write error to %s: %v", remote, werr))
	}
	return ContinueError(fmt.Sprintf("%s failed: %v", command, err))
}

func parseTransfer(parts []string) (fromID, toID, amount uint64, ok bool) {
	if len(parts) != 4 {
		return 0, 0, 0, false
	}
	fromID, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return 0, 0, 0, false
	}
	toID, err = strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		return 0, 0, 0, false
	}
	amount, err = strconv.ParseUint(parts[3], 10, 64)
	if err != nil {
		return 0, 0, 0, false
	}
	return fromID, toID, amount, true
}

func handleTransfer(parts []string, conn net.Conn, remote string, db *sql.DB) CommandError {
	fromID, toID, amount, ok := parseTransfer(parts)
	if !ok {
		return invalidRequest(conn, remote)
	}
	if err := transfer(db, fromID, toID, amount); err != nil {
		return failedCommand("transfer", err, conn, remote)
	}
	if _, err := fmt.Fprintf(conn, "OK\n"); err != nil {
		return DisconnectError(fmt.Sprintf("write error to %s: %v", remote, err))
	}
	return nil
}

func parseGet(parts []string) (accountID uint64, ok bool) {
	if len(parts) != 2 {
		return 0, false
	}
	accountID, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return 0, false
	}
	return accountID, true
}

func handleGet(parts []string, conn net.Conn, remote string, db *sql.DB) CommandError {
	accountID, ok := parseGet(parts)
	if !ok {
		return invalidRequest(conn, remote)
	}
	balance, err := getBalance(db, accountID)
	if err != nil {
		return failedCommand("get", err, conn, remote)
	}
	if _, err := fmt.Fprintf(conn, "%d\n", balance); err != nil {
		return DisconnectError(fmt.Sprintf("write error to %s: %v", remote, err))
	}
	return nil
}

func handleGetTotal(parts []string, conn net.Conn, remote string, db *sql.DB) CommandError {
	if len(parts) != 1 {
		return invalidRequest(conn, remote)
	}
	total, err := getTotal(db)
	if err != nil {
		return failedCommand("get_total", err, conn, remote)
	}
	if _, err := fmt.Fprintf(conn, "%d\n", total); err != nil {
		return DisconnectError(fmt.Sprintf("write error to %s: %v", remote, err))
	}
	return nil
}

func parseCountAccounts(parts []string) (fromID, toID uint64, ok bool) {
	if len(parts) != 3 {
		return 0, 0, false
	}
	fromID, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return 0, 0, false
	}
	toID, err = strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		return 0, 0, false
	}
	return fromID, toID, true
}

func handleCountAccounts(parts []string, conn net.Conn, remote string, db *sql.DB) CommandError {
	fromID, toID, ok := parseCountAccounts(parts)
	if !ok {
		return invalidRequest(conn, remote)
	}
	count, err := countAccounts(db, fromID, toID)
	if err != nil {
		return failedCommand("count_accounts", err, conn, remote)
	}
	if _, err := fmt.Fprintf(conn, "%d\n", count); err != nil {
		return DisconnectError(fmt.Sprintf("write error to %s: %v", remote, err))
	}
	return nil
}

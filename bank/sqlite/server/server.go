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
		switch parts[0] {
		case "transfer":
			fromID, toID, amount, ok := parseTransfer(parts)
			if !ok {
				fmt.Fprintf(conn, "ERROR\n")
				log.Printf("invalid request from %s, disconnecting", remote)
				return
			}
			if err := transfer(db, fromID, toID, amount); err != nil {
				if _, werr := fmt.Fprintf(conn, "%s\n", err); werr != nil {
					log.Printf("write error to %s: %v", remote, werr)
					return
				}
				continue
			}
			if _, err := fmt.Fprintf(conn, "OK\n"); err != nil {
				log.Printf("write error to %s: %v", remote, err)
				return
			}
		case "get":
			accountID, ok := parseGet(parts)
			if !ok {
				fmt.Fprintf(conn, "ERROR\n")
				log.Printf("invalid request from %s, disconnecting", remote)
				return
			}
			balance, err := getBalance(db, accountID)
			if err != nil {
				if _, werr := fmt.Fprintf(conn, "%s\n", err); werr != nil {
					log.Printf("write error to %s: %v", remote, werr)
					return
				}
				continue
			}
			if _, err := fmt.Fprintf(conn, "%d\n", balance); err != nil {
				log.Printf("write error to %s: %v", remote, err)
				return
			}
		default:
			fmt.Fprintf(conn, "ERROR\n")
			log.Printf("invalid request from %s, disconnecting", remote)
			return
		}
	}
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

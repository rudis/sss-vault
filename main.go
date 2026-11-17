// sss-vault is a CLI tool that implements Shamir's Secret Sharing (SSS). The
// actual implementation is copied from the shamir/ directory of HashiCorp's
// Vault [1].
//
// [1] https://github.com/hashicorp/vault

// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2026  Simon Ruderich

package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"sss-vault/shamir"
)

func main() {
	log.SetFlags(0)
	usage := func() {
		log.Fatalf("usage: %[1]s split <parts> <threshold>\n"+
			"       %[1]s combine\n", os.Args[0])
	}
	if len(os.Args) < 2 {
		usage()
	}

	if os.Args[1] == "split" && len(os.Args) == 4 {
		parts, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("invalid parts: %v", err)
		}
		threshold, err := strconv.Atoi(os.Args[3])
		if err != nil {
			log.Fatalf("invalid threshold: %v", err)
		}
		_, err = fmt.Fprintln(os.Stderr,
			"Waiting for secret on stdin until EOF")
		if err != nil {
			log.Fatal(err)
		}

		secret, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		res, err := shamir.Split(secret, parts, threshold)
		if err != nil {
			log.Fatal(err)
		}

		_, err = fmt.Printf(
			"\nGenerated %d parts with threshold %d\n\n",
			parts, threshold)
		if err != nil {
			log.Fatal(err)
		}
		for i, x := range res {
			_, err := fmt.Printf("part %d:\n%s\n\n", i+1,
				base64.StdEncoding.EncodeToString(x))
			if err != nil {
				log.Fatal(err)
			}
		}

	} else if os.Args[1] == "combine" && len(os.Args) == 2 {
		_, err := fmt.Fprintln(os.Stderr,
			"Waiting for parts on stdin until EOF (one per line)")
		if err != nil {
			log.Fatal(err)
		}

		stdin, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		var parts [][]byte
		for i, x := range strings.Fields(string(stdin)) {
			dec, err := base64.StdEncoding.DecodeString(x)
			if err != nil {
				log.Fatalf("part %d is invalid base64: %v", i+1, err)
			}
			// Reject part with x=0 which causes Combine to only use this
			// part's value as secret, discarding all other parts. This is
			// only relevant if the part is manipulated.
			if len(dec) > 0 && dec[len(dec)-1] == 0 {
				log.Fatalf("part %d is invalid: x coordinate is 0", i+1)
			}
			parts = append(parts, dec)
		}

		secret, err := shamir.Combine(parts)
		if err != nil {
			log.Fatal(err)
		}
		format := "\nRecovered secret:\n%s\n"
		if !isTerminalSafe(string(secret)) {
			format = "\nRecovered secret (escaped):\n%q\n"
		}
		_, err = fmt.Printf(format, string(secret))
		if err != nil {
			log.Fatal(err)
		}

	} else {
		usage()
	}
}

func isTerminalSafe(x string) bool {
	if !utf8.ValidString(x) {
		return false
	}
	for _, r := range x {
		if r == '\n' || r == '\t' {
			continue
		}
		if unicode.IsControl(r) || unicode.Is(unicode.Bidi_Control, r) {
			return false
		}
	}
	return true
}

// vi: set noet ts=4 sw=4 sts=4:

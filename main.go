// Copyright (c) 2012, Jeramey Crawford <jeramey@antihe.ro>
// Copyright (c) 2013, Jonas mg
// Copyright (c) 2015, Peter Olds <peter@olds.co>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//   * Redistributions of source code must retain the above copyright
//     notice, this list of conditions and the following disclaimer.
//
//   * Redistributions in binary form must reproduce the above copyright
//     notice, this list of conditions and the following disclaimer in
//     the documentation and/or other materials provided with the
//     distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/kless/osutil/user/crypt"
	"github.com/polds/cli"
	"golang.org/x/crypto/ssh/terminal"

	_ "github.com/kless/osutil/user/crypt/md5_crypt"
	_ "github.com/kless/osutil/user/crypt/sha256_crypt"
	_ "github.com/kless/osutil/user/crypt/sha512_crypt"
)

func main() {
	app := cli.NewApp()
	app.Name = "grub-crypt"
	app.Usage = "Encrypt a password."
	app.Version = "0.1"
	app.Copyright = fmt.Sprintf("%s\n", copyright)
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Jeramey Crawford",
			Email: "jeramey@antihe.ro",
		},
		cli.Author{
			Name: "Jonas mg",
		},
		cli.Author{
			Name:  "Peter Olds",
			Email: "peter@olds.co",
		},
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "md5",
			Usage: "Use MD5 to encrypt the password.",
		},
		cli.BoolFlag{
			Name:  "sha-256",
			Usage: "Use SHA-256 to encrypt the password.",
		},
		cli.BoolTFlag{
			Name:  "sha-512",
			Usage: "Use SHA-512 to encrypt the password (default).",
		},
	}

	app.Action = func(c *cli.Context) {
		password := prompt()

		if c.Bool("md5") {
			fmt.Println(gen(crypt.MD5, password))
			return
		}

		if c.Bool("sha-256") {
			fmt.Println(gen(crypt.SHA256, password))
			return
		}

		if c.Bool("sha-512") {
			fmt.Println(gen(crypt.SHA512, password))
			return
		}
	}

	app.Run(os.Args)
}

func prompt() []byte {
	fmt.Print("Password: ")
	password, err := terminal.ReadPassword(0)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}

	fmt.Printf("\nRetype Password: ")
	confirmation, err := terminal.ReadPassword(0)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}

	fmt.Printf("\n")
	if !bytes.Equal(password, confirmation) {
		fmt.Println("Sorry, passwords do not match.")
		os.Exit(1)
	}

	return password
}

func gen(c crypt.Crypt, key []byte) string {
	crypter := crypt.New(c)
	hash, err := crypter.Generate(key, nil)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}

	return hash
}

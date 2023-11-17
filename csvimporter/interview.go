package csvimporter

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

func Import(filename string) map[string]int {
	emailCh := make(chan string, 100)
	hostnameMap := make(map[string]int, 0)

	go readCSV(filename, emailCh)

	for email := range emailCh {
		emailParts := strings.Split(email, "@")
		if len(emailParts) != 2 {
			log.Printf("%q is not an email", email)
			continue
		}

		if err := IsValidHostname(emailParts[1]); err != nil {
			log.Printf("%q is not an email: %s", email, err.Error())
			continue
		}

		hostnameMap[emailParts[1]] += 1
	}

	return hostnameMap
}

func readCSV(filename string, emailCh chan string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer close(emailCh)
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 5
	headerPos := -1

	header, err := reader.Read()
	if err != nil {
		log.Fatal("Failed to parse csv header")
	}

	for idx, headerName := range header {
		if headerName == "email" {
			headerPos = idx
			break
		}
	}

	if headerPos == -1 {
		log.Fatal("csv header does not have email field")
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			return
		}

		if err != nil {
			log.Println(err)
			continue
		}

		emailCh <- record[headerPos]
	}
}

func IsValidHostname(name string) error {
	if name == "" {
		return fmt.Errorf("hostname can't be empty")
	}

	hostname := []byte(name)

	const hostnameMaxLength = 253
	if len(hostname) > hostnameMaxLength {
		return fmt.Errorf("hostname is too long (>%d): %s", hostnameMaxLength, name)
	}

	if hostname[len(hostname)-1] == '.' {
		return fmt.Errorf("hostname ends with trailing dot")
	}

	const hostnameMinLength = 1
	for len(hostname) >= hostnameMinLength {
		var (
			label []byte
			err   error
		)

		label, hostname, err = nextLabel(hostname)
		if err != nil {
			return err
		}

		const labelMinimumLength = 1
		if len(label) < labelMinimumLength {
			return fmt.Errorf("label is too short")
		}

		const labelMaximumLength = 63
		if len(label) > labelMaximumLength {
			return fmt.Errorf("label is too long (>%d): %s", labelMaximumLength, label)
		}

		if label[0] == '-' {
			return fmt.Errorf("label begins with a hyphen: %s", label)
		}

		if label[len(label)-1] == '-' {
			return fmt.Errorf("label ends with a hyphen: %s", label)
		}
	}

	return nil
}

func nextLabel(address []byte) (label, remaining []byte, err error) {
	for i, b := range address {
		if b == '.' {
			return address[:i], address[i+1:], nil
		}

		if !(b >= 'a' && b <= 'z' || b >= '0' && b <= '9' || b == '-' || b >= 'A' && b <= 'Z') {
			c, _ := utf8.DecodeRuneInString(string(address[i:]))
			if c == utf8.RuneError {
				return nil, address, fmt.Errorf("invalid rune at offset %d", i)
			}

			return nil, address, fmt.Errorf("invalid character '%c' at offset %d", c, i)
		}
	}

	return address, nil, nil
}

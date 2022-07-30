package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	sslArgs := os.Args[1:]
	C := sslArgs[0]  // Country
	ST := sslArgs[1] // State
	L := sslArgs[2]  // Locality
	O := sslArgs[3]  // Organization
	OU := sslArgs[4] // Organizational Unit
	CN := sslArgs[5] // Common name (address for cert)

	ssl := `openssl`
	// Certificate request
	csrOpts := fmt.Sprintf(`/C=%s/ST=%s/L=%s/O=%s/OU=%s/CN=%s`, C, ST, L, O, OU, CN)
	csrArgs := []string{
		"req", "-nodes", "-newkey", "rsa:2048", "-keyout", "./dat/cert/" + CN + ".key", "-out", "./dat/cert/" + CN + ".csr", `-subj`, csrOpts,
	}
	csrCmd := exec.Command(ssl, csrArgs...)
	if output, err := csrCmd.CombinedOutput(); err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
	}
	// Certificate sign
	crtArgs := []string{
		"x509", "-req", "-days", fmt.Sprint(7), "-in", "./dat/cert/" + CN + ".csr", "-signkey", "./dat/cert/" + CN + ".key", "-out", "./dat/cert/" + CN + ".crt",
	}
	crtCmd := exec.Command(ssl, crtArgs...)
	if output, err := crtCmd.CombinedOutput(); err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
	}
}

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
	O := sslArgs[3]
	OU := sslArgs[4]
	CN := sslArgs[5] // Common name (address for cert)

	ssl := `openssl`
	// Certificate request
	csrOpts := fmt.Sprintf(`/C=%s/ST=%s/L=%s/O=%s/OU=%s/CN=%s`, C, ST, L, O, OU, CN)
	csrArgs := []string{
		"req", "-nodes", "-newkey", "rsa:2048", "-keyout", "./cert/" + CN + ".key", "-out", "./cert/" + CN + ".csr", `-subj`, csrOpts,
	}
	csrCmd := exec.Command(ssl, csrArgs...)
	if output, err := csrCmd.CombinedOutput(); err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
	}
	// Certificate sign
	crtArgs := []string{
		"x509", "-req", "-days", fmt.Sprint(7), "-in", "./cert/" + CN + ".csr", "-signkey", "./cert/" + CN + ".key", "-out", "./cert/" + CN + ".crt",
	}
	crtCmd := exec.Command(ssl, crtArgs...)
	if output, err := crtCmd.CombinedOutput(); err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
	}
}

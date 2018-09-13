package main

import (
	"fmt"
	"os"
	"crypto/rand"
	"crypto/dsa"
	"crypto/md5"
	"hash"
	"io"
	"math/big"
	pdfcontent "github.com/unidoc/unidoc/pdf/contentstream"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

func generateKeyPair() (*dsa.PrivateKey, dsa.PublicKey){
   
   params := new(dsa.Parameters)
   if err := dsa.GenerateParameters(params, rand.Reader, dsa.L1024N160); err != nil {
      fmt.Println(err)
      os.Exit(1)
   }

   privatekey := new(dsa.PrivateKey)
   privatekey.PublicKey.Parameters = *params
   dsa.GenerateKey(privatekey, rand.Reader) // this generates a public & private key pair

   var pubkey dsa.PublicKey
   pubkey = privatekey.PublicKey

   //fmt.Println("Private Key :")
   //fmt.Printf("%x \n", *privatekey)

   //fmt.Println("Public Key :")
   //fmt.Printf("%x \n",pubkey)
   
	return privatekey, pubkey 
}

func digitallySignPDF(inputText string, privatekey *dsa.PrivateKey) ( []byte , *big.Int, *big.Int) {
   var h hash.Hash
   h = md5.New()
   r := big.NewInt(0)
   s := big.NewInt(0)

   io.WriteString(h, inputText)
   signhash := h.Sum(nil)

   r, s, err := dsa.Sign(rand.Reader, privatekey, signhash)
   if err != nil {
      fmt.Println(err)
   }
		
	return signhash, r, s
}
func verify(signhash []byte, r *big.Int, s *big.Int, pubkey dsa.PublicKey) bool{

   verifystatus := dsa.Verify(&pubkey, signhash, r, s)
   return verifystatus
}
func readPDFContent(inputPath string) string {
	f, err := os.Open(inputPath)
	if err != nil {
		return "error"
	}

	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return "error"
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return "error"
	}

	if isEncrypted {
		_, err = pdfReader.Decrypt([]byte(""))
		if err != nil {
			return "error"
		}
	}

	numPages, err := pdfReader.GetNumPages()
	//if err != nil {
		//return err
	//}

	
	pdfText := ""
	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return "error"
		}

		contentStreams, err := page.GetContentStreams()
		if err != nil {
			return "error"
		}

		pageContentStr := ""
		for _, cstream := range contentStreams {
			pageContentStr += cstream
		}

		cstreamParser := pdfcontent.NewContentStreamParser(pageContentStr)
		txt, err := cstreamParser.ExtractText()
		if err != nil {
			return "error"
		}
		pdfText +=txt
	}

	return pdfText
}

func main() {
	
	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run DigitallySignPDF.go input.pdf\n")
		os.Exit(1)
	}
	inputPath := os.Args[1]
	privKey, pubKey := generateKeyPair()
	inputTxt := readPDFContent(inputPath)
	r, s, hash := digitallySignPDF(inputTxt, privKey)
	isVerified :=verify(r, s, hash, pubKey)
	fmt.Printf("%v",isVerified)
 }

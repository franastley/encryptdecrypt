package main

import (
	"fmt"
	"regexp"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
)
const (
	Base64 string = "^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$"
)

var (
	rxBase64 = regexp.MustCompile(Base64) //Regex to detect if string is Base64 check if string is possibly encrypted.
	bytes := []byte("ABCDEFGHIJKLMNOPQRSQWERW") //We hardcode the key for now however this should be stored somewhere else.
)

// IsBase64 check if a string is base64 encoded.
func IsBase64(str string) bool {
	return rxBase64.MatchString(str)
}

func encrypt(stringToEncrypt string, keyString string) (encryptedString string) {

	//Since the key is in string, we need to convert decode it to bytes
	key, _ := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}

func decrypt(rw http.ResponseWriter, encryptedString string, keyString string) (decryptedString string) {
	//Both the Key and encrypted string are converted to hexadecimals
	key, _ := hex.DecodeString(keyString)
	enc, _ := hex.DecodeString(encryptedString)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	
	//Create a new GCM (Galois/Counter Mode)
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	//Get the nonce size number only used once 
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	return fmt.Sprintf("%s", plaintext)
}

//defines the struct that we will store the json received
type value_struct struct {
	Value string
}


//function for our encrypt endpoint
func encryptvalue(rw http.ResponseWriter, request *http.Request) {
	
	key := hex.EncodeToString(bytes) //encode key in bytes to string and keep as secret, put in a vault
	//decodes the body from request to a json 
	decoder := json.NewDecoder(request.Body)
	
	var t value_struct
	
	//stores the decoded json to our struct, keys must match
	err := decoder.Decode(&t)
	
	//check if there were any issues decoding the json 
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	//check if there is a value from the json 
	if t.Value == "" {
		http.Error(rw, "Json Input Must have a key named Value", http.StatusInternalServerError)
		return
	}
	
	
	//run the encryption function
	encrypted := encrypt(t.Value, key)
	
	//write the encrypted string to the response
	rw.Write([]byte(encrypted))

}

//declare function for out decryption endpoint
func dencryptvalue(rw http.ResponseWriter, request *http.Request) {
	
	key := hex.EncodeToString(bytes) //encode key in bytes to string and keep as secret, put in a vault
	
	decoder := json.NewDecoder(request.Body)
	
	var t value_struct
	err := decoder.Decode(&t)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if t.Value == "" {
		http.Error(rw, "Json Input Must have a key named Value and a value which cannot be blank", http.StatusInternalServerError)
		return
	}
	
	if !IsBase64(t.Value){
		http.Error(rw, "Not a valid encryption string, please use a Base64 string", http.StatusInternalServerError)
		return
	}
	decrypted := decrypt(rw,t.Value, key)
	rw.Write([]byte(decrypted))
	

}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Homepage. You can enjoy our encrypt and decrypt rest api endpoints</h1>")
}


func main() {
	
	r := mux.NewRouter()

	r.HandleFunc("/", homepage)
	
	//Set Mux handle for our encryption endpoint
	r.HandleFunc("/encrypt", encryptvalue).Methods("GET", "OPTIONS").Name("Encrypt")

	//Set Mux handle for our decryption endpoint
	r.HandleFunc("/decrypt", dencryptvalue).Methods("GET", "OPTIONS").Name("Decrypt")


    http.ListenAndServe(":80", r)
}



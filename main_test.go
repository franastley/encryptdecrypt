package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"bytes"
	"strings"
	"github.com/gorilla/mux"
)

type App struct {
    Router *mux.Router
}

var a App

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
    rr := httptest.NewRecorder()
    a.Router.ServeHTTP(rr, req)

    return rr
}

func checkError(err error, t *testing.T) {
	if err != nil {
		t.Errorf("An error occurred. %v", err)
	}
}


func Test_Decrypt(t *testing.T) {
	encrypt := []byte(`{"Value":"672395b5f580812674d75f914dc07eb7a827a0f282b080d21bfaa863769dc5b7"}`)
	
	req, err := http.NewRequest("GET", "/decrypt", bytes.NewBuffer(encrypt))
	checkError(err, t)
	
	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	
	r.HandleFunc("/decrypt", dencryptvalue).Methods("GET")

	r.ServeHTTP(rr, req)
	expectedcode := 200
	if status := rr.Code; status != expectedcode {
		t.Errorf("Status code differs. Expected %d.\n Got %d", http.StatusOK, status)
	}
	
	expected := "test"
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",rr.Body.String(), expected)
    }

}

func Test_DecryptWithoutValue(t *testing.T) {
	encrypt := []byte(`{"Test":"Test"}`)
	
	req, err := http.NewRequest("GET", "/decrypt", bytes.NewBuffer(encrypt))
	checkError(err, t)
	
	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	
	r.HandleFunc("/decrypt", dencryptvalue).Methods("GET")

	r.ServeHTTP(rr, req)
	expectedcode := 500
	if status := rr.Code; status != expectedcode {
		t.Errorf("Status code differs. Expected %d.\n Got %d", http.StatusOK, status)
	}
	
	expected := "Json Input Must have a key named Value and a value which cannot be blank"
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",rr.Body.String(), expected)
    }

}

func Test_DecryptWithInvalidValue(t *testing.T) {
	encrypt := []byte(`{"Value":"This is a Test"}`)
	
	req, err := http.NewRequest("GET", "/decrypt", bytes.NewBuffer(encrypt))
	checkError(err, t)
	
	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	
	r.HandleFunc("/decrypt", dencryptvalue).Methods("GET")

	r.ServeHTTP(rr, req)
	expectedcode := 500
	if status := rr.Code; status != expectedcode {
		t.Errorf("Status code differs. Expected %d.\n Got %d", http.StatusOK, status)
	}
	
	expected := "Not a valid encryption string, please use a Base64 string"
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",rr.Body.String(), expected)
    }

}

func Test_Encrypt(t *testing.T) {
	encrypt := []byte(`{"Value":"test"}`)
	
	req, err := http.NewRequest("GET", "/encrypt", bytes.NewBuffer(encrypt))
	checkError(err, t)
	
	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	
	r.HandleFunc("/encrypt", encryptvalue).Methods("GET")

	r.ServeHTTP(rr, req)
	expectedcode:=500
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d.\n Got %d", http.StatusOK, expectedcode)
	}
	if !IsBase64(strings.TrimSpace(rr.Body.String())){
		t.Errorf("handler returned non base64 body: %v ",rr.Body.String())
	}
	
	
}

func Test_EncryptWithoutValue(t *testing.T) {
	encrypt := []byte(`{"Test":"test"}`)
	
	req, err := http.NewRequest("GET", "/encrypt", bytes.NewBuffer(encrypt))
	checkError(err, t)
	
	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	
	r.HandleFunc("/encrypt", encryptvalue).Methods("GET")
	
	r.ServeHTTP(rr, req)
	expectedcode := 500
	if status := rr.Code; status != expectedcode {
		t.Errorf("Status code differs. Expected %d.\n Got %d", http.StatusOK, expectedcode)
	}
	
	expected := "Json Input Must have a key named Value"
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",rr.Body.String(), expected)
    }

}

func Test_EncryptInvalidJSON(t *testing.T) {
	encrypt := []byte(`"Test"`)
	
	req, err := http.NewRequest("GET", "/encrypt", bytes.NewBuffer(encrypt))
	checkError(err, t)
	
	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	
	r.HandleFunc("/encrypt", encryptvalue).Methods("GET")
	
	r.ServeHTTP(rr, req)
	expectedcode := 500
	if status := rr.Code; status != expectedcode {
		t.Errorf("Status code differs. Expected %d.\n Got %d", http.StatusOK, expectedcode)
	}
	
	expected := "json: cannot unmarshal string into Go value of type main.value_struct"
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",rr.Body.String(), expected)
    }

}

func Test_DecryptInvalidJSON(t *testing.T) {
	encrypt := []byte(`"Test"`)
	
	req, err := http.NewRequest("GET", "/decrypt", bytes.NewBuffer(encrypt))
	checkError(err, t)
	
	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	
	r.HandleFunc("/decrypt", dencryptvalue).Methods("GET")
	
	r.ServeHTTP(rr, req)
	expectedcode := 500
	if status := rr.Code; status != expectedcode {
		t.Errorf("Status code differs. Expected %d.\n Got %d", http.StatusOK, expectedcode)
	}
	
	expected := "json: cannot unmarshal string into Go value of type main.value_struct"
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",rr.Body.String(), expected)
    }

}

func Test_Homepage(t *testing.T) {

	req, err := http.NewRequest("GET", "/", nil)
	checkError(err, t)
	
	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	
	r.HandleFunc("/", homepage).Methods("GET")
	
	r.ServeHTTP(rr, req)
	expectedcode := 200
	if status := rr.Code; status != expectedcode {
		t.Errorf("Status code differs. Expected %d.\n Got %d", http.StatusOK, expectedcode)
	}
	
	expected := "<h1>Homepage. You can enjoy our encrypt and decrypt rest api endpoints</h1>"
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",rr.Body.String(), expected)
    }

}

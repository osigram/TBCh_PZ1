package main

import (
	"PZ1/internal/domain"
	"PZ1/internal/keystorage"
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func main() {
	keyStorage := keystorage.MustNewKeyStorage("key.json")
	key := keyStorage.Key()
	me := domain.Client(key.PublicKey)
	myKey, _ := me.MarshalJSON()
	fmt.Println("Our key: ", string(myKey))

	registerTransaction, err := domain.NewRegisterTransaction(&key)
	if err != nil {
		panic(err)
	}
	if err := sendTransaction(registerTransaction); err != nil {
		fmt.Println("already registered", err)
	}

	for i := 0; i < 2; i++ {
		account, err := getAccount(myKey)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("Your account: %v\n", account)

		client, amount, err := getInputValues()

		transaction, err := domain.NewSendTransaction(&key, *client, int64(amount))
		if err != nil {
			fmt.Println(err)
			continue
		}

		if err = sendTransaction(transaction); err != nil {
			fmt.Println(err)
		}

		fmt.Println("Sent!")
	}
}

func sendTransaction(transaction *domain.Transaction) error {
	data, err := json.Marshal(transaction)
	if err != nil {
		return err
	}

	resp, err := http.Post("http://localhost:34578/transaction",
		"application/json",
		bytes.NewReader(data))
	if err != nil {
		errorString := err.Error()
		if resp != nil {
			errorString = resp.Status + ": " + errorString
		}
		return errors.New(errorString)
	}

	return nil
}

func getAccount(clientData []byte) (int, error) {
	resp, err := http.Get("http://localhost:34578/account?key=" + string(clientData))
	if err != nil {
		errorString := err.Error()
		if resp != nil {
			errorString = resp.Status + ": " + errorString
		}
		return 0, errors.New(errorString)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var account int
	if err := json.Unmarshal(data, &account); err != nil {
		return 0, err
	}

	return account, nil
}

func getInputValues() (client *domain.Client, amount int, err error) {
	client = &domain.Client{}
	fmt.Println("Enter Key of recipient: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return nil, 0, errors.New("unable to get your message")
	}
	clientData := scanner.Text()
	err = client.UnmarshalJSON([]byte(clientData))
	if err != nil {
		return nil, 0, err
	}

	fmt.Println("Enter amount: ")
	scanner = bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		fmt.Printf("unable to get your message")
	}
	amount, err = strconv.Atoi(scanner.Text())
	if err != nil {
		return nil, 0, errors.New("wrong amount")
	}

	return client, amount, nil
}

package json

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

type Person struct {
	Name string
	Age  int
}

func Test(t *testing.T) {
	// Marshal:从Go结构体到JSON
	p1 := Person{Name: "abc", Age: 12}
	jsonData, err := json.Marshal(p1)
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	fmt.Printf("Marshaled JSON: %s cl\n", jsonData)

	// UnMarshal:从JSON到GO结构体
	jsonString := `{"Name":"abc", "Age":12}`
	var p2 Person
	err = json.Unmarshal([]byte(jsonString), &p2)
	if err != nil {
		log.Fatalf("JSON Unmashal failed: %s", err)
	}
	fmt.Printf("Unmarshaled Person: Name=%s, Age=%d\n", p2.Name, p2.Age)
}

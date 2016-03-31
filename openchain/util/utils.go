/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package util

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	gp "google/protobuf"

	"golang.org/x/crypto/sha3"
)

// ComputeCryptoHash should be used in openchain code so that we can change the actual algo used for crypto-hash at one place
func ComputeCryptoHash(data []byte) (hash []byte) {
	hash = make([]byte, 64)
	sha3.ShakeSum256(hash, data)
	return
}

// EncodeToB64 base64-encodes a byte slice
func EncodeToB64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// GenerateUUID returns a UUID based on RFC 4112
func GenerateUUID() string {
	uuid := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, uuid)
	if err != nil {
		panic(fmt.Sprintf("Error generating UUID: %s", err))
	}

	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80

	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40

	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

// CreateUtcTimestamp returns a google/protobuf/Timestamp in UTC
func CreateUtcTimestamp() *gp.Timestamp {
	now := time.Now().UTC()
	secs := now.Unix()
	nanos := int32(now.UnixNano() - (secs * 1000000000))
	return &(gp.Timestamp{Seconds: secs, Nanos: nanos})
}

// LoadFromDisk loads a file from disk
func LoadFromDisk(filename string) (data []byte, err error) {
	data, err = ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Unable to load file %v: %v", filename, err)
	}
	return
}

// SaveToDisk saves a byte slice to disk
func SaveToDisk(filename string, data []byte) error {
	err := ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("Unable to write to file %v: %v", filename, err)
	}
	return nil
}

// EncodeSaveToDisk encodes an object via the gob package and saves it to disk
func EncodeSaveToDisk(filename string, object interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Unable to create file %v: %v", filename, err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(object)
	if err != nil {
		return fmt.Errorf("Unable to encode object before saving to file %v: %v", filename, err)
	}
	return nil
}

// LoadDecodeFromDisk loads a file from disk and decodes it via the gob package
func LoadDecodeFromDisk(filename string, object interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Unable to load file %v: %v", filename, err)
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(object)
	if err != nil {
		return fmt.Errorf("Unable to decode loaded file %v: %v", filename, err)
	}
	return nil
}

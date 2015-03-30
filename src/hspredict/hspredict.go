/*
Tool for predicting the HSDirs responsible for a particular hidden service at a
given time using the rendevous v2 scheme as specified in rend-spec.txt. Heavily
influenced by Tor source and Donncha O'Cearbhaill's retrieve_hs_descriptor.py

author: George Tankersley <george.tankersley@gmail.com>
*/

package hspredict

import (
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
)

const (
	REPLICAS                          = 2
	REND_TIME_PERIOD_V2_DESC_VALIDITY = 24 * 60 * 60 // 86400
)

func ComputeRendV2DescId(service_id string, replica byte, time int64, desc_cookie string) (string, error) {
	var err error

	/* Convert service ID to binary. */
	service_id_binary, err := base32.StdEncoding.DecodeString(service_id)

	if err != nil {
		fmt.Println("err: ", err)
		return "", err
	}

	/* Calculate current time-period. */
	time_period := get_time_period(time, 0, service_id_binary)

	/* Calculate secret-id-part = h(time-period | cookie | replica). */
	secret_id_part := get_secret_id_part_bytes(time_period, desc_cookie, replica)

	/* Calculate descriptor ID. */
	desc_id := rend_get_descriptor_id_bytes(service_id_binary, secret_id_part)

	desc_encode := string(base32.StdEncoding.EncodeToString(desc_id))

	return desc_encode, nil
}

func get_time_period(time int64, deviation int64, service_id_binary []byte) int64 {
	return (time+int64(service_id_binary[0])*REND_TIME_PERIOD_V2_DESC_VALIDITY/256)/REND_TIME_PERIOD_V2_DESC_VALIDITY + int64(deviation)
}

func get_secret_id_part_bytes(time_period int64, desc_cookie string, replica byte) []byte {
	h := sha1.New()
	htonl_time := make([]byte, 4)
	binary.BigEndian.PutUint32(htonl_time, uint32(time_period))
	h.Write(htonl_time)
	if desc_cookie != "" {
		h.Write([]byte(desc_cookie))
	}
	h.Write([]byte{replica})
	return h.Sum(nil)
}

func rend_get_descriptor_id_bytes(service_id_binary, secret_id_part []byte) []byte {
	h := sha1.New()
	h.Write(service_id_binary)
	h.Write(secret_id_part)
	return h.Sum(nil)
}

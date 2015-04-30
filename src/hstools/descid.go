// Tool for predicting the HSDirs responsible for a particular hidden service at a
// given time using the rendevous v2 scheme as specified in rend-spec.txt. Heavily
// influenced by Tor source and Donncha O'Cearbhaill's retrieve_hs_descriptor.py
//
// author: George Tankersley <george.tankersley@gmail.com>
// author: Filippo Valsorda <hi@filippo.io>

package hstools

import (
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"strings"
	"time"
)

const (
	REPLICAS                          = 2
	REND_TIME_PERIOD_V2_DESC_VALIDITY = 24 * 60 * 60 // 86400
)

func OnionToDescID(onion string, t time.Time) ([][]byte, error) {
	onion = strings.ToLower(onion)

	switch len(onion) {
	case 16 + len(".onion"):
		if onion[16:] == ".onion" {
			onion = onion[:16]
		} else {
			return nil, errors.New("wrong suffix")
		}
	case 16:
		// good like this
	default:
		return nil, errors.New("wrong length")
	}

	first, err := ComputeRendV2DescID(onion, 0, t.Unix(), "")
	if err != nil {
		return nil, err
	}
	second, err := ComputeRendV2DescID(onion, 1, t.Unix(), "")
	if err != nil {
		return nil, err
	}

	return [][]byte{first, second}, nil
}

func ComputeRendV2DescID(serviceID string, replica byte, time int64, descCookie string) ([]byte, error) {
	// Convert service ID to binary.
	serviceIDBinary, err := FromBase32(serviceID)
	if err != nil {
		return nil, err
	}

	// Calculate current time-period.
	timePeriod := getTimePeriod(time, 0, serviceIDBinary)

	// Calculate secret-id-part = h(time-period | cookie | replica).
	secretIDPart := getSecretIDPartBytes(timePeriod, descCookie, replica)

	// Calculate descriptor ID.
	descID := rendGetDescriptorIDBytes(serviceIDBinary, secretIDPart)

	return descID, nil
}

func getTimePeriod(time int64, deviation int64, serviceIDBinary []byte) int64 {
	return (time+int64(serviceIDBinary[0])*REND_TIME_PERIOD_V2_DESC_VALIDITY/256)/REND_TIME_PERIOD_V2_DESC_VALIDITY + int64(deviation)
}

func getSecretIDPartBytes(timePeriod int64, descCookie string, replica byte) []byte {
	h := sha1.New()
	htonlTime := make([]byte, 4)
	binary.BigEndian.PutUint32(htonlTime, uint32(timePeriod))
	h.Write(htonlTime)
	if descCookie != "" {
		h.Write([]byte(descCookie))
	}
	h.Write([]byte{replica})
	return h.Sum(nil)
}

func rendGetDescriptorIDBytes(serviceIDBinary, secretIDPart []byte) []byte {
	h := sha1.New()
	h.Write(serviceIDBinary)
	h.Write(secretIDPart)
	return h.Sum(nil)
}

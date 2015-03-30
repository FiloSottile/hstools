package hspredict

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"testing"
	"time"
)

func mkServiceId() string {
	rand_id := make([]byte, 20)
	rand.Read(rand_id)
	hash := sha1.Sum(rand_id)
	service_id := base32.StdEncoding.EncodeToString(hash[:10])
	return service_id
}

func TestRendComputeV2DescId(t *testing.T) {
	service_id := "FACEBOOKCOREWWWI" //mkServiceId()
	if len(service_id) != 16 {
		t.Errorf("generated invalid service id")
	}
	time := time.Now().Unix()
	desc_id1, _ := ComputeRendV2DescId(service_id, 0, time, "")
	desc_id2, _ := ComputeRendV2DescId(service_id, 1, time, "")
	fmt.Println("desc_id: ", desc_id1, desc_id2)
}

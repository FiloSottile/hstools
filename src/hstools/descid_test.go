package hstools

import (
	"testing"
	"time"
)

func TestFacebookOnion(t *testing.T) {
	tt, _ := time.Parse(time.RFC3339, "2015-04-11T19:30:00Z")
	desc, err := OnionToDescID("facebookcorewwwi.onion", tt)
	if err != nil {
		t.Fatal(err)
	}
	if desc[0] != "e4jiuabozanwqxdobx44w47mx2hi2auz" {
		t.Errorf("Wrong desc[0]: %v (!= e4jiuabozanwqxdobx44w47mx2hi2auz)", desc[0])
	}
	if desc[1] != "tyvtyaqd4trmgoopqktv4aawelu6skes" {
		t.Errorf("Wrong desc[0]: %v (!= tyvtyaqd4trmgoopqktv4aawelu6skes)", desc[1])
	}
}

func TestCurrentOnion(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	desc, err := OnionToDescID("facebookcorewwwi.onion", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(desc[0], desc[1])
}

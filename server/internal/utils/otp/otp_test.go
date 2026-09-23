package otp

import (
	"regexp"
	"testing"
)

func TestGenerateOTP(t *testing.T) {
	for range 100 {
		otp, err := GenerateOTP()
		if err != nil {
			t.Fatalf("failed to generate otp: %v", err)
		}

		if len(otp) != 6 {
			t.Errorf("expected otp length 6 got %v value(%v)", len(otp), otp)
		}
		if matched, _ := regexp.MatchString(`^[0-9]{6}`, otp); !matched {
			t.Errorf("otp %v contains non numeric characters", otp)
		}
	}
}

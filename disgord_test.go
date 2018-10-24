package disgord

import "testing"

func TestValidateUsername(t *testing.T) {
	var err error

	if err = ValidateUsername(""); err == nil {
		t.Error("expected empty error")
	}

	if err = ValidateUsername("a"); err == nil {
		t.Error("expected username to be too short")
	}

	if err = ValidateUsername("gk523526hdfgdfjdghlkjdhfglksjhdfg"); err == nil {
		t.Error("expected username to be too long")
	}

	if err = ValidateUsername("  anders"); err == nil {
		t.Error("expected username to have whitespace prefix error")
	}

	if err = ValidateUsername("anders  "); err == nil {
		t.Error("expected username to have whitespace suffix error")
	}

	if err = ValidateUsername("and  ers"); err == nil {
		t.Error("expected username to have excessive whitespaces error")
	}

	if err = ValidateUsername("@anders"); err == nil {
		t.Error("expected illegal char error")
	}

	if err = ValidateUsername("#anders"); err == nil {
		t.Error("expected illegal char error")
	}

	if err = ValidateUsername("and:ers"); err == nil {
		t.Error("expected illegal char error")
	}

	if err = ValidateUsername("and```ers"); err == nil {
		t.Error("expected illegal char error")
	}

	if err = ValidateUsername("discordtag"); err == nil {
		t.Error("expected illegal username error")
	}

	if err = ValidateUsername("everyone"); err == nil {
		t.Error("expected illegal username error")
	}

	if err = ValidateUsername("here"); err == nil {
		t.Error("expected illegal username error")
	}
}

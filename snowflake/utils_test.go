package snowflake

import "testing"

func TestParseSnowflakeString(t *testing.T) {
	// test panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ParseID did panic")
		}
	}()

	id := "435639843545"
	if ParseSnowflakeString(id).String() != id {
		t.Errorf("Incorrect string parsing for ID, base 10. Wants %s, got %s", id, ParseSnowflakeString(id).String())
	}
}

func TestParseSnowflakeStringWithPanicTriggerLetters(t *testing.T) {
	// test panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("ParseID did not panic")
		}
	}()

	id := "435639sd843545gf453s"
	if ParseSnowflakeString(id).String() != id {
		t.Errorf("Incorrect string parsing for ID, base 10. Wants %s, got %s", id, ParseSnowflakeString(id).String())
	}
}

func TestParseSnowflakeStringWithPanicTriggerOverflow(t *testing.T) {
	// test panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("ParseID did not panic")
		}
	}()

	id := "184467440737095516151" // string(uint64(0) - 1) + "1"
	if ParseSnowflakeString(id).String() != id {
		t.Errorf("Incorrect string parsing for ID, base 10. Wants %s, got %s", id, ParseSnowflakeString(id).String())
	}
}

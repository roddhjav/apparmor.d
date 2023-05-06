// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package logging

import "testing"

func TestPrint(t *testing.T) {
	msg := "Print message"
	wantN := 13

	gotN := Print(msg)
	if gotN != wantN {
		t.Errorf("Print() = %v, want %v", gotN, wantN)
	}
}

func TestPrintln(t *testing.T) {
	msg := "Print message"
	wantN := 14
	gotN := Println(msg)
	if gotN != wantN {
		t.Errorf("Println() = %v, want %v", gotN, wantN)
	}
}

func TestBulletf(t *testing.T) {
	msg := "Bullet message"
	want := "\033[1m ⋅ \033[0mBullet message\n"
	if got := Bulletf(msg); got != want {
		t.Errorf("Bulletf() = %v, want %v", got, want)
	}
}

func TestBullet(t *testing.T) {
	msg := "Bullet message"
	wantN := 28
	gotN := Bullet(msg)
	if gotN != wantN {
		t.Errorf("Bullet() = %v, want %v", gotN, wantN)
	}
}

func TestStepf(t *testing.T) {
	msg := "Step message"
	want := "\033[1;32mStep message\033[0m\n"
	if got := Stepf(msg); got != want {
		t.Errorf("Stepf() = %v, want %v", got, want)
	}
}

func TestStep(t *testing.T) {
	msg := "Step message"
	wantN := 24
	gotN := Step(msg)
	if gotN != wantN {
		t.Errorf("Step() = %v, want %v", gotN, wantN)
	}
}

func TestSuccessf(t *testing.T) {
	msg := "Success message"
	want := "\033[1;32m ✓ \033[0mSuccess message\n"
	if got := Successf(msg); got != want {
		t.Errorf("Successf() = %v, want %v", got, want)
	}
}

func TestSuccess(t *testing.T) {
	msg := "Success message"
	wantN := 32
	gotN := Success(msg)
	if gotN != wantN {
		t.Errorf("Success() = %v, want %v", gotN, wantN)
	}
}

func TestWarningf(t *testing.T) {
	msg := "Warning message"
	want := "\033[1;33m ‼ \033[0mWarning message\n"
	if got := Warningf(msg); got != want {
		t.Errorf("Warningf() = %v, want %v", got, want)
	}
}

func TestWarning(t *testing.T) {
	msg := "Warning message"
	wantN := 32
	gotN := Warning(msg)
	if gotN != wantN {
		t.Errorf("Warning() = %v, want %v", gotN, wantN)
	}
}

func TestError(t *testing.T) {
	msg := "Error message"
	wantN := 30
	gotN := Error(msg)
	if gotN != wantN {
		t.Errorf("Error() = %v, want %v", gotN, wantN)
	}
}

func TestFatalf(t *testing.T) {
	msg := "Error message"
	want := "\033[1;31m ✗ Error: \033[0mError message\n"
	if got := Fatalf(msg); got != want {
		t.Errorf("Fatalf() = %v, want %v", got, want)
	}
}

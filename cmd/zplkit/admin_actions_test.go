package main

import (
	"regexp"
	"testing"
)

func TestBuildClockZPL(t *testing.T) {
	zpl, err := buildClockZPL(nil)
	if err != nil {
		t.Fatalf("buildClockZPL(nil) returned error: %v", err)
	}
	want := regexp.MustCompile(`^\^XA\^ST\d{1,2},\d{1,2},\d{4},\d{1,2},\d{1,2},\d{1,2},M\^XZ$`)
	if !want.MatchString(zpl) {
		t.Errorf("buildClockZPL(nil) = %q, want match of %s", zpl, want.String())
	}

	if _, err := buildClockZPL([]string{"unexpected"}); err == nil {
		t.Error("buildClockZPL with an argument: want error, got nil")
	}
}

func TestBuildDarknessZPL(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr bool
	}{
		{"min", []string{"0"}, "~SD00^XA^JUS^XZ", false},
		{"typical", []string{"15"}, "~SD15^XA^JUS^XZ", false},
		{"max", []string{"30"}, "~SD30^XA^JUS^XZ", false},
		{"below range", []string{"-1"}, "", true},
		{"above range", []string{"31"}, "", true},
		{"not a number", []string{"abc"}, "", true},
		{"no arguments", nil, "", true},
		{"too many arguments", []string{"15", "20"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildDarknessZPL(tt.args)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("buildDarknessZPL(%v): want error, got nil (result %q)", tt.args, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("buildDarknessZPL(%v) returned error: %v", tt.args, err)
			}
			if got != tt.want {
				t.Errorf("buildDarknessZPL(%v) = %q, want %q", tt.args, got, tt.want)
			}
		})
	}
}

func TestBuildCalibrateZPL(t *testing.T) {
	got, err := buildCalibrateZPL(nil)
	if err != nil {
		t.Fatalf("buildCalibrateZPL(nil) returned error: %v", err)
	}
	if want := "~JC"; got != want {
		t.Errorf("buildCalibrateZPL(nil) = %q, want %q", got, want)
	}

	if _, err := buildCalibrateZPL([]string{"unexpected"}); err == nil {
		t.Error("buildCalibrateZPL with an argument: want error, got nil")
	}
}

func TestBuildSpeedZPL(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr bool
	}{
		{"min", []string{"1"}, "^XA^PR1^JUS^XZ", false},
		{"typical", []string{"4"}, "^XA^PR4^JUS^XZ", false},
		{"max", []string{"14"}, "^XA^PR14^JUS^XZ", false},
		{"below range", []string{"0"}, "", true},
		{"above range", []string{"15"}, "", true},
		{"not a number", []string{"abc"}, "", true},
		{"no arguments", nil, "", true},
		{"too many arguments", []string{"4", "4"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildSpeedZPL(tt.args)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("buildSpeedZPL(%v): want error, got nil (result %q)", tt.args, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("buildSpeedZPL(%v) returned error: %v", tt.args, err)
			}
			if got != tt.want {
				t.Errorf("buildSpeedZPL(%v) = %q, want %q", tt.args, got, tt.want)
			}
		})
	}
}

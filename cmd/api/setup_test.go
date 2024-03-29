package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	appInit()
	//  Make this work, all you have to do is
	// M.Run will handle running tests after app setup
	os.Exit(m.Run())
}

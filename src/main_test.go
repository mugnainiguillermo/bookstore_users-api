package src

import (
	"fmt"
	"testing"
	"time"
)

func Test(t *testing.T) {
	now := time.Now().UTC()
	fmt.Println(now.Format(time.RFC3339))
}

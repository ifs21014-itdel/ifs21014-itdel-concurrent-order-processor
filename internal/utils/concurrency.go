package utils

import "fmt"

func SafeGoroutine(name string, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("❌ Recovered in:", name, "error:", r)
		} else {
			fmt.Println("✅ Goroutine", name, "done")
		}
	}()
	fmt.Println("🚀 Start goroutine:", name)
	fn()
}

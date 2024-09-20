package utils

import (
	"fmt"
	"runtime"
	"time"
)

func PrintMemoryUsage(intervalSeconds int) {
	ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)

		// Conversion des valeurs en MB pour plus de clarté
		fmt.Printf("Mémoire allouée: %.2f MB\n", float64(memStats.Alloc)/1024/1024)
		fmt.Printf("Mémoire totale allouée: %.2f MB\n", float64(memStats.TotalAlloc)/1024/1024)
		fmt.Printf("Mémoire système obtenue: %.2f MB\n", float64(memStats.Sys)/1024/1024)
		fmt.Printf("Nombre de Garbage Collection : %d\n", memStats.NumGC)
		fmt.Println("-----------------------------------")
	}
}

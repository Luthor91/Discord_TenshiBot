package utils

import (
	"fmt"
	"strconv"
	"time"
)

// parseDuration parse la durée au format '10s', '5m', '2h', '1d'
func ParseDuration(durationStr string) (time.Duration, error) {
	if len(durationStr) < 2 {
		return 0, fmt.Errorf("durée invalide")
	}

	unit := durationStr[len(durationStr)-1]
	valueStr := durationStr[:len(durationStr)-1]
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("valeur numérique invalide")
	}

	switch unit {
	case 's': // secondes
		return time.Duration(value) * time.Second, nil
	case 'm': // minutes
		return time.Duration(value) * time.Minute, nil
	case 'h': // heures
		return time.Duration(value) * time.Hour, nil
	case 'd': // jours
		return time.Duration(value) * time.Hour * 24, nil
	default:
		return 0, fmt.Errorf("unité de temps invalide")
	}
}

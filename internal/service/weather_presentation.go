package service

func mapWeatherCode(code int) string {
	switch code {
	case 0:
		return "Ясно"
	case 1, 2, 3:
		return "Переменная облачность"
	case 45, 48:
		return "Туман"
	case 51, 53, 55:
		return "Морось"
	case 61, 63, 65:
		return "Дождь"
	case 71, 73, 75:
		return "Снег"
	case 95:
		return "Гроза"
	default:
		return "Неизвестно"
	}
}

func getClothing(temp float64) string {
	if temp < coldThreshold {
		return "Тёплая одежда"
	}
	if temp < warmThreshold {
		return "Куртка"
	}
	return "Лёгкая одежда"
}

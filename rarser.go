package rarser

func GetRequestParser() RequestParser {
	once.Do(func() {
		instance = &requestParser{}
	})
	return instance
}
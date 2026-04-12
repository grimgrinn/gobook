package wait

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// WaitForServer пытается соедениться с сервером заданного URL.
// Попытки предпринимаются в течении минуты с растущими интервалами.
// Сообщает об ошибке, если все попытки неудачны.
func WaitForServer(url string) error {
	const timeout = 1 * time.Minute
	deadline := time.Now().Add(timeout)
	for tries := 0; time.Now().Before(deadline); tries++ {
		_, err := http.Head(url)
		if err == nil {
			return nil // Успешное соеденение
		}
		log.Printf("Сервер не отвечает (%s); повтор...", err)
		time.Sleep(time.Second << uint(tries)) // Увеличение задержки
	}
	return fmt.Errorf("Сервер %s не отвечает; время %s ", url, timeout)
}

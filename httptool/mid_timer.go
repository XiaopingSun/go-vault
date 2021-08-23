package httptool

import (
	"fmt"
	"net/http"
	"time"
)

type Mid_timer struct {}

func (m *Mid_timer)mid_handle(next http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {

		timeStart := time.Now()

		next.ServeHTTP(w, r)

		timeElapsed := time.Since(timeStart)
		fmt.Println("Time Elapsed:", timeElapsed)
		fmt.Println("Timer Middle Ware Work Done.")
	}
	return http.HandlerFunc(handler)
}

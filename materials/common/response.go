package common

import "net/http"

type ResponseJson struct {
	StatusCode int         `json:"statusCode"`
	Data       interface{} `json:"data"`
	Status     int         `json:"status"`
	Message    string      `json:"message"`
}

func FormatResult(list interface{}, postFlg bool, err error) (int, ResponseJson) {

	// fmt.Printf("the list is", list)
	// fmt.Println("...................")
	hs := 500
	if postFlg {
		hs = http.StatusCreated
	} else {
		hs = http.StatusOK
	}

	rj := ResponseJson{}
	if err != nil {
		hs = http.StatusInternalServerError
		rj.StatusCode = http.StatusInternalServerError
		rj.Message = err.Error()
	} else {
		rj.StatusCode = hs
		rj.Data = list
		rj.Message = "done"
	}

	// fmt.Printf("the final list is", rj)
	// fmt.Println(".............")
	return hs, rj
}

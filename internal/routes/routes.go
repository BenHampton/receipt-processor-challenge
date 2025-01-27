package routes

import (
	"encoding/json"
	"fetch-interview/internal/model"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var receiptMap = make(map[uuid.UUID]int)

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/receipts/process", create).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", findById)

	return router
}

func findById(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key := vars["id"]

	uuidErr := uuid.Validate(key)

	if uuidErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := uuid.MustParse(key)

	receiptMapValue, ok := receiptMap[id]

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var response = model.Point{Points: receiptMapValue}

	json.NewEncoder(w).Encode(response)
}

func create(w http.ResponseWriter, r *http.Request) {

	requestBody, error := io.ReadAll(r.Body)

	if error != nil || len(requestBody) == 0 {
		http.Error(w, "Error reading requestBody", http.StatusBadRequest)
		return
	}

	var receipt model.Receipt
	err := json.Unmarshal(requestBody, &receipt)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response := createReceiptsPoints(receipt)

	json.NewEncoder(w).Encode(response)
}

func createReceiptsPoints(receipt model.Receipt) model.ReceiptsPoints {

	var totalPoints = 0

	totalPoints += getRetailPoints(receipt.Retailer)
	totalPoints += getPointsFromTotal(receipt.Total)
	totalPoints += getPointsFromItems(receipt.Items)
	totalPoints += getPointsFromDateTime(receipt.PurchaseDate, receipt.PurchaseTime)

	fmt.Printf("totalPoints: %d\n\n", totalPoints)

	var receiptsPoints = model.ReceiptsPoints{Id: uuid.New()}

	receiptMap[receiptsPoints.Id] = totalPoints

	return receiptsPoints
}

func getRetailPoints(name string) int {

	if len(name) == 0 {
		return 0
	}

	var points = 0
	for _, ch := range name {

		if unicode.IsLetter(ch) || unicode.IsNumber(ch) {
			points++
		}
	}

	fmt.Printf("%d points - retailer name \"%s\" has %d characters\n", points, name, points)

	return points
}

func getPointsFromTotal(totalStr string) int {

	if len(totalStr) == 0 {
		return 0
	}

	var points = 0
	totalFloat, err := strconv.ParseFloat(totalStr, 64)

	if err == nil && math.Floor(totalFloat) == totalFloat {

		points += 50
		fmt.Printf("50 points - total is a round dollar amount\n")

		if isMultipleOfFour(totalFloat) {
			points += 25
			fmt.Printf("25 points - total is a multiple of 0.25\n")
		}
	}

	return points
}

func getPointsFromItems(items []model.Item) int {

	if len(items) == 0 {
		return 0
	}

	var itemPairs = len(items) / 2
	var pairPoints = itemPairs * 5
	fmt.Printf("%d points - %d items (%d pairs @ 5 points each)\n", pairPoints, len(items), itemPairs)

	var descriptionPoints = 0

	for _, item := range items {

		var trimmedDescription = strings.TrimSpace(item.ShortDescription)
		var descriptionLength = len(trimmedDescription)

		if descriptionLength == 0 {
			continue
		}

		if descriptionLength%3 == 0 {

			priceFloat, err := strconv.ParseFloat(item.Price, 64)

			if err == nil {
				var p = priceFloat * 0.2
				var priceTotal = math.Ceil(p)
				descriptionPoints += int(priceTotal)

				fmt.Printf("%d points - \"%s\" is %d characters (a multiple of 3) item price of %v * 0.2 = %v, rounded up is 3 points\n", descriptionPoints, trimmedDescription, descriptionLength, priceFloat, p)

			}
		}
	}

	return descriptionPoints + pairPoints
}

func isMultipleOfFour(f float64) bool {

	shifted := f * 4
	return math.Floor(shifted) == shifted
}

func getPointsFromDateTime(purchaseDate string, purchaseTime string) int {

	var dateTime = purchaseDate + " " + purchaseTime
	d, err := time.Parse("2006-01-02 15:04", dateTime)
	if err != nil {
		return 0
	}

	var datePoints = 0
	var timePoints = 0

	day := d.Day()
	if day%2 != 0 {
		datePoints += 6
		fmt.Printf("6 points - purchase day \"%d\" is odd\n", day)
	}

	if d.Hour() >= 14 && d.Hour() <= 16 {
		timePoints += 10
		fmt.Printf("10 points - %d is between 2:00pm and 4:00pm\n", d.Hour())

	}

	return datePoints + timePoints
}

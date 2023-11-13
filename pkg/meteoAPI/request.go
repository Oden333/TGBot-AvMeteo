package meteo

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Определение структур для XML
type AviationProducts struct {
	XMLName   xml.Name                    `xml:"aviationProducts"`
	Forecasts []TerminalAerodromeForecast `xml:"terminalAerodromeForecast"`
}

type TerminalAerodromeForecast struct {
	ICAOAirportIdentifier string      `xml:"icaoAirportIdentifier"`
	IssuedTime            TimeInstant `xml:"issuedTime>TimeInstant"`
	ValidPeriod           ValidPeriod `xml:"validPeriod"`
	NAISHeader            string      `xml:"naisHeader"`
	TAFText               string      `xml:"tafText"`
}

type TimeInstant struct {
	ID           string `xml:"id,attr"`
	TimePosition string `xml:"timePosition"`
}

type ValidPeriod struct {
	BeginPosition string `xml:"beginPosition"`
	EndPosition   string `xml:"endPosition"`
}

func TAFRequest() (string, error) {
	url := "https://api.met.no/weatherapi/tafmetar/1.0/taf.xml?icao=UUDD"

	resp, err := http.Get(url)
	if err != nil {
		return "Ошибка при отправке запроса:", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Ошибка при чтении ответа:", err

	}

	var aviationProducts AviationProducts
	err = xml.Unmarshal(body, &aviationProducts)
	if err != nil {
		return "Ошибка при разборе XML:", err

	}

	var message string
	var max string
	var current *TerminalAerodromeForecast
	for _, forecast := range aviationProducts.Forecasts {
		if forecast.IssuedTime.ID > max {
			max = forecast.IssuedTime.ID
			current = &forecast
		}
	}

	message =
		`ICAO Airport Identifier: %s
Issued Time: %s
Valid Period: %s - %s
TAF:
		%s
		`
	del := strings.Trim(current.TAFText, "\n")
	del = strings.Trim(del, "	")
	del = strings.Trim(del, "=")
	message = fmt.Sprintf(message,
		current.ICAOAirportIdentifier,
		strings.ReplaceAll(current.IssuedTime.TimePosition, "T", "	"),
		strings.ReplaceAll(current.ValidPeriod.BeginPosition, "T", " "),
		strings.ReplaceAll(current.ValidPeriod.EndPosition, "T", " "),
		del)

	return message, nil
}

func MeteoRequest() (string, error) {
	url := "https://api.met.no/weatherapi/tafmetar/1.0/metar.txt?icao=UUDD"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	fin := string(body)
	a := strings.Split(fin, "=\n")
	fin = `METAR:
	%s`
	fin = fmt.Sprintf(fin, a[len(a)-2])
	return fin, err

}

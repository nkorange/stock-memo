package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/nkorange/stock-memo/pkg/stock"
	"github.com/nkorange/stock-memo/pkg/trade_strategy"
	"html/template"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type Server struct {
	*http.Server
	addr string
}

func NewServer(addr string) (*Server, error) {
	s := &Server{}
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/analyze", s.analyze)
	serverMux.HandleFunc("/", s.dashboard)
	s.Server = &http.Server{
		Addr:    addr,
		Handler: serverMux,
	}
	return s, nil
}

func (s *Server) Run(ctx context.Context) error {
	fmt.Println("listening to", s.Server.Addr)
	ln, err := net.Listen("tcp", s.Server.Addr)
	if err != nil {
		return err
	}
	return s.Server.Serve(ln)
}

func (s *Server) dashboard(w http.ResponseWriter, req *http.Request) {
	tmpl, _ := template.ParseFiles("html/dashboard.html")
	tmpl.Execute(w, nil)
}

func (s *Server) analyze(w http.ResponseWriter, req *http.Request) {

	//startMoney, err := strconv.ParseFloat(req.Form.Get("startMoney"), 64)
	//if err != nil {
	//	panic(err)
	//}
	//startDay, err := strconv.Atoi(req.Form.Get("startDay"))
	//if err != nil {
	//	panic(err)
	//}
	startMoney := 10000.0
	startDay := 1
	err := req.ParseMultipartForm(5 * 1024 * 1024)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	file, _, err := req.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	io.Copy(&buf, file)
	contents := buf.String()
	lines := strings.Split(contents, "\n")
	history := &stock.PriceHistory{
		Prices: make([]*stock.Price, len(lines)-1),
	}
	ind := len(history.Prices) - 1
	for _, line := range lines[1:] {
		elements := strings.Split(line, "\"")
		closePrice, _ := strconv.ParseFloat(elements[3], 64)
		openPrice, _ := strconv.ParseFloat(elements[5], 64)
		highestPrice, _ := strconv.ParseFloat(elements[7], 64)
		lowestPrice, _ := strconv.ParseFloat(elements[9], 64)
		price := &stock.Price{
			Date:         elements[1],
			ClosePrice:   closePrice,
			OpenPrice:    openPrice,
			HighestPrice: highestPrice,
			LowestPrice:  lowestPrice,
		}
		history.Prices[ind] = price
		ind--
	}

	tradeStrategy := trade_strategy.NewPercentageStrategy()
	profit := tradeStrategy.Trade(history, startMoney, startDay)

	writeData(w, profit)
}

func writeError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	_, err = w.Write([]byte(err.Error()))
	if err != nil {
		fmt.Println(err)
	}
}

func writeInternalServerError(w http.ResponseWriter, err error) {
	writeError(w, http.StatusInternalServerError, err)
}

func writeData(w http.ResponseWriter, data interface{}) {
	js, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		fmt.Println("invalid json data", data, err)
		writeInternalServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(js)
	if err != nil {
		fmt.Println("fail to write data", err)
	}
}

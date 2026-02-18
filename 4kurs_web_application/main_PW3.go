package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

type PageData struct {
	Pc     string
	Sigma1 string
	Sigma2 string
	Price  string
	Res1   string
	Res2   string
	Profit string
}

func main() {
	http.HandleFunc("/", calcHandler)
	fmt.Println("Server running: http://localhost:8081")
	http.ListenAndServe(":8081", nil)
}

func calcHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Pc:     "5.0",
		Sigma1: "1.0",
		Sigma2: "0.25",
		Price:  "7.0",
	}

	if r.Method == http.MethodPost {
		data.Pc = r.FormValue("pc")
		data.Sigma1 = r.FormValue("sigma1")
		data.Sigma2 = r.FormValue("sigma2")
		data.Price = r.FormValue("price")

		pc := parse(data.Pc)
		s1 := parse(data.Sigma1)
		s2 := parse(data.Sigma2)
		p := parse(data.Price)

		profit1 := calculate(pc, s1, p)
		profit2 := calculate(pc, s2, p)

		data.Res1 = fmt.Sprintf("%.2f", profit1)
		data.Res2 = fmt.Sprintf("%.2f", profit2)
		data.Profit = fmt.Sprintf("%.2f", profit2-profit1)
	}

	t, _ := template.New("p").Parse(htmlTemplate)
	t.Execute(w, data)
}

func calculate(pc, sigma, price float64) float64 {
	delta := pc * 0.05
	p1 := pc - delta
	p2 := pc + delta

	integral := 0.5 * (math.Erf((p2-pc)/(sigma*math.Sqrt(2))) - math.Erf((p1-pc)/(sigma*math.Sqrt(2))))
	
	w1 := pc * 24 * integral
	w2 := pc * 24 * (1 - integral)

	profit := (w1 * price * 1000) - (w2 * price * 1000)
	return profit / 1000
}

func parse(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Lab 3 - Var 5</title>
    <style>
        body { font-family: sans-serif; padding: 20px; background: #f0f2f5; }
        .box { background: white; padding: 30px; width: 400px; margin: 0 auto; border-radius: 8px; box-shadow: 0 2px 5px rgba(0,0,0,0.1); }
        h2 { text-align: center; color: #333; }
        label { display: block; margin-top: 15px; color: #555; }
        input { width: 100%; padding: 8px; margin-top: 5px; border: 1px solid #ddd; border-radius: 4px; box-sizing: border-box;}
        button { width: 100%; background: #007bff; color: white; padding: 10px; border: none; margin-top: 20px; border-radius: 4px; cursor: pointer; }
        .res { margin-top: 20px; padding: 15px; background: #e8f5e9; border: 1px solid #c8e6c9; border-radius: 4px; }
    </style>
</head>
<body>
    <div class="box">
        <h2>Розрахунок СЕС (Варіант 5)</h2>
        <form method="POST">
            <label>Потужність (Pc), МВт</label>
            <input type="text" name="pc" value="{{.Pc}}">
            <label>Похибка 1 (σ1)</label>
            <input type="text" name="sigma1" value="{{.Sigma1}}">
            <label>Похибка 2 (σ2)</label>
            <input type="text" name="sigma2" value="{{.Sigma2}}">
            <label>Вартість (B), грн/кВт·год</label>
            <input type="text" name="price" value="{{.Price}}">
            <button type="submit">Розрахувати</button>
        </form>
        {{if .Res1}}
        <div class="res">
            <p>Прибуток (система 1): {{.Res1}} тис. грн</p>
            <p>Прибуток (система 2): {{.Res2}} тис. грн</p>
            <p><strong>Виграш: {{.Profit}} тис. грн</strong></p>
        </div>
        {{end}}
    </div>
</body>
</html>
`
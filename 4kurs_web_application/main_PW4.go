package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

type PageData struct {
	Sm      string
	Unom    string
	Ik      string
	Tf      string
	Ct      string
	Jek     string
	Sk      string
	Ucn     string
	SnomT   string
	Uk      string
	Rcn     string
	Xcn     string
	Rcmin   string
	Xcmin   string
	Uvn     string
	Unn     string
	SnomT3  string
	Ukmax   string
	Res1    string
	Res2    string
	Res3    string
}

func main() {
	http.HandleFunc("/", calcHandler)
	http.ListenAndServe(":8081", nil)
}

func calcHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Sm: "1300", Unom: "10", Ik: "2500", Tf: "2.5", Ct: "92", Jek: "1.4",
		Sk: "200", Ucn: "10.5", SnomT: "6.3", Uk: "10.5",
		Rcn: "10.65", Xcn: "24.02", Rcmin: "34.88", Xcmin: "65.68", Uvn: "115", Unn: "11", SnomT3: "6.3", Ukmax: "11.1",
	}

	if r.Method == http.MethodPost {
		data.Sm = r.FormValue("sm")
		data.Unom = r.FormValue("unom")
		data.Ik = r.FormValue("ik")
		data.Tf = r.FormValue("tf")
		data.Ct = r.FormValue("ct")
		data.Jek = r.FormValue("jek")
		data.Sk = r.FormValue("sk")
		data.Ucn = r.FormValue("ucn")
		data.SnomT = r.FormValue("snomt")
		data.Uk = r.FormValue("uk")
		data.Rcn = r.FormValue("rcn")
		data.Xcn = r.FormValue("xcn")
		data.Rcmin = r.FormValue("rcmin")
		data.Xcmin = r.FormValue("xcmin")
		data.Uvn = r.FormValue("uvn")
		data.Unn = r.FormValue("unn")
		data.SnomT3 = r.FormValue("snomt3")
		data.Ukmax = r.FormValue("ukmax")

		sm := parse(data.Sm)
		unom := parse(data.Unom)
		ik := parse(data.Ik)
		tf := parse(data.Tf)
		ct := parse(data.Ct)
		jek := parse(data.Jek)

		im := (sm / 2) / (math.Sqrt(3) * unom)
		impa := 2 * im
		sek := im / jek
		smin := (ik * math.Sqrt(tf)) / ct

		data.Res1 = fmt.Sprintf("Струм норм. режиму: %.2f А\nСтрум післяаварійний: %.2f А\nЕкономічний переріз: %.2f мм2\nМінімальний переріз (терм. стійкість): %.2f мм2", im, impa, sek, smin)

		sk := parse(data.Sk)
		ucn := parse(data.Ucn)
		snomt := parse(data.SnomT)
		uk := parse(data.Uk)

		xc := math.Pow(ucn, 2) / sk
		xt := (uk / 100) * (math.Pow(ucn, 2) / snomt)
		xsum := xc + xt
		ip0 := ucn / (math.Sqrt(3) * xsum)

		data.Res2 = fmt.Sprintf("Опір системи Xc: %.2f Ом\nОпір трансформатора Xt: %.2f Ом\nСумарний опір: %.2f Ом\nПочатковий струм КЗ: %.2f кА", xc, xt, xsum, ip0)

		rcn := parse(data.Rcn)
		xcn := parse(data.Xcn)
		rcmin := parse(data.Rcmin)
		xcmin := parse(data.Xcmin)
		uvn := parse(data.Uvn)
		unn := parse(data.Unn)
		snomt3 := parse(data.SnomT3)
		ukmax := parse(data.Ukmax)

		xt3 := (ukmax * math.Pow(uvn, 2)) / (100 * snomt3)
		rn := rcn
		xn := xcn + xt3
		rmin := rcmin
		xmin := xcmin + xt3

		kpr := math.Pow(unn, 2) / math.Pow(uvn, 2)
		rshn := rn * kpr
		xshn := xn * kpr
		zshn := math.Sqrt(math.Pow(rshn, 2) + math.Pow(xshn, 2))

		rshmin := rmin * kpr
		xshmin := xmin * kpr
		zshmin := math.Sqrt(math.Pow(rshmin, 2) + math.Pow(xshmin, 2))

		ishn3 := (unn * 1000) / (math.Sqrt(3) * zshn)
		ishn2 := ishn3 * (math.Sqrt(3) / 2)
		ishmin3 := (unn * 1000) / (math.Sqrt(3) * zshmin)
		ishmin2 := ishmin3 * (math.Sqrt(3) / 2)

		data.Res3 = fmt.Sprintf("Нормальний режим: I(3) = %.0f А, I(2) = %.0f А\nМінімальний режим: I(3) = %.0f А, I(2) = %.0f А", ishn3, ishn2, ishmin3, ishmin2)
	}

	t, _ := template.New("p").Parse(htmlTemplate)
	t.Execute(w, data)
}

func parse(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
<title>Lab 4 - Var</title>
<style>
body { font-family: sans-serif; padding: 20px; background: #f0f2f5; }
.box { background: white; padding: 30px; width: 600px; margin: 0 auto; border-radius: 8px; box-shadow: 0 2px 5px rgba(0,0,0,0.1); }
h2, h3 { text-align: center; color: #333; }
label { display: block; margin-top: 15px; color: #555; font-size: 14px;}
input { width: 100%; padding: 8px; margin-top: 5px; border: 1px solid #ddd; border-radius: 4px; box-sizing: border-box;}
button { width: 100%; background: #007bff; color: white; padding: 10px; border: none; margin-top: 20px; border-radius: 4px; cursor: pointer; }
.res { margin-top: 20px; padding: 15px; background: #e8f5e9; border: 1px solid #c8e6c9; border-radius: 4px; white-space: pre-line; }
.grid { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }
</style>
</head>
<body>
<div class="box">
<h2>Розрахунок струмів КЗ</h2>
<form method="POST">
<h3>Завдання 1</h3>
<div class="grid">
<div><label>Sm (кВА)</label><input type="text" name="sm" value="{{.Sm}}"></div>
<div><label>Unom (кВ)</label><input type="text" name="unom" value="{{.Unom}}"></div>
<div><label>Ik (А)</label><input type="text" name="ik" value="{{.Ik}}"></div>
<div><label>tФ (с)</label><input type="text" name="tf" value="{{.Tf}}"></div>
<div><label>Ct</label><input type="text" name="ct" value="{{.Ct}}"></div>
<div><label>jek</label><input type="text" name="jek" value="{{.Jek}}"></div>
</div>

<h3>Завдання 2</h3>
<div class="grid">
<div><label>Sk (МВА)</label><input type="text" name="sk" value="{{.Sk}}"></div>
<div><label>Ucn (кВ)</label><input type="text" name="ucn" value="{{.Ucn}}"></div>
<div><label>SnomT (МВА)</label><input type="text" name="snomt" value="{{.SnomT}}"></div>
<div><label>Uk (%)</label><input type="text" name="uk" value="{{.Uk}}"></div>
</div>

<h3>Завдання 3</h3>
<div class="grid">
<div><label>Rc.n</label><input type="text" name="rcn" value="{{.Rcn}}"></div>
<div><label>Xc.n</label><input type="text" name="xcn" value="{{.Xcn}}"></div>
<div><label>Rc.min</label><input type="text" name="rcmin" value="{{.Rcmin}}"></div>
<div><label>Xc.min</label><input type="text" name="xcmin" value="{{.Xcmin}}"></div>
<div><label>Uvn</label><input type="text" name="uvn" value="{{.Uvn}}"></div>
<div><label>Unn</label><input type="text" name="unn" value="{{.Unn}}"></div>
<div><label>SnomT</label><input type="text" name="snomt3" value="{{.SnomT3}}"></div>
<div><label>Uk.max</label><input type="text" name="ukmax" value="{{.Ukmax}}"></div>
</div>

<button type="submit">Розрахувати</button>
</form>

{{if .Res1}}
<div class="res">
<strong>Результат Завдання 1:</strong>
{{.Res1}}
</div>
<div class="res">
<strong>Результат Завдання 2:</strong>
{{.Res2}}
</div>
<div class="res">
<strong>Результат Завдання 3:</strong>
{{.Res3}}
</div>
{{end}}

</div>
</body>
</html>
`
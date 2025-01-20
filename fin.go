package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type PageData struct {
	Result       string
	UnitTypes    []string
	FromUnits    []string
	ToUnits      []string
	SelectedType string
	Value        string
	FromUnit     string
	ToUnit       string
}

const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>UnitLab</title>
    <style>
        @import url('https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@300;400&display=swap');

        body {
            font-family: 'Space Grotesk', sans-serif;
            margin: 0;
            padding: 0;
            min-height: 100vh;
            background-color: #000000;
            color: #ffffff;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
        }

        .container {
            width: 100%;
            max-width: 800px;
            text-align: center;
            padding: 20px;
        }

        h1 {
            font-size: 4.5rem;
            font-weight: 300;
            margin: 0;
            letter-spacing: 2px;
            animation: fadeIn 1s ease-out;
        }

        .subtitle {
            font-size: 2.5rem;
            font-weight: 300;
            margin-bottom: 6rem;
            opacity: 0.9;
            animation: fadeIn 1s ease-out 0.5s both;
        }

        .conversion-form {
            display: flex;
            flex-direction: column;
            gap: 2rem;
            align-items: center;
            position: relative;
        }

        .input-group {
            display: flex;
            gap: 1rem;
            width: 100%;
            justify-content: center;
        }

        input, select {
            background: #ffffff;
            border: none;
            padding: 1.2rem;
            font-size: 1.2rem;
            font-family: 'Space Grotesk', sans-serif;
            color: #000000;
            width: 300px;
            border-radius: 0;
            -webkit-appearance: none;
            appearance: none;
            transition: all 0.3s ease;
            animation: slideIn 0.5s ease-out both;
        }

        select {
            width: 150px;
            cursor: pointer;
            padding-right: 2rem;
            background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='black' stroke-width='1.5' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
            background-repeat: no-repeat;
            background-position: right 8px center;
            background-size: 12px;
        }

        .converting-text {
            font-size: 1.2rem;
            letter-spacing: 3px;
            margin: 1.5rem 0;
            font-weight: 300;
            animation: pulse 2s infinite;
        }

        #unitType {
            position: relative;
            margin: 0 auto 2rem auto;
            width: auto;
            min-width: 150px;
            padding: 0.8rem 2rem 0.8rem 1rem;
            font-size: 1rem;
            background: transparent;
            color: white;
            border: 1px solid rgba(255,255,255,0.3);
            animation: fadeIn 1s ease-out 1s both;
        }

        #unitType::before {
            content: "Unit Type: ";
            opacity: 0.7;
        }

        input:focus, select:focus {
            outline: none;
            box-shadow: 0 0 0 2px rgba(255,255,255,0.2);
        }

        input[readonly] {
            background: #ffffff;
            color: #000000;
        }

        @keyframes fadeIn {
            from { opacity: 0; }
            to { opacity: 1; }
        }

        @keyframes slideIn {
            from { transform: translateY(20px); opacity: 0; }
            to { transform: translateY(0); opacity: 1; }
        }

        @keyframes pulse {
            0% { opacity: 1; }
            50% { opacity: 0.5; }
            100% { opacity: 1; }
        }

        .input-group:nth-child(1) { animation-delay: 0.2s; }
        .input-group:nth-child(2) { animation-delay: 0.4s; }

        .convert-button {
            background: #ffffff;
            border: none;
            padding: 1.2rem 2.5rem;
            font-size: 1.2rem;
            font-family: 'Space Grotesk', sans-serif;
            color: #000000;
            cursor: pointer;
            transition: all 0.3s ease;
            animation: slideIn 0.5s ease-out both;
            animation-delay: 0.6s;
        }

        .convert-button:hover {
            background: rgba(255,255,255,0.9);
            transform: translateY(-2px);
        }

        .convert-button:active {
            transform: translateY(0);
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>UnitLab</h1>
        <div class="subtitle">bernininini</div>

        <form method="POST" class="conversion-form" id="conversionForm">
            <select name="unitType" id="unitType">
                {{range .UnitTypes}}
                    <option value="{{.}}" {{if eq . $.SelectedType}}selected{{end}}>{{.}}</option>
                {{end}}
            </select>

            <div class="input-group">
                <input type="number" 
                       name="value" 
                       step="any" 
                       required 
                       placeholder="Enter value" 
                       value="{{.Value}}">
                <select name="from">
                    {{range .FromUnits}}
                        <option value="{{.}}" {{if eq . $.FromUnit}}selected{{end}}>{{.}}</option>
                    {{end}}
                </select>
            </div>

            <div class="converting-text">CONVERTING TO.......</div>

            <div class="input-group">
                <input type="text" 
                       value="{{if .Result}}{{.Result}}{{else}}?{{end}}" 
                       readonly>
                <select name="to">
                    {{range .ToUnits}}
                        <option value="{{.}}" {{if eq . $.ToUnit}}selected{{end}}>{{.}}</option>
                    {{end}}
                </select>
            </div>

            <button type="submit" class="convert-button">Convert</button>
        </form>
    </div>

    <script>
        window.onload = function() {
            document.querySelector('input[name="value"]').focus();
            document.querySelector('#unitType').addEventListener('change', function() {
                document.getElementById('conversionForm').submit();
            });
        }
    </script>
</body>
</html>
`

func main() {
	http.HandleFunc("/", handleConvert)
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func handleConvert(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("converter").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	unitTypes := []string{"length", "mass", "temperature", "time"}
	data := PageData{
		UnitTypes: unitTypes,
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		unitType := r.FormValue("unitType")
		value := r.FormValue("value")
		from := r.FormValue("from")
		to := r.FormValue("to")

		data.SelectedType = unitType
		data.Value = value
		data.FromUnit = from
		data.ToUnit = to
		data.FromUnits, data.ToUnits = getUnitsForType(unitType)

		if value != "" {
			val, err := strconv.ParseFloat(value, 64)
			if err == nil {
				result, err := convert(val, from, to, unitType)
				if err != nil {
					data.Result = "?"
				} else {
					data.Result = fmt.Sprintf("%.2f", result)
				}
			}
		}
	} else {
		data.SelectedType = "length"
		data.FromUnits, data.ToUnits = getUnitsForType(data.SelectedType)
	}

	tmpl.Execute(w, data)
}

func getUnitsForType(unitType string) ([]string, []string) {
	units := map[string][]string{
		"length":      {"meters", "feet", "inches", "kilometers", "centimeters", "miles", "yards"},
		"mass":        {"kilograms", "pounds", "grams", "ounces", "tons"},
		"temperature": {"celsius", "fahrenheit", "kelvin"},
		"time":        {"seconds", "minutes", "hours", "days", "weeks"},
	}

	if unitList, ok := units[unitType]; ok {
		return unitList, unitList
	}
	return []string{}, []string{}
}

func convert(value float64, from, to, unitType string) (float64, error) {
	conversions := map[string]map[string]float64{
		"length": {
			"meters":      1.0,
			"feet":        0.3048,
			"inches":      0.0254,
			"kilometers":  1000.0,
			"centimeters": 0.01,
			"miles":       1609.34,
			"yards":       0.9144,
		},
		"mass": {
			"kilograms": 1.0,
			"pounds":    0.453592,
			"grams":     0.001,
			"ounces":    0.0283495,
			"tons":      907.185,
		},
		"time": {
			"seconds": 1.0,
			"minutes": 60.0,
			"hours":   3600.0,
			"days":    86400.0,
			"weeks":   604800.0,
		},
	}

	if unitType == "temperature" {
		return convertTemperature(value, from, to)
	}

	if factors, ok := conversions[unitType]; ok {
		if fromFactor, ok := factors[from]; ok {
			if toFactor, ok := factors[to]; ok {
				baseValue := value * fromFactor
				return baseValue / toFactor, nil
			}
		}
	}

	return 0, fmt.Errorf("invalid conversion")
}

func convertTemperature(value float64, from, to string) (float64, error) {
	var celsius float64
	switch from {
	case "celsius":
		celsius = value
	case "fahrenheit":
		celsius = (value - 32) * 5 / 9
	case "kelvin":
		celsius = value - 273.15
	default:
		return 0, fmt.Errorf("invalid temperature unit")
	}

	switch to {
	case "celsius":
		return celsius, nil
	case "fahrenheit":
		return celsius*9/5 + 32, nil
	case "kelvin":
		return celsius + 273.15, nil
	default:
		return 0, fmt.Errorf("invalid temperature unit")
	}
}

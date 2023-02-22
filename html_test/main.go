package main

import (
	"html/template"
	"os"
)

func main() {
	// Create the map[string][]string
	data := make(map[string][]string)
	data["Compute"] = []string{"BGP", "Pass", "Founded"}
	data["Storage"] = []string{"LLDP", "Fail", "asdfasdfadsgftrgwertgwaerfafasfadfareaqwefastrsertsergewsrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr"}

	// Create the HTML template
	htmlTemplate := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Summary Report</title>
			<style>
				table {
					border-collapse: collapse;
					width: 100%;
				}
				th, td {
					border: 1px solid black;
					padding: 8px;
					text-align: center;
					word-break: break-all; 
				}
				th {
					background-color: #6c6c6d;
				}
				.big-cell {
					width: 60%;
				}
				.small-cell {
					width: 20%; 
				}
				.pass { background-color: green; }
				.fail { background-color: red; }
			</style>
		</head>
		<body>
		{{range $key, $value := .}}
		<h1>{{$key}}</h1>
		<table>
		<tr>
		  <th class="small-cell">Column 1</th>
		  <th class="small-cell">Column 2</th>
		  <th class="big-cell">Column 3</th>
		</tr>
			<tr class="{{if eq (index $value 1) "Pass"}}pass{{else if eq (index $value 1) "Fail"}}fail{{end}}">
				{{range $index, $element := $value}}
					<td>{{$element}}</td>
				{{end}}
			</tr>
	  	</table>
	  	{{end}}
		</body>
		</html>
	`

	// Parse the HTML template
	tmpl, err := template.New("htmlTemplate").Parse(htmlTemplate)
	if err != nil {
		panic(err)
	}

	// Execute the template and write the output to a file
	file, err := os.Create("table.html")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = tmpl.Execute(file, data)
	if err != nil {
		panic(err)
	}
}

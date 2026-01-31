module github.com/mlctrez/lotterydisplay

go 1.24

//replace github.com/briandowns/openweathermap v0.21.1 => /home/mattman/golang/briandowns/openweathermap

require (
	github.com/briandowns/openweathermap v0.21.1
	github.com/kardianos/service v1.2.2
	github.com/mlctrez/servicego v1.4.10
	github.com/robfig/cron/v3 v3.0.1
)

require golang.org/x/sys v0.33.0 // indirect

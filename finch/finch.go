package finch

import (
	"appengine"
	"appengine/urlfetch"
	"html/template"
	"net/http"
)

const homeTemplateHTML = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>Finch.gae</title>
  <link rel="stylesheet" href="//netdna.bootstrapcdn.com/twitter-bootstrap/2.3.0/css/bootstrap-combined.min.css">
</head>
<body>
  <main class="container">
    <header class="jumbo subhead">
      <div class="container"><h1>Finch.gae</h1></div>
    </header>
    <table class="table table-bordered">
      <tr>
        <th>URL</th>
        <th>Status</th>
      </tr>
      {{range .}}
      <tr class="{{.Status}}">
        <td>{{.Url}}</td>
        <td>{{.Status}}</td>
      </tr>
      {{end}}
    </table>
  </main>
</body>
</html>
`

var homeTemplate = template.Must(template.New("home").Parse(homeTemplateHTML))

var hosts = []string{"http://google.com/", "http://linuxfr.org/", "http://down.example.net/"}

type Service struct {
	Url    string
	Status string
}

func check(c appengine.Context, url string) string {
	client := urlfetch.Client(c)
	_, err := client.Get(url)
	if err != nil {
		//return err.Error()
		return "error"
	}
	return "success"
}

func services(c appengine.Context) []Service {
	services := []Service{}
	for _, host := range hosts {
		status := check(c, host)
		services = append(services, Service{host, status})
	}
	return services
}

func handler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	err := homeTemplate.Execute(w, services(c))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func init() {
	http.HandleFunc("/", handler)
}

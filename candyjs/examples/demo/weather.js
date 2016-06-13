var http = CandyJS.require('net/http');
var ioutil = CandyJS.require('io/ioutil');

resp = http.get('http://api.openweathermap.org/data/2.5/weather?q=Madrid,ES&units=metric');

if (resp.statusCode == 200) {
  var json = ioutil.readAll(resp.body);
  var obj = JSON.parse(json);

  print(obj.name + " temperature " + Math.round(obj.main.temp) + "°C");
} else {
  print('Request failed, status code:', resp.statusCode);
}
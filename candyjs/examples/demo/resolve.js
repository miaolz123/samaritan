var net = CandyJS.require('net');

var host = 'google.com';

print(host + ' resolve to:');
net.lookupHost(host).forEach(function(ip) {
  print(ip)
});
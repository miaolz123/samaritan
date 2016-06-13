# Let's build a interpreter for our JavaScripts.
ccat main.go

clear
go generate .
go build .

clear
# Now we run weather.js, its makes a HTTP request.
ccat weather.js

./example weather.js

clear
# And now we run resolve.js, a totally different code.
ccat resolve.js

./example resolve.js
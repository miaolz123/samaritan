var time = CandyJS.require('time');
var gin = CandyJS.require('github.com/gin-gonic/gin');

var engine = gin.default();
engine.get("/back", CandyJS.proxy(function(ctx) {
  var future = time.date(2015, 10, 21, 4, 29 ,0, 0, time.UTC);
  var now = time.now();

  ctx.json(200, {
    future: future.string(),
    now: now.string(),
    nsecs: future.sub(now)
  });
}));

engine.run(':8080');
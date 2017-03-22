// Simple server for hosting the client view.
var livereload = require('express-livereload')
var express = require('express')
var app = express()

/*** Live reloading for development ***/
var livereload = require('livereload');
var server = livereload.createServer();
server.watch(__dirname + "/public");

app.use(express.static(__dirname + '/public'));

app.get('*', function(req, res) {
    res.sendFile(__dirname + '/public/index.html');
});

const port = 3000;

app.listen(port, function() {
    console.log('Client hosting server listening on port %d!', port)
})

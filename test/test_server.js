var http = require('http');

//We need a function which handles requests and send response
function handleRequest(request, response){
  console.log("Request");
  console.log(request.body);
  console.log(request.headers);
  response.writeHead(200, {'Content-Type': 'text/json'});
  response.end('{"message": "Yay"}');
}

//Create a server
var server = http.createServer(handleRequest);

//Lets start our server
server.listen(8080, function(){
    console.log("Server listening on: http://localhost:%s", 8080);
});

/**
 * Created by Carl on 6/6/2017.
 */
// Reusable HTTP request object.
var serverRequest = new XMLHttpRequest();
/* Call to register given user ID with KoAuth. */
function enableKoAuth(request) {
    // @TODO: Change to your server's IP here.
    serverRequest.open("POST", "http://localhost:8080/user/new", true);
    serverRequest.setRequestHeader("Content-Type", "application/json");
    serverRequest.addEventListener("readystatechange", onEnableRequestResponse, false);
    serverRequest.send(JSON.stringify(request));
}
function onEnableRequestResponse(e) {
    if (serverRequest.readyState === 4) {
        if (serverRequest.status === 200) {
            var response = JSON.parse(serverRequest.response);
            document.getElementById("serial").innerText = "Serial: " + response.Serial;
        }
        else {
            console.log(e);
        }
    }
}
function onSubmit() {
    var input = document.getElementById("input").value;
}

/**
 * Created by Carl on 6/6/2017.
 */
// Reusable HTTP request object.
var serverRequest = new XMLHttpRequest();
/* Call to register given user ID with KoAuth. */
function enableKoAuth(userID) {
    // @TODO: Change to your server's IP here.
    serverRequest.open("POST", "http://localhost:8080/user/new", true);
    serverRequest.addEventListener("readystatechange", onEnableRequestResponse, false);
    serverRequest.send(userID);
}
function onEnableRequestResponse(e) {
    if (serverRequest.readyState === 4 && serverRequest.status === 200) {
        var response = JSON.parse(serverRequest.responseText);
        console.log(response);
    }
}

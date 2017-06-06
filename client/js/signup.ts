/**
 * Created by Carl on 6/6/2017.
 */

// Reusable HTTP request object.
let serverRequest: XMLHttpRequest = new XMLHttpRequest();

interface RequestBody {
    userId: number
}

/* Call to register given user ID with KoAuth. */
function enableKoAuth(request: RequestBody) {
    // @TODO: Change to your server's IP here.
    serverRequest.open("POST", "http://localhost:8080/user/new", true);
    serverRequest.addEventListener("readystatechange", onEnableRequestResponse, false);
    serverRequest.send(request);
}

function onEnableRequestResponse(e) {
    if (serverRequest.readyState === 4 && serverRequest.status === 200) {
        const response = JSON.parse(serverRequest.response);
        console.log(response);
    }
}
/**
 * Created by Carl on 6/6/2017.
 */

// Reusable HTTP request object.
let serverRequest: XMLHttpRequest = new XMLHttpRequest();

interface RequestBody {
    Email: string,
    DeviceSerial?: string,
    Serial?: string
}

/* Call to register given user ID with KoAuth. */
function enableKoAuth(request: RequestBody) {
    // @TODO: Change to your server's IP here.
    serverRequest.open("POST", "http://localhost:8080/user/new", true);
    serverRequest.setRequestHeader("Content-Type", "application/json");
    serverRequest.addEventListener("readystatechange", onEnableRequestResponse, false);
    serverRequest.send(JSON.stringify(request));
}

function onEnableRequestResponse(e) {
    if (serverRequest.readyState === 4) {
        if(serverRequest.status === 200) {
            const response = JSON.parse(serverRequest.response);
            document.getElementById("serial").innerText = "Serial: " + response.Serial;
            alert(JSON.stringify(response));
        } else {
            alert(e.toString());
        }
    }
}

function onSubmit() {
    let input = (<HTMLInputElement>document.getElementById("input")).value;
    alert(input);
}
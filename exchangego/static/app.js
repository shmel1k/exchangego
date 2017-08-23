let canvas = document.getElementById("canvas");

let base = 10;

let history = null;

function generateArray(history) {
    let ar = [];

    let len = history.length;
    for (let i = 0; i < len; i++) {
        // TODO x time
        ar.push([i, history[i]]);
    }

    return ar;
}

function draw(graph) {
    $.plot($("#placeholder"), [graph]);
}

$.get("/get?size=10", function(data) {
    history = data["history"];

    draw( generateArray(history) );
    let ws = new WebSocket('ws://localhost:4242/ws');

    ws.addEventListener('message', function(e) {
        let msg = JSON.parse(e.data);

        history = history.slice(1);
        history.push(msg["value"]);

        draw( generateArray(history) );
    });

});
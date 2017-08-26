let canvas = document.getElementById("canvas");
let balanceElement = document.getElementById("money");
let statusElement = document.getElementById("status");

let lastCurrency;

let base = 10;
let history = null;

let line = null;

function generateArray(history) {
    let ar = [];
    let ar_line = [];

    let len = history.length;
    for (let i = 0; i < len; i++) {
        // TODO x time
        ar.push([i, history[i]]);
        if (line != null) {
            ar_line.push([i, line])
        }
    }
    return [ar, ar_line];
}

function getRealHost() {
    return window.location.hostname + (window.location.port ? ':'+ window.location.port : '');
}

window.startGame = function(type) {
    /* TODO loading */
    $.get("/game?type=" + type + "&seconds=10", function (e) {
        statusElement.innerHTML = "Start Go " + type + "!";

        line = lastCurrency;
    });
};

function draw(graph) {
    $.plot($("#placeholder"), graph);
}

$.get("/get?size=20", function(data) {
    history = data["history"];

    lastCurrency = history[history.length - 1];
    draw(generateArray(history));
    let ws = new WebSocket('ws://' + getRealHost() + '/ws');

    ws.addEventListener('message', function(e) {
        let msg = JSON.parse(e.data);
        console.log(msg);
        if ("status" in msg) {
            line = null;

            balanceElement.innerHTML = msg["money"] + "";

            let bonus = "";
            if (msg["status"] == "win")
                bonus = "Get +100";
            else
                bonus = "Lost -100";

            statusElement.innerHTML = "You " + msg["status"] + ". " + bonus;
            return
        }

        lastCurrency = history[history.length - 1];

        history = history.slice(1);
        history.push(msg["value"]);

        draw( generateArray(history) );
    });

});
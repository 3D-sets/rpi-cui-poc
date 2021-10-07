
function setAngle() {
    var angle = $("#servo-slider").roundSlider("option", "value");
    $.ajax({
        url: "servo",
        type: "get",
        data: {
            angle: angle
        }
    });
}

function setSpeed() {
    var speed = $("#motor-slider").roundSlider("option", "value");
    $.ajax({
        url: "motor",
        type: "get",
        data: {
            angle: speed
        }
    });
}

$(document).ready(function() {
    $("#servo-slider").roundSlider({
        radius: 360,
        circleShape: "custom-quarter",
        min: -45,
        max: 45,
        value: 0,
        startAngle: "45",
        handleShape: "dot",
        update: "setAngle"
    });
    $("#motor-slider").roundSlider({
        radius: 180,
        circleShape: "custom-quarter",
        min: -45,
        max: 45,
        value: 0,
        startAngle: "135",
        handleShape: "dot",
        update: "setSpeed"
    }); 
});
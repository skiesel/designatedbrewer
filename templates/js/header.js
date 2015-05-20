function sendNotification(message, sound, timeout, remote) {
	showLocalMessage(message, timeout);
	
	if(sound === "ALERT") { $("#alertSound")[0].play(); }
	else if(sound === "DONE") { $("#doneSound")[0].play(); }

	if(remote) {
		sendAlert(message, message);
	}
}

function showLocalMessage(message, closeTimeout) {
	$("#messageArea").html(message);
	$("#alert").fadeIn();

	if(typeof closeTimeout === "number") {
		window.setTimeout(function() {
			closeLocalMessage();
		}, closeTimeout);
	}
}

function closeLocalMessage() {
	$("#alert").fadeOut();
}

$(document).ready(function() {
	$("#closeLocalMessage").click(closeLocalMessage);
});
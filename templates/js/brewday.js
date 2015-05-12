var step = 0;
var autoadvance = false;
var timer = null;
var timeRemaining = null;
var timeElapsed = null;
var temp0InRange = false;
var temp1InRange = false;

$("#autoadvance").change(function() { autoadvance = !autoadvance; });

function formatTime(t) {
	var str = "";
	if(t.hours() < 10) { str += "0" + t.hours(); }
	else { str += t.hours(); }
	str += ":";
	if(t.minutes() < 10) { str += "0" + t.minutes(); }
	else { str += t.minutes(); }
	str += ":";
	if(t.seconds() < 10) { str += "0" + t.seconds(); }
	else { str += t.seconds(); }
	return str;
}

function getFirstSpargeTemp() {
	for(var i = 0; i < schedule.length; i++) {
		if(schedule[i].Type === "SPARGE") {
			return schedule[i].Temperature;
		}
	}
	return -1;
}

function showCurrentStep() {
	if(step > 0) {
		$($(".scheduleStep")[step-1]).removeClass("warning");
	}

	$($(".scheduleStep")[step]).addClass("warning");

	if(schedule[step].Type === "CHILL") {
		$("#temperature0Label").html("Wort Temp");
		$("#temperature1Label").html("Chill Water Temp");

		timeElapsed = moment.duration(0);
		$("#timer").html(formatTime(timeElapsed));
	} else {
		timeRemaining = moment.duration(parseInt(schedule[step].Duration), "minutes");
		$("#timer").html(formatTime(timeRemaining));

		if(schedule[step].Type === "BOIL") {
			$("#temperature0Label").html("Wort Temp");
			$("#temperature1Label").html("Unused");
		} else if(schedule[step].Type === "MASH" || schedule[step].Type === "SPARGE") {
			$("#temperature0Label").html("Mash Temp");
			$("#temperature1Label").html("Sparge Temp");
		}
	}
}
showCurrentStep();

function checkTemperatureRanges(t0, t1) {
	if(schedule[step].Type === "MASH") {
		var spargeTemp = getFirstSpargeTemp();

		if(!temp0InRange && Math.abs(schedule[step].Temperature - t0) <= temperatureThreshold) {
			sendNotification("Mash temperature is within desired range.", null, null, false);
			temp0InRange = true;
		}
		else if(!temp1InRange && Math.abs(spargeTemp - t1) <= temperatureThreshold) {
			sendNotification("Sparge temperature is within desired range.", null, null, false);
			temp1InRange = true;
		}
		else if(temp0InRange && Math.abs(schedule[step].Temperature - t0) > temperatureThreshold) {
			temp0InRange = false;
			if(schedule[step].Temperature - t0 < temperatureThreshold) {
				sendNotification("Mash temperature exceeds threshold", "ALERT", null, false);
			} else if(schedule[step].Temperature - t0 > -temperatureThreshold) {
				sendNotification("Mash temperature dropped below threshold", "ALERT", null, false);
			}
		}
		else if(temp1InRange && Math.abs(spargeTemp - t1) > temperatureThreshold) {
			if(spargeTemp - t1 < temperatureThreshold) {
				sendNotification("Sparge temperature exceeds threshold", "ALERT", null, false);
			} else if(spargeTemp - t1 > -temperatureThreshold) {
				sendNotification("Sparge temperature dropped below threshold", "ALERT", null, false);
			}
		}		
	}
	else if(schedule[step].Type === "SPARGE") {
		temp0InRange = false;

		if(schedule[step+1].Type === "SPARGE") {
			var spargeTemp = schedule[step+1].Temperature;

			if(!temp1InRange && Math.abs(spargeTemp - t1) <= temperatureThreshold) {
				sendNotification("Sparge temperature is within desired range.", null, null, false);
				temp1InRange = true;
			}
			else if(temp1InRange && Math.abs(spargeTemp - t1) > temperatureThreshold) {
			if(spargeTemp - t1 < temperatureThreshold) {
				sendNotification("Sparge temperature exceeds threshold", "ALERT", null, false);
			} else if(spargeTemp - t1 > -temperatureThreshold) {
				sendNotification("Sparge temperature dropped below threshold", "ALERT", null, false);
			}
		}	
		}
	}
	else if(schedule[step].Type === "BOIL") {
		if(!temp0InRange && t0 >= 200) {
			temp0InRange = true;
			sendNotification("Approaching boiling temperature.");
		}
	}
	else if(schedule[step].Type === "CHILL") {
		if(schedule[step].Temperature > t0) {
			sendNotification("Chill temperature is within desired range.");
		}
	}
}

function refreshSensorReadings() {
	$.get(
		"/get-temperature-readings",
		function(data) {
			var temps = $.parseJSON(data);
			$("#temperature0").html(temps[0] + "&deg;F");
			$("#temperature1").html(temps[1] + "&deg;F");
			checkTemperatureRanges(temps[0], temps[1]);
		}
	);
}
window.setInterval(refreshSensorReadings(),5000);

function countdownTimer() {
	$("#timer").html(formatTime(timeRemaining));
	timeRemaining.subtract(1, "seconds");
	if(timeRemaining.seconds() <= 0) {
		if(schedule[step].Type === "MASH") { sendNotification("Mash Step Complete.", "DONE", null, null); }
		else if(schedule[step].Type === "SPARGE") { sendNotification("Sparge Step Complete.", "DONE", null, null); }
		else if(schedule[step].Type === "BOIL") { sendNotification("Boil Step Complete.", "DONE", null, null); }
		step++;
		showCurrentStep();

		if(autoadvance) {
			startTimer();
		} else {
			stopTimer()
			$("#start").show();
			$("#pause").hide();
			$("#resume").hide();
		}
	}
}

function countupTimer() {
	$("#timer").html(formatTime(timeElapsed));
	timeElapsed.add(1, "seconds")
}

function startTimer() {
	clearInterval(timer);
	if(schedule[step].Type === "CHILL") {
		timer = window.setInterval(countupTimer, 1000);
	} else {
		timer = window.setInterval(countdownTimer, 1000);
	}
}

function stopTimer() {
	clearInterval(timer);
}

$("#start").click(function() {
	$("#start").hide();
	$("#pause").show();
	startTimer();
});

$("#pause").click(function() {
	$("#resume").show();
	$("#pause").hide();
	stopTimer();
});

$("#resume").click(function() {
	$("#resume").hide();
	$("#pause").show();
	startTimer();
});

function sendAlert(subject, message) {
	$.post(
		"/send-alert",
		JSON.stringify({ Subject : subject, Message : message }),
		function(data) {}
	);
}
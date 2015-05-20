var schedule = {
	MashSteps : [],
	SpargeSteps : [],
	BoilSteps : [],
	ChillSteps : []
}

function mashUI() {
	$("#previous").off("click").hide();
	$("#add").off("click").click(addMashStep).html("add mash step");
	$("#next").off("click").click(spargeUI).html("go to sparge");
}

function addMashStep() {
	$row = $(".scheduleStep.template").clone();
	$row.removeClass("template").addClass("mashStep");;;
	$row.find(".stepType").html("Mash Step #" + (schedule.MashSteps.length + 1));
	$row.find(".temperature").html($("#temperature").val() + "&deg;F");
	$row.find(".duration").html($("#duration").val() + "min");
	$row.show();
	$("tr.mashStep:last").after($row);

	schedule.MashSteps.push({ Temperature : $("#temperature").val(), Duration : $("#duration").val() } );
}

function spargeUI() {
	$(".temperatureInput").show();
	$(".durationInput").show();
	$("#previous").show().off("click").click(mashUI).html("back to mash");
	$("#add").off("click").click(addSpargeStep).html("add sparge step");
	$("#next").off("click").click(boilUI).html("go to boil");
}

function addSpargeStep() {
	$row = $(".scheduleStep.template").clone();
	$row.removeClass("template").addClass("spargeStep");
	$row.find(".stepType").html("Sparge Step #" + (schedule.SpargeSteps.length + 1));
	$row.find(".temperature").html($("#temperature").val() + "&deg;F");
	$row.find(".duration").html($("#duration").val() + "min");
	$row.show();
	$("tr.spargeStep:last").after($row);

	schedule.SpargeSteps.push({ Temperature : $("#temperature").val(), Duration : $("#duration").val() });
}

function boilUI() {
	$(".temperatureInput").hide();
	$(".durationInput").show();
	$("#previous").show().off("click").click(spargeUI).html("back to sparge");;
	$("#add").off("click").click(addBoilStep).html("add boil step");
	$("#next").off("click").click(chillUI).html("go to chill");
}

function addBoilStep() {
	$row = $(".scheduleStep.template").clone();
	$row.removeClass("template").addClass("boilStep");
	$row.find(".stepType").html("Boil Step #" + (schedule.BoilSteps.length + 1));
	$row.find(".temperature").html("--");
	$row.find(".duration").html($("#duration").val() + "min");
	$row.show();
	$("tr.boilStep:last").after($row);

	schedule.BoilSteps.push($("#duration").val());
}

function chillUI() {
	$(".temperatureInput").show();
	$(".durationInput").hide();
	$("#previous").show().off("click").click(boilUI).html("back to boil");;
	$("#add").off("click").click(addChillStep).html("add chill step");
	$("#next").off("click").click(finishAndSave).html("save schedule");
}

function addChillStep() {
	$row = $(".scheduleStep.template").clone();
	$row.removeClass("template").addClass("chillStep");
	$row.find(".stepType").html("Chill Step #" + (schedule.ChillSteps.length + 1));
	$row.find(".temperature").html($("#temperature").val() + "&deg;F");
	$row.find(".duration").html("--");
	$row.show();
	$("tr.chillStep:last").after($row);

	schedule.ChillSteps.push($("#temperature").val());
}

function finishAndSave() {
	$.post(
		"/save-schedule",
		JSON.stringify(schedule),
		function(data) {
			if(data === "success") {
				sendNotification("Schedule Saved.");
			}
		}
	);
}

mashUI();

$("#temperature").mousemove(function() { $("#temperatureValue").html($("#temperature").val()); });
$("#temperature").change(function() { $("#temperatureValue").html($("#temperature").val()); });
$("#duration").mousemove(function() { $("#durationValue").html($("#duration").val()); });
$("#duration").change(function() { $("#durationValue").html($("#duration").val()); });

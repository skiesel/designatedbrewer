var selectedFile = "";

$("#select").click(function() {
	$.post(
		"load-file",
		JSON.stringify({ Filename : $("#file").val() }),
		function(data, textStatus, jqXHR) {
			$("#mashSteps").empty();
			$("#spargeSteps").empty();
			$("#boilSteps").empty();
			$("#chillSteps").empty();

			$("#confirm").show();
			selectedFile = $("#file").val();
			displayFile(JSON.parse(data));
		}
	);
});

$("#confirm").click(function() {
	$("#hiddenPostForm input").val(selectedFile);
	$("#hiddenPostForm").submit();
});

function displayFile(file) {
	$.each(file.MashSteps, function(index, value) {
		$row = $(".scheduleStep.template").clone();
		$row.removeClass("template").addClass("mashStep");
		$row.find(".stepType").html("Mash Step #" + (index + 1));
		$row.find(".temperature").html(value.Temperature + "&deg;F");
		$row.find(".duration").html(value.Duration + "min");
		$row.show();
		$("tr.mashStep:last").after($row);
	});

	$.each(file.SpargeSteps, function(index, value) {
		$row = $(".scheduleStep.template").clone();
		$row.removeClass("template").addClass("spargeStep");
		$row.find(".stepType").html("Sparge Step #" + (index + 1));
		$row.find(".temperature").html(value.Temperature + "&deg;F");
		$row.find(".duration").html(value.Duration + "min");
		$row.show();
		$("tr.spargeStep:last").after($row);
	});

	$.each(file.BoilSteps, function(index, value) {
		$row = $(".scheduleStep.template").clone();
		$row.removeClass("template").addClass("boilStep");
		$row.find(".stepType").html("Boil Step #" + (index + 1));
		$row.find(".temperature").html("--");
		$row.find(".duration").html(value + "min");
		$row.show();
		$("tr.boilStep:last").after($row);
	});

	$.each(file.ChillSteps, function(index, value) {
		$row = $(".scheduleStep.template").clone();
		$row.removeClass("template").addClass("chillStep");
		$row.find(".stepType").html("Chill Step #" + (index + 1));
		$row.find(".temperature").html(value + "&deg;F");
		$row.find(".duration").html("--");
		$row.show();
		$("tr.chillStep:last").after($row);
	});
}
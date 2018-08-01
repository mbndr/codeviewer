package codeviewer

// html template to display
var html = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <title>Codeviewer - {{.Filename}}</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
	<link id="theme" rel="stylesheet" href="style/{{.Theme}}.min.css">
    <style>
        body {
            margin: 0;
            padding: 0;
		}
		pre {
			margin: 0;
		}
		#overlay {
			background-color: #2A2A2A;
			color: #FFFFFF;
			font-family: monospace;
			padding: 10px;
			position: absolute;
			top: 10px;
			right: 10px;
			border-radius: 5px;
			border: solid 1px #202020;
		}
		#overlay hr {
			border: 0;
			height: 2px;
			margin: 10px 0;
			background: #202020;
		}
    </style>
</head>
<body class="hljs">

	<div id="overlay">
		Choose your theme<br><br>
		<select id="theme-chooser">
			{{range $n, $l := .ThemeList}}
				<option value="{{$n}}" {{if eq $n $.Theme}}selected{{end}} >{{$l}}</option>
			{{end}}
		</select>
		<hr>
		<!-- Metadata -->
		File: {{.Filename}}<br>
		Language: {{.Language}}
	</div>

    <pre><code class="{{.Language}}">{{.Content}}</code></pre>

	<script src="highlight.min.js"></script>
	<script src="lang/{{.Language}}.min.js"></script>
	<script>hljs.initHighlightingOnLoad();</script>
	<script>
	window.onload = function() {
		// THEME CHOOSER
		var themeChooser = document.getElementById("theme-chooser");
		themeChooser.onchange = function() {
			var val = themeChooser.options[themeChooser.selectedIndex].value;
			setStyle(val);
		};
	
		// OVERLAY
		var overlay = document.getElementById("overlay");
	
		document.addEventListener("keydown", function(e) {
			if (e.keyCode != 17) { // "CTRL"
				return;
			}
	
			console.log("toggle overlay")
	
			if (overlay.style.visibility == "hidden") {
				overlay.style.visibility = "visible";
			} else {
				overlay.style.visibility = "hidden"
			}
		});
	}
	
	function setStyle(name) {
		var styleLink = document.getElementById("theme");
		styleLink.href = "style/" + name + ".min.css";
	
		var xhr = new XMLHttpRequest();
		xhr.open('GET', "set-style/" + name, true);
		xhr.send();
	}	
	</script>
</body>
</html>
`
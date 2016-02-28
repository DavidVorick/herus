<html>
<head>
	<title>{{.Title}}</title>
</head>

<body>
	<a href='{{.MediaPrefix}}{{.Hash}}'>{{.Hash}}</a><br>
	<br>
	<h2>Annotations and Elaborations</h2>
	{{range .Elaborations}}
		<a href='{{$.ElaborationPrefix}}{{.Hash}}'>{{.Title}}</a><br>
	{{end}}
</body>
</html>

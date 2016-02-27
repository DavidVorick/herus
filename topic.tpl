<html>
<head>
	<title>{{.TopicTitle}}</title>
</head>

<body>
	<center><h1>{{.TopicTitle}}</h1></center>
	{{range .SubmittedMedia}}
		<a href='{{.MediaHash}}'>{{.MediaTitle}}</a><br>
	{{end}}
</body>
</html>

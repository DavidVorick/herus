<html>
<head>
	<title>{{.TopicTitle}}</title>
</head>

<body>
	<center><h1>{{.TopicTitle}}</h1></center>
	{{range .AssociatedMedia}}
		<a href='{{.Hash}}'>{{.Title}}</a><br>
	{{end}}
</body>
</html>

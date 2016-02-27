<html>
<head>
	<title>{{.TopicTitle}}</title>
</head>

<body>
	<center><h1>{{.TopicTitle}}</h1></center>
	{{range .AssociatedMedia}}
		<a href='{{$.MediaPrefix}}{{.Hash}}'>{{.Title}}</a><br>
	{{end}}
	<br><br>
	<center><h2>Related Pages</h2></center>
	{{range .RelatedTopics}}
		<a href='/t/{{.TopicTitle}}'>{{.TopicTitle}}</a><br>
	{{end}}
</body>
</html>

	<center><h1>{{.Title}}</h1></center>
	{{range .AssociatedMedia}}
		<a href='{{$.ElaborationPrefix}}{{.Hash}}'>{{.Title}}</a><br>
	{{end}}
	<br><br>
	<center><h2>Related Pages</h2></center>
	{{range .RelatedTopics}}
		<a href='{{$.TopicPrefix}}{{.Title}}'>{{.Title}}</a><br>
	{{end}}

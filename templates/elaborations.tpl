	<br><br><br>
	<center><a href='{{.MediaPrefix}}{{.Hash}}'>{{.Hash}}</a><br></center>
	<br>
	<center><h2>Annotations and Elaborations</h2></center>
	<center>
		{{range .Elaborations}}
		<a href='{{$.ElaborationPrefix}}{{.Hash}}'>{{.Title}}</a><br>
		{{end}}
	</center>
	<br><br><br>
	<br><br><br>
	<br><br><br>

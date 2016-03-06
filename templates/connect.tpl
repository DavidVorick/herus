	{{if .ErrorExists}}
	<center><h2>{{.Error}}</h2></center><br><br>
	{{end}}
	{{if .PostWithoutError}}
	<center><h2>Thanks!</h2></center><br><br>
	{{end}}

	<form action='connect.go' method='post'>
		Source Topic: <input type='test' name='sourceTopic'><br>
		Destination Topic: <input type='text' name='destinationTopic'></br>
		<input type='submit' value='Connect'><br>
	</form>

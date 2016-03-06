	<form action='connect.go' method='post'>
		Source Topic: <input type='test' name='sourceTopic'><br>
		Destination Topic: <input type='text' name='destinationTopic'></br>
		<input type='submit' value='Connect'><br>
	</form>
	{{if .}}<center><h1>Thanks!</h1></center><br><br>{{end}}

<html>
<head>
	<title>Upload Media</title>
</head>

<body>
<p>Uploads must have either a parent topic or a parent media. The same media <br>
can have multiple parents, but must be added as separate uploads.</p>
<form enctype="multipart/form-data" action="upload" method="post">
	File: <input type="file" name="upload"><br>
	File Title: <input type="test" name="title"><br>
	Parent Topic: <input type="text" name="parentTopic"></br>
	Parent Media: <input type="text" name="parentMedia"></br>
	<input type="submit" value="upload"><br>
</form>
{{if .}}<center><h1>Thanks!</h1></center><br><br>{{end}}
</body>
</html>

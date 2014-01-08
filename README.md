## Simple appengine datastore dump and restore

To use this dump and restore service, import the go part in your program as follows:

	import _ github.com/jum/dsdump

This will initialize an HTTP handler for /admin/dsdump in your app engine app. You should use something like the following in your app.yaml to protect any URL under /admin:

	handlers:
	- url: /admin/.*
	  script: _go_app
	  login: admin

The two shell scripts named dsdump and dsrestore in the scripts folder can be used to dump and restore data from the server. Use at follows for a local development server:

	dsdump http://localhost:8080 >data.json

or add admin email and password for a real appengine instance:

	dsdump https://appid.appspot.com admin@domain.com password >data.gob

And accordingly to restore to a local developsment server:

	dsrestore http://localhost:8080 data.gob
	
And to restore to an appengine instance:

	dsrestore https://appid.appspot.com admin@domain.com password data.gob

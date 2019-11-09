# GO NOTIFIER
This package can be used to send desktop notifications from any go applications or packages.

##### How to use?
1.  Get the package
		go get -v github.com/satheesh1997/go-notifier
2.  Import to your source code file
		import (
			"github.com/satheesh1997/go-notifier"
		)

##### How to send notification?
1.  To send notification use
		notifier.Notify(title, message, level)
2. Levels can be:
	1.  notifier.LOW
	2.  notifier.NORMAL
	3.  notifier.CRITICAL
	
##### Supported OS:
1.  Linux

we will extend our support soon for other operating system

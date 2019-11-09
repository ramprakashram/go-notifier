# go-notifier
Go package that can be used to send notifications from go application


To use in your application:
  1. use go get -v github.com/satheesh1997/go-notifier
  2. add 
    import (
      ...
      "github.com/satheesh1997/go-notifier"
   )
  3. to send notification
     notifier.Notify(title, message, level)

# ocgen
This is an Objective C boilerplate code generator. It has been developed by demand.
Right now only generates the methods needed to conform the `NSCopying` and `NSCoding` protocols and 
it only takes care of properties declared in the class interface.

## How to use
If you don't have Go installed on your system, you can grab directly the executable for 
your working platform. These are supported platforms:

* Mac OSX 64bits
* Linux 64bits

If you have a distribution of Go installed, you can "go get" it directly

    go get github.com/alvaroloes/ocgen
    
## How it works
Just call it with the name of your `.h` class file as parameter.

    ocgen /path/to/class.h
    
It will output the code needed to conform the protocols. Just copy and paste it into you `.m` class file
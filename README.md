# OCGen
This is an Objective C boilerplate code generator. It has been developed by demand.
Right now only generates the methods needed to conform the `NSCopying` and `NSCoding` protocols.
It's in an early state and there are still many things to add, as you can see in the TODO section.

Any pull request will be highly appreciated.

## Usage

	ocgen [options] directory1 [directory2,...]
  		-backup=true: Whether to create a backup of all files before modifying them
  		-backupDir="./.ocgen": The directory where the backups will be placed if 'backup=true'

It will generate the methods to conform the `NSCopying` and `NSCoding` protocols for the classes under each directory. Only the classes tagged with "OCGEN_AUTO" will be considered. Right now all properties are taken into account for each generated method

## How to tag a class with OCGEN_AUTO?
You first need to create an empty macro:

	#define OCGEN_AUTO

Then you need to use that macro to tag the class interface, putting it at the end of the `@interface` line:

	@interface MyClass : NSObject <NSCopying, NSCoding> OCGEN_AUTO
		(...)
	@end


## TODO
### High priority
* Call the super in the generated methods if it responds to the method selector
* Restore the backed file if there was an error in the write operation inside the `GenerateMethods`
* Decide how to copy items based on property attributes (and type?)
* How to handle classes that conforms to the protocols indirectly through another protocol?
* Merge properties from header and implementation file (taking care of readonly ones)
* Add instruction about how to execute it in every compilation

### Medium priority
* Add concurrency
* Allow to specify a directory to store the backups

### Lower priority
* Allow installing through Alcatraz
* Provide a header file with the macros used for tagging

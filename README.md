# OCGen
This is an Objective C boilerplate code generator. It has been developed by demand.
Right now only generates the methods needed to conform the `NSCopying` and `NSCoding` protocols.
It is under heavy development and in an unfinished state. You can use the previous version (a different approach) in the branch `initial_approach`

## TODO
### High priority
* Finish main function, defining the supported flags
* Call the super in the generated methods if it responds to the method selector
* Restore the backed file if there was an error in the write operation inside the `GenerateMethods`
* Decide how to copy items based on property attributes (and type?)
* How to handle classes that conforms to the protocols indirectly through another protocol?

### Medium priority
* Add concurrency
* Allow to specify a directory to store the backups

### Lower priority
* Allow installing through Alcatraz

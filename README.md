# OCGen
The Objective C code generator.

Developed by demand due to my job needs, this tool aims to generate boring and repetitive code to let you focus on more
interesting and pleasant things.

Right now it is able to generate all the code a class needs to conform the `NSCopying` and `NSCoding` protocols.
In the future, it will be able to generate more code, such as:
 
* Generating the corresponding Objective C classes from a `JSON/XML` object
* Generating the code needed to parse from/to `JSON/XML`. Useful when you don't want to use runtime techniques (*aka* reflexion) 
because of performance or any other reason
* Generating an Objective C mobile SDK from an API specification. This is, by far, one of the most useful (and hard!) things this 
tool could ever do :-)
* Generating code for other mobile languages, such us Android, Swift or Go
* &lt;Put your awesome idea here&gt;

## Disclaimer
*This tool is under development and unexpected behavior may arise. Please be sure your project is using a version control system, like git, 
before using OCGen, so you can revert any unwanted changes the tool has done.*


## Usage
### Download
There are two ways of getting OCGen:

1. Go to _Releases_ github section and download the latest version. 
1. If you have the Go compiler, you can either compile it from source or "go-get" it:
```
go get github.com/alvaroloes/ocgen
```

*NOTE:* Please read the section _"Contribution or compiling from source"_ at the bottom of the page before choosing the second option.


### Basic
Just call `ocgen` specifying the directory under which your classes are (it could be the root of your project or any other
more specific folder, such as `/models`)
	
	ocgen ./MyProject/models

All the classes that conforms `NSCopying` or `NSCoding` will now have the corresponding methods autogenerated.
For example, if we had the class `User` under the previous directory with the following interface:
```objective-c
@interface User : NSObject <NSCoding, NSCopying>    

    @property (nonatomic, copy) NSString *name;
    @property (nonatomic, strong) NSNumber *age;
    @property (nonatomic, strong) NSArray *telephones;
    @property (nonatomic, assign) BOOL isAdmin;
    @property (nonatomic, strong, readonly) NSString *gender;
    @property (nonatomic, weak) id<GroupProtocol> belongingGroup;

@end
```
    
The implementation now looks like this:
```objective-c
@implementation User

// (...)

- (instancetype)copyWithZone:(NSZone *)zone
{
    // OCGEN: Autogenerated method
    User *copy = [[[self class] allocWithZone:zone] init];
    if (copy != nil)
    { 
        copy.name = [self.name copyWithZone:zone];
        copy.age = [self.age copyWithZone:zone];
        copy.telephones = [[NSArray alloc] initWithArray:self.telephones copyItems:YES];
        copy.isAdmin = self.isAdmin;
        copy->_gender = [self.gender copyWithZone:zone];
        copy.belongingGroup = self.belongingGroup;
    }
    return copy;
}

- (instancetype)initWithCoder:(NSCoder *)decoder
{
    // OCGEN: Autogenerated method
    self = [super init];
    if(self != nil)
    { 
        _name = [decoder decodeObjectForKey:@"name"];
        _age = [decoder decodeObjectForKey:@"age"];
        _telephones = [decoder decodeObjectForKey:@"telephones"];
        _isAdmin = [decoder decodeBoolForKey:@"isAdmin"];
        _gender = [decoder decodeObjectForKey:@"gender"];
        _belongingGroup = [decoder decodeObjectForKey:@"belongingGroup"];
    }
    return self;
}

- (void)encodeWithCoder:(NSCoder *)coder
{
    // OCGEN: Autogenerated method
    [coder encodeObject:self.name forKey:@"name"];
    [coder encodeObject:self.age forKey:@"age"];
    [coder encodeObject:self.telephones forKey:@"telephones"];
    [coder encodeBool:self.isAdmin forKey:@"isAdmin"];
    [coder encodeObject:self.gender forKey:@"gender"];
    [coder encodeConditionalObject:self.belongingGroup forKey:@"belongingGroup"];
}

// (...) 

@end
```

As you can see, OCGen takes into account the type of the property (and if it is readonly or not) to generate valid code.
If the class have previous declarations of the methods, they will be replaced.

### Ignoring classes
If you want OCGen to ignore a class (because you want to add special coding or copying behavior, for example), you
can tag it with the `OCGEN_IGNORE` macro and it won't be processed.

To do so, you first need to define an empty macro with the tag name (you can do it in the `.pch` file to make it globally visible):
```objective-c
#define OCGEN_IGNORE
```
    
Then write that tag at the end of the line where you have the `@interface` declaration, like so:
```objective-c
@interface User : NSObject <NSCoding, NSCopying> OCGEN_IGNORE   

    @property (nonatomic, copy) NSString *name;
    @property (nonatomic, strong) NSNumber *age;
    @property (nonatomic, strong) NSArray *telephones;
    @property (nonatomic, assign) BOOL isAdmin;
    @property (nonatomic, strong, readonly) NSString *gender;

@end
```
    
This way the class will be ignored 

### Advanced

Here you can see the usage and all the supported options: 
```bash
ocgen [options] directory1 [directory2,...]
  -backup
        Whether to create a backup of all files before modifying them
        
  -backupDir string
        The directory where the backups will be placed if 'backup' is present (default "./.ocgen")
        
  -extraNSCodingProtocols string
        A comma separated list (without spaces) of protocol names that will be considered as if they were NSCoding. 
        This is useful if your class does not conform NSCoding directly, but through another protocol that conforms it. 
        Example: extraNSCodingProtocols="MyProtocolThatConformsNSCoding,OtherProtocolThatConformsNSCoding"
        
  -extraNSCopyingProtocols string
        A comma separated list (without spaces) of protocol names that will be considered as if they were NSCopying. 
        This is useful if your class does not conform NSCopying directly, but through another protocol that conforms it. 
        Example: extraNSCopyingProtocols="MyProtocolThatConformsNSCopying,OtherProtocolThatConformsNSCopying"
        
  -h	Prints the usage
  
  -v	Prints the current version
```

## Use OCGen as a build phase
You can set OCGen as a build phase to always have your `NSCoding` or `NSCopying` code up to date.
You can follow the instructions here: http://www.runscriptbuildphase.com/ 

Now, each change you do to your classes properties will be reflected in the generated code on each compilation. 

## Limitations
* The `.m` file must be in the same directory than the `.h` file.
* The *ivars* that correspond to each property must have the same name than the property prefixed by `_` (this is the default
behavior in Objective C)
* Read only properties not backed by an *ivar* (with a defined getter) are not supported.

## Contribution or compiling from source
If you find any bug or want a feature to be added, feel free to file an issue. 

If you want to compile the code from source or contribute with a pull request, please take into account that *you need to 
build Go directly from the master branch*. The reason behind this is that some improvements where added
to the `text/template` package after the Go 1.5 release that make the template code much more readable (specifically, the addition of {{- }} and {{ -}}
that allows newlines to be stripped from the generated text)

To build Go from source, you can follow the steps in the official site: https://golang.org/doc/install/source
 
Thanks!

## TODO
* Add executable options to customize the tags used to mark the classes
* Be able to restore a class from the backed up version
* Add concurrency
* Allow installing through Alcatraz
* Provide a header file with the macros used for tagging

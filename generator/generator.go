package generator
import (
    "github.com/alvaroloes/ocgen/parser"
    "bytes")


func NSCopying(info *parser.ObjCClassInfo) (string, error) {
    var res bytes.Buffer;
    err := NSCopyingTpl.Execute(&res, info);
    return res.String(), err
}

func NSCoding(info *parser.ObjCClassInfo) (string, error) {
    var res bytes.Buffer;
    err := NSCodingInitTpl.Execute(&res, info);
    return res.String(), err
}

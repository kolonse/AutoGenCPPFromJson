// AutoGenCPPFromJson project doc.go

/*
AutoGenCPPFromJson document
*/
package main

/**
使用方式 :AutoGenCPPFromJson.exe [json file] [GenName] [GenDir]
[json file] 表示需要转化成CPP代码的 json源文件
[GenName] 表示需要生成的类名
[GenDir] 表示类生成的路径

细节注意:
1.对于json数组  暂时只支持数组类型为单一的类型  不支持混合类型  所以限定了源文件的数组 必须且只有一个元素
example:
[
	"test"
]

2.对于整数或者浮点数类型  必须填写值 0 1 2 3 4 其中之一
分别表示 int,uint,int64,uint64,double
example:
{
	"int":0,
}

*/

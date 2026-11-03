^XA
^FX asset-tag/2x1/001 - 2 x 1 in @ 203 dpi: company name, Code 128, human-readable tag
^CI28
^PW406
^LL203

^FX --- Company name ---
^FO34,22
{{$sz := fitfont .CompanyName 360 26 16}}^A0N,{{$sz}},{{$sz}}
^FD{{.CompanyName}}^FS

^FX --- Code 128 asset tag ---
^FO34,56
^BY2
^BCN,80,N,N,N,A
^FD{{.Tag}}^FS

^FX --- Human-readable asset tag ---
^FO34,145
^A0N,40,40
^FD{{.Tag}}^FS

^XZ

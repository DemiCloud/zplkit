^XA
^FX asset-tag/2x1/003 - 2 x 1 in @ 203 dpi: company name, return contact, Code 128, tag
^CI28
^PW406
^LL203

^FX --- Company name (fitfont keeps it off the label edge) ---
^FO34,18
{{$sz := fitfont .CompanyName 360 24 14}}^A0N,{{$sz}},{{$sz}}
^FD{{.CompanyName}}^FS

^FX --- Return contact ---
^FO34,45
^A0N,22,22
^FDIf found: found@example.com^FS

^FX --- Code 128 asset tag ---
^FO34,73
^BY2
^BCN,70,N,N,N,A
^FD{{.Tag}}^FS

^FX --- Human-readable asset tag ---
^FO34,151
^A0N,36,36
^FD{{.Tag}}^FS

^XZ

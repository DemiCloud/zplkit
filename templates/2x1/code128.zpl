^XA
^FX code128 - 2 x 1 in @ 203 dpi: company name, Code 128, human-readable tag
^CI28
^PW406
^LL203

^FX --- Company name ---
^FO30,15
^A0N,26,26
^FD{{.CompanyName}}^FS

^FX --- Code 128 asset tag ---
^FO30,50
^BY2
^BCN,80,N,N,N,A
^FD{{.Tag}}^FS

^FX --- Human-readable asset tag ---
^FO30,145
^A0N,40,40
^FD{{.Tag}}^FS

^XZ

^XA
^FX contact-code128 - 2 x 1 in @ 203 dpi: company name, return contact, Code 128, tag
^CI28
^PW406
^LL203

^FX --- Company name ---
^FO30,10
^A0N,24,24
^FD{{.CompanyName}}^FS

^FX --- Return contact ---
^FO30,36
^A0N,22,22
^FDIf found: found@example.com^FS

^FX --- Code 128 asset tag ---
^FO30,64
^BY2
^BCN,70,N,N,N,A
^FD{{.Tag}}^FS

^FX --- Human-readable asset tag ---
^FO30,145
^A0N,36,36
^FD{{.Tag}}^FS

^XZ

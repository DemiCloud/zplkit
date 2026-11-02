^XA
^FX mail-from/2.25x1.25/002 - 2.25 x 1.25 in @ 203 dpi: square border, company name, name, street address, city/state/ZIP
^CI28
^PW456
^LL254
^LH0,0

^FX --- Border ---
^FO12,12
^GB432,230,4^FS

^FX --- Company name ---
^FO38,69
^A0N,26,26
^FD{{.CompanyName}}^FS

^FX --- Sender name ---
^FO38,99
^A0N,30,30
^FD{{.Name}}^FS

^FX --- Street address ---
^FO38,133
^A0N,24,24
^FD{{.Address1}}^FS

^FX --- City, state, ZIP ---
^FO38,161
^A0N,24,24
^FD{{.City}}, {{.State}} {{.ZIP}}^FS

^XZ

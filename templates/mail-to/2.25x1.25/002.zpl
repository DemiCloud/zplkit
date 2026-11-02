^XA
^FX mail-to/2.25x1.25/002 - 2.25 x 1.25 in @ 203 dpi: square border, name, street address, city/state/ZIP
^CI28
^PW456
^LL254
^LH0,0

^FX --- Border ---
^FO12,12
^GB432,230,4^FS

^FX --- Recipient name ---
^FO38,72
^A0N,32,32
^FD{{.Name}}^FS

^FX --- Street address ---
^FO38,112
^A0N,26,26
^FD{{.Address1}}^FS

^FX --- City, state, ZIP ---
^FO38,146
^A0N,26,26
^FD{{.City}}, {{.State}} {{.ZIP}}^FS

^XZ

^XA
^FX mail-to/2.25x1.25/002 - 2.25 x 1.25 in @ 203 dpi: square border, name, street address, city/state/ZIP
^CI28
^PW456
^LL254
^LH0,0

^FX --- Border ---
^FO12,12
^GB432,230,4^FS

^FX --- Recipient name (fitfont keeps it off the border) ---
^FO38,72
{{$sz := fitfont .Name 402 32 20}}^A0N,{{$sz}},{{$sz}}
^FD{{.Name}}^FS

^FX --- Street address (fitfont keeps it off the border) ---
^FO38,112
{{$sz := fitfont .Address1 402 26 16}}^A0N,{{$sz}},{{$sz}}
^FD{{.Address1}}^FS

^FX --- City, state, ZIP (fitfont keeps it off the border) ---
^FO38,146
{{$sz := fitfont (printf "%s, %s %s" .City .State .ZIP) 402 26 16}}^A0N,{{$sz}},{{$sz}}
^FD{{.City}}, {{.State}} {{.ZIP}}^FS

^XZ

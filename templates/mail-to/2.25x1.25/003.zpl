^XA
^FX mail-to/2.25x1.25/003 - 2.25 x 1.25 in @ 203 dpi: rounded border, name, street address, city/state/ZIP
^CI28
^PW456
^LL254
^LH0,0

^FX --- Border with rounded corners ---
^FO12,12
^GB432,230,4,B,3^FS

^FX --- Recipient name ---
^FO38,72
{{$sz := fitfont .Name 402 32 20}}^A0N,{{$sz}},{{$sz}}
^FD{{.Name}}^FS

^FX --- Street address ---
^FO38,112
{{$sz := fitfont .Address1 402 26 16}}^A0N,{{$sz}},{{$sz}}
^FD{{.Address1}}^FS

^FX --- City, state, ZIP ---
^FO38,146
{{$sz := fitfont (printf "%s, %s %s" .City .State .ZIP) 402 26 16}}^A0N,{{$sz}},{{$sz}}
^FD{{.City}}, {{.State}} {{.ZIP}}^FS

^XZ

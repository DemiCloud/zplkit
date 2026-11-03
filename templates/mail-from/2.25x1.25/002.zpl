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
{{$sz := fitfont .CompanyName 402 26 16}}^A0N,{{$sz}},{{$sz}}
^FD{{.CompanyName}}^FS

^FX --- Sender name ---
^FO38,99
{{$sz := fitfont .Name 402 30 18}}^A0N,{{$sz}},{{$sz}}
^FD{{.Name}}^FS

^FX --- Street address ---
^FO38,133
{{$sz := fitfont .Address1 402 24 14}}^A0N,{{$sz}},{{$sz}}
^FD{{.Address1}}^FS

^FX --- City, state, ZIP ---
^FO38,161
{{$sz := fitfont (printf "%s, %s %s" .City .State .ZIP) 402 24 14}}^A0N,{{$sz}},{{$sz}}
^FD{{.City}}, {{.State}} {{.ZIP}}^FS

^XZ

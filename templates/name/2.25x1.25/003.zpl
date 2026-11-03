^XA
^FX name/2.25x1.25/003 - 2.25 x 1.25 in @ 203 dpi: rounded border, name, centered
^CI28
^PW456
^LL254
^LH0,0

^FX --- Border with rounded corners ---
^FO12,12
^GB432,230,4,B,3^FS

^FX --- Name (fitfont sizes to the longest line; a real newline in the
^FX field, e.g. -f Name=$'Cables:\nUSB-C', becomes a second line) ---
{{$n := linecount .Name}}^FO16,{{if gt $n 1}}79{{else}}103{{end}}
{{$sz := fitfont .Name 424 48 20}}^A0N,{{$sz}},{{$sz}}
^FB424,2,0,C
^FD{{zplbreak .Name}}^FS

^XZ

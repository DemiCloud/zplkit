^XA
^FX name/2.25x1.25/002 - 2.25 x 1.25 in @ 203 dpi: square border, name, centered
^CI28
^PW456
^LL254
^LH0,0

^FX --- Border ---
^FO12,12
^GB432,230,4^FS

^FX --- Name (fitfont shrinks it to stay on one line within 424 dots) ---
^FO16,103
{{$sz := fitfont .Name 424 48 20}}^A0N,{{$sz}},{{$sz}}
^FB424,1,0,C
^FD{{.Name}}^FS

^XZ

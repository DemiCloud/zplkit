^XA
^FX name/2.25x1.25/001 - 2.25 x 1.25 in @ 203 dpi: name, centered
^CI28
^PW456
^LL254
^LH0,0

^FX --- Name (fitfont shrinks it to stay on one line within 424 dots) ---
^FO0,103
{{$sz := fitfont .Name 424 48 20}}^A0N,{{$sz}},{{$sz}}
^FB456,1,0,C
^FD{{.Name}}^FS

^XZ

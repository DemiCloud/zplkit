^XA
^FX name/2x1/001 - 2 x 1 in @ 203 dpi: name, centered
^CI28
^PW406
^LL203

^FX --- Name (fitfont shrinks it to stay on one line within 406 dots) ---
^FO0,70
{{$sz := fitfont .Name 406 48 20}}^A0N,{{$sz}},{{$sz}}
^FB406,1,0,C
^FD{{.Name}}^FS

^XZ

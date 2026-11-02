^XA
^FX name/2x1/003 - 2 x 1 in @ 203 dpi: rounded border, name, centered
^CI28
^PW406
^LL203

^FX --- Border with rounded corners ---
^FO8,8
^GB390,187,4,B,3^FS

^FX --- Name (fitfont shrinks it to stay on one line within 390 dots) ---
^FO8,70
{{$sz := fitfont .Name 390 48 20}}^A0N,{{$sz}},{{$sz}}
^FB390,1,0,C
^FD{{.Name}}^FS

^XZ

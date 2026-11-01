^XA
^FX code128-qr-border - 2.25 x 1.25 in @ 203 dpi: border, company name, human-readable tag, Code 128, QR to asset URL
^CI28
^PW456
^LL254
^LH0,0

^FX --- Border ---
^FO4,4
^GB448,246,4^FS

^FX --- Branding ---
^FO20,18
^A0N,30,30
^FD{{.CompanyName}} Asset^FS

^FX --- Human-readable asset tag ---
^FO20,60
^A0N,45,38
^FD{{.Tag}}^FS

^FX --- Code 128 asset tag ---
^FO20,125
^BY2,2,72
^BCN,72,N,N,N,A
^FD{{.Tag}}^FS

^FX --- QR direct asset URL ---
^FO300,28
^BQN,2,4
^FDQA,{{.AssetURL}}{{.Tag}}^FS

^XZ
